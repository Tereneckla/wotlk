package priest

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/proto"
	"github.com/Tereneckla/wotlk/sim/core/stats"
)

var TalentTreeSizes = [3]int{28, 27, 27}

type Priest struct {
	core.Character
	SelfBuffs
	Talents *proto.PriestTalents

	SurgeOfLight bool

	Latency float64

	ShadowfiendAura *core.Aura
	ShadowfiendPet  *Shadowfiend

	// cached cast stuff
	// TODO: aoe multi-target situations will need multiple spells ticking for each target.
	InnerFocusAura     *core.Aura
	MiseryAura         *core.Aura
	ShadowWeavingAura  *core.Aura
	ShadowyInsightAura *core.Aura
	ImprovedSpiritTap  *core.Aura
	DispersionAura     *core.Aura

	SurgeOfLightProcAura *core.Aura

	BindingHeal     *core.Spell
	CircleOfHealing *core.Spell
	DevouringPlague *core.Spell
	FlashHeal       *core.Spell
	GreaterHeal     *core.Spell
	HolyFire        *core.Spell
	InnerFocus      *core.Spell
	ShadowWordPain  *core.Spell
	MindBlast       *core.Spell
	MindFlay        []*core.Spell
	//MindSear        []*core.Spell
	Penance         *core.Spell
	PenanceHeal     *core.Spell
	PowerWordShield *core.Spell
	PrayerOfHealing *core.Spell
	PrayerOfMending *core.Spell
	Renew           *core.Spell
	EmpoweredRenew  *core.Spell
	ShadowWordDeath *core.Spell
	Shadowfiend     *core.Spell
	Smite           *core.Spell
	Starshards      *core.Spell
	VampiricTouch   *core.Spell
	Dispersion      *core.Spell

	PWSShields    []*core.Shield
	WeakenedSouls core.AuraArray

	ProcPrayerOfMending core.ApplySpellResults

	DpInitMultiplier float64
}

type SelfBuffs struct {
	UseShadowfiend bool
	UseInnerFire   bool

	PowerInfusionTarget *proto.RaidTarget
}

func (priest *Priest) GetCharacter() *core.Character {
	return &priest.Character
}

func (priest *Priest) HasMajorGlyph(glyph proto.PriestMajorGlyph) bool {
	return priest.HasGlyph(int32(glyph))
}
func (priest *Priest) HasMinorGlyph(glyph proto.PriestMinorGlyph) bool {
	return priest.HasGlyph(int32(glyph))
}

func (priest *Priest) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ShadowProtection = true
	raidBuffs.DivineSpirit = true

	raidBuffs.PowerWordFortitude = core.MaxTristate(raidBuffs.PowerWordFortitude, core.MakeTristateValue(
		true,
		priest.Talents.ImprovedPowerWordFortitude == 2))
}

func (priest *Priest) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (priest *Priest) Initialize() {
	// Shadow Insight gained from Glyph of Shadow
	// Finalized spirit off gear and not dynamic spirit (e.g. Spirit Tap does not increase this)
	priest.ShadowyInsightAura = priest.NewTemporaryStatsAura(
		"Shadowy Insight",
		core.ActionID{SpellID: 61792},
		stats.Stats{stats.SpellPower: priest.GetStat(stats.Spirit) * 0.30},
		time.Second*10,
	)

	priest.registerDevouringPlagueSpell()
	priest.registerShadowWordPainSpell()
	priest.registerMindBlastSpell()
	priest.registerShadowWordDeathSpell()
	priest.registerShadowfiendSpell()
	priest.registerStarshardsSpell()
	priest.registerVampiricTouchSpell()
	priest.registerDispersionSpell()

	priest.registerPowerInfusionCD()

	priest.MindFlay = []*core.Spell{
		nil, // So we can use # of ticks as the index
		priest.newMindFlaySpell(1),
		priest.newMindFlaySpell(2),
		priest.newMindFlaySpell(3),
	}
	/*priest.MindSear = []*core.Spell{
		nil, // So we can use # of ticks as the index
		priest.newMindSearSpell(1),
		priest.newMindSearSpell(2),
		priest.newMindSearSpell(3),
		priest.newMindSearSpell(4),
		priest.newMindSearSpell(5),
	}*/
}

