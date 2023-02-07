package warlock

import (
	"time"

	"github.com/Tereneckla/wotlk70/sim/core"
	"github.com/Tereneckla/wotlk70/sim/core/proto"
	"github.com/Tereneckla/wotlk70/sim/core/stats"
)

var TalentTreeSizes = [3]int{28, 27, 26}

type Warlock struct {
	core.Character
	Talents  *proto.WarlockTalents
	Options  *proto.Warlock_Options
	Rotation *proto.Warlock_Rotation

	procTrackers []*ProcTracker
	majorCds     []*core.MajorCooldown

	Pet *WarlockPet

	ShadowBolt         *core.Spell
	Incinerate         *core.Spell
	Immolate           *core.Spell
	ImmolateDot        *core.Dot
	UnstableAffliction *core.Spell
	Corruption         *core.Spell
	Haunt              *core.Spell
	LifeTap            *core.Spell
	DarkPact           *core.Spell
	ChaosBolt          *core.Spell
	SoulFire           *core.Spell
	Conflagrate        *core.Spell
	ConflagrateDot     *core.Dot
	DrainSoul          *core.Spell
	Shadowburn         *core.Spell

	CurseOfElements     *core.Spell
	CurseOfElementsAura *core.Aura
	CurseOfWeakness     *core.Spell
	CurseOfWeaknessAura *core.Aura
	CurseOfTongues      *core.Spell
	CurseOfTonguesAura  *core.Aura
	CurseOfAgony        *core.Spell
	CurseOfDoom         *core.Spell
	Seed                *core.Spell
	SeedDamageTracker   []float64

	NightfallProcAura      *core.Aura
	EradicationAura        *core.Aura
	DemonicEmpowerment     *core.Spell
	DemonicEmpowermentAura *core.Aura
	DemonicPactAura        *core.Aura
	DemonicSoulAura        *core.Aura
	Metamorphosis          *core.Spell
	MetamorphosisAura      *core.Aura
	ImmolationAura         *core.Spell
	HauntDebuffAura        *core.Aura
	MoltenCoreAura         *core.Aura
	DecimationAura         *core.Aura
	PyroclasmAura          *core.Aura
	BackdraftAura          *core.Aura
	EmpoweredImpAura       *core.Aura
	GlyphOfLifeTapAura     *core.Aura
	SpiritsoftheDamnedAura *core.Aura

	Infernal *InfernalPet
	Inferno  *core.Spell

	// Rotation related memory
	CorruptionRolloverPower float64
	DrainSoulRolloverPower  float64
	// The sum total of demonic pact spell power * seconds.
	DPSPAggregate  float64
	PreviousTime   time.Duration
	SpellsRotation []SpellRotation

	petStmBonusSP                float64
	masterDemonologistFireCrit   float64
	masterDemonologistShadowCrit float64

	CritDebuffCategory *core.ExclusiveCategory
}

type SpellRotation struct {
	Spell    *core.Spell
	CastIn   CastReadyness
	Priority int
}

type CastReadyness func(*core.Simulation) time.Duration

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func (warlock *Warlock) GrandSpellstoneBonus() float64 {
	return core.TernaryFloat64(warlock.Options.WeaponImbue == proto.Warlock_Options_GrandSpellstone, 0.01, 0)
}
func (warlock *Warlock) GrandFirestoneBonus() float64 {
	return core.TernaryFloat64(warlock.Options.WeaponImbue == proto.Warlock_Options_GrandFirestone, 0.01, 0)
}

func (warlock *Warlock) Initialize() {

	warlock.registerIncinerateSpell()
	warlock.registerShadowBoltSpell()
	warlock.registerImmolateSpell()
	warlock.registerCorruptionSpell()
	warlock.registerCurseOfElementsSpell()
	warlock.registerCurseOfWeaknessSpell()
	warlock.registerCurseOfTonguesSpell()
	warlock.registerCurseOfAgonySpell()
	warlock.registerCurseOfDoomSpell()
	warlock.registerLifeTapSpell()
	warlock.registerSeedSpell()
	warlock.registerSoulFireSpell()
	warlock.registerUnstableAfflictionSpell()
	warlock.registerDrainSoulSpell()
	warlock.registerConflagrateSpell()
	warlock.registerHauntSpell()
	warlock.registerChaosBoltSpell()

	warlock.registerDemonicEmpowermentSpell()
	if warlock.Talents.Metamorphosis {
		warlock.registerMetamorphosisSpell()
		warlock.registerImmolationAuraSpell()
	}
	warlock.registerDarkPactSpell()
	warlock.registerShadowBurnSpell()
	warlock.registerInfernoSpell()

	warlock.defineRotation()
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.BloodPact = core.MaxTristate(raidBuffs.BloodPact, core.MakeTristateValue(
		warlock.Options.Summon == proto.Warlock_Options_Imp,
		warlock.Talents.ImprovedImp == 2,
	))

	raidBuffs.FelIntelligence = core.MaxTristate(raidBuffs.FelIntelligence, core.MakeTristateValue(
		warlock.Options.Summon == proto.Warlock_Options_Felhunter,
		warlock.Talents.ImprovedFelhunter == 2,
	))
}