func (priest *Priest) RegisterHealingSpells() {
	priest.registerPenanceHealSpell()
	priest.registerBindingHealSpell()
	priest.registerCircleOfHealingSpell()
	priest.registerFlashHealSpell()
	priest.registerGreaterHealSpell()
	priest.registerPowerWordShieldSpell()
	priest.registerPrayerOfHealingSpell()
	priest.registerPrayerOfMendingSpell()
	priest.registerRenewSpell()
}

func (priest *Priest) AddShadowWeavingStack(sim *core.Simulation) {
	if priest.ShadowWeavingAura != nil {
		priest.ShadowWeavingAura.Activate(sim)
		priest.ShadowWeavingAura.AddStack(sim)
	}
}

func (priest *Priest) Reset(_ *core.Simulation) {
}

func New(char core.Character, selfBuffs SelfBuffs, talents string) *Priest {
	priest := &Priest{
		Character: char,
		SelfBuffs: selfBuffs,
		Talents:   &proto.PriestTalents{},
	}
	core.FillTalentsProto(priest.Talents.ProtoReflect(), talents, TalentTreeSizes)
	priest.AddStatDependency(stats.Intellect, stats.SpellCrit, (1/80)*core.CritRatingPerCritChance)
	priest.EnableManaBar()
	priest.ShadowfiendPet = priest.NewShadowfiend()

	if selfBuffs.UseInnerFire {
		multi := 1 + float64(priest.Talents.ImprovedInnerFire)*0.15
		sp := 120.0 * multi
		armor := 2440 * multi * core.TernaryFloat64(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfInnerFire), 1.5, 1)
		priest.AddStat(stats.SpellPower, sp)
		priest.AddStat(stats.Armor, armor)
	}

	return priest
}

func init() {
	//const basecrit = 3.29 * core.CritRatingPerCritChance
	const basespellcrit = 1.24 * core.CritRatingPerCritChance
	const basehealth = 3391
	const basemana = 2953

	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    basehealth,
		stats.Mana:      basemana,
		stats.Strength:  39,
		stats.Agility:   45,
		stats.Stamina:   58,
		stats.Intellect: 145,
		stats.Spirit:    151,
		stats.SpellCrit: basespellcrit,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    basehealth,
		stats.Mana:      basemana,
		stats.Strength:  44,
		stats.Agility:   41,
		stats.Stamina:   59,
		stats.Intellect: 144,
		stats.Spirit:    150,
		stats.SpellCrit: basespellcrit,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    basehealth,
		stats.Mana:      basemana,
		stats.Strength:  35,
		stats.Agility:   49,
		stats.Stamina:   58,
		stats.Intellect: 145,
		stats.Spirit:    151,
		stats.SpellCrit: basespellcrit,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    basehealth,
		stats.Mana:      basemana,
		stats.Strength:  40,
		stats.Agility:   42,
		stats.Stamina:   58,
		stats.Intellect: 145,
		stats.Spirit:    153,
		stats.SpellCrit: basespellcrit,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    basehealth,
		stats.Mana:      basemana,
		stats.Strength:  38,
		stats.Agility:   43,
		stats.Stamina:   58,
		stats.Intellect: 143,
		stats.Spirit:    156,
		stats.SpellCrit: basespellcrit,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    basehealth,
		stats.Mana:      basemana,
		stats.Strength:  40,
		stats.Agility:   47,
		stats.Stamina:   58,
		stats.Intellect: 141,
		stats.Spirit:    152,
		stats.SpellCrit: basespellcrit,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassPriest}] = stats.Stats{
		stats.Health:    basehealth,
		stats.Mana:      basemana,
		stats.Strength:  36,
		stats.Agility:   47,
		stats.Stamina:   58,
		stats.Intellect: 148,
		stats.Spirit:    149,
		stats.SpellCrit: basespellcrit,
	}
}

// Agent is a generic way to access underlying priest on any of the agents.
type PriestAgent interface {
	GetPriest() *Priest
}