func (warlock *Warlock) Prepull(sim *core.Simulation) {
	spellChoice := warlock.ShadowBolt
	if warlock.Rotation.Type == proto.Warlock_Rotation_Destruction {
		spellChoice = warlock.SoulFire
	}

	delay := (warlock.ApplyCastSpeed(core.GCDDefault) + warlock.ApplyCastSpeed(spellChoice.DefaultCast.CastTime))
	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) {
		warlock.GlyphOfLifeTapAura.Activate(sim)
		warlock.GlyphOfLifeTapAura.UpdateExpires(warlock.GlyphOfLifeTapAura.Duration - delay)
	}

	if warlock.SpiritsoftheDamnedAura != nil {
		warlock.SpiritsoftheDamnedAura.Activate(sim)
		warlock.SpiritsoftheDamnedAura.UpdateExpires(warlock.SpiritsoftheDamnedAura.Duration - delay)
	}

	warlock.SpendMana(sim, spellChoice.DefaultCast.Cost, spellChoice.Cost.(*core.ManaCost).ResourceMetrics)
	spellChoice.CD.UsePrePull(sim, warlock.ApplyCastSpeed(spellChoice.DefaultCast.CastTime))
	spellChoice.SkipCastAndApplyEffects(sim, warlock.CurrentTarget)
}

func (warlock *Warlock) Reset(sim *core.Simulation) {
	if sim.CurrentTime == 0 {
		warlock.petStmBonusSP = 0
	}
}

func NewWarlock(character core.Character, options *proto.Player) *Warlock {
	warlockOptions := options.GetWarlock()

	warlock := &Warlock{
		Character: character,
		Talents:   &proto.WarlockTalents{},
		Options:   warlockOptions.Options,
		Rotation:  warlockOptions.Rotation,
		// manaTracker:           common.NewManaSpendingRateTracker(),
	}
	core.FillTalentsProto(warlock.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)
	warlock.EnableManaBar()

	warlock.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	warlock.AddStatDependency(stats.Intellect, stats.SpellCrit, (1/81.92)*core.CritRatingPerCritChance)
	if warlock.Options.Armor == proto.Warlock_Options_FelArmor {
		demonicAegisMultiplier := 1 + float64(warlock.Talents.DemonicAegis)*0.1
		amount := 100.0 * demonicAegisMultiplier
		warlock.AddStat(stats.SpellPower, amount)
		warlock.AddStatDependency(stats.Spirit, stats.SpellPower, 0.3*demonicAegisMultiplier)
	}

	if warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		warlock.Pet = warlock.NewWarlockPet()
	}

	if warlock.Rotation.UseInfernal {
		warlock.Infernal = warlock.NewInfernal()
	}

	warlock.applyWeaponImbue()

	return warlock
}

func RegisterWarlock() {
	core.RegisterAgentFactory(
		proto.Player_Warlock{},
		proto.Spec_SpecWarlock,
		func(character core.Character, options *proto.Player) core.Agent {
			return NewWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Warlock)
			if !ok {
				panic("Invalid spec value for Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:    7164,
		stats.Strength:  48,
		stats.Agility:   60,
		stats.Stamina:   76,
		stats.Intellect: 136,
		stats.Spirit:    137,
		stats.Mana:      3856,
		stats.SpellCrit: 1.697 * core.CritRatingPerCritChance,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:    7164,
		stats.Strength:  54,
		stats.Agility:   55,
		stats.Stamina:   77,
		stats.Intellect: 130,
		stats.Spirit:    141,
		stats.Mana:      3856,
		stats.SpellCrit: 1.697 * core.CritRatingPerCritChance,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:    7164,
		stats.Strength:  50,
		stats.Agility:   56,
		stats.Stamina:   76,
		stats.Intellect: 131,
		stats.Spirit:    144,
		stats.Mana:      3856,
		stats.SpellCrit: 1.697 * core.CritRatingPerCritChance,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:    7164,
		stats.Strength:  51,
		stats.Agility:   58,
		stats.Stamina:   76,
		stats.Intellect: 133,
		stats.Spirit:    139, // racial makes this 170
		stats.Mana:      3856,
		stats.SpellCrit: 1.697 * core.CritRatingPerCritChance,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassWarlock}] = stats.Stats{
		stats.Health:    7164,
		stats.Strength:  46,
		stats.Agility:   60,
		stats.Stamina:   76,
		stats.Intellect: 136, // racial makes this 170
		stats.Spirit:    139,
		stats.Mana:      3856,
		stats.SpellCrit: 1.697 * core.CritRatingPerCritChance,
		// Not sure how stats modify the crit chance.
		// stats.MeleeCrit:   4.43 * core.CritRatingPerCritChance,
	}
}

// Agent is a generic way to access underlying warlock on any of the agents.
type WarlockAgent interface {
	GetWarlock() *Warlock
}

func (warlock *Warlock) HasMajorGlyph(glyph proto.WarlockMajorGlyph) bool {
	return warlock.HasGlyph(int32(glyph))
}

func (warlock *Warlock) HasMinorGlyph(glyph proto.WarlockMinorGlyph) bool {
	return warlock.HasGlyph(int32(glyph))
}
