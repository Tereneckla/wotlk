package druid

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/proto"
)

func (druid *Druid) registerMoonfireSpell() {
	numTicks := druid.moonfireTicks()

	bonusSpellpower := core.TernaryFloat64(druid.Equip[proto.ItemSlot_ItemSlotRanged].ID == 23197, 33, 0)

	starfireBonusCrit := float64(druid.Talents.ImprovedInsectSwarm) * core.CritRatingPerCritChance

	baseDamageMultiplier := 1 +
		0.05*float64(druid.Talents.ImprovedMoonfire) +
		[]float64{0.0, 0.03, 0.06, 0.1}[druid.Talents.Moonfury]

	malusInitialDamageMultiplier := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMoonfire), 0.9, 0)

	bonusPeriodicDamageMultiplier := 0 +
		0.01*float64(druid.Talents.Genesis) +
		core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMoonfire), 0.75, 0)

	druid.Moonfire = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48463},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagNaturesGrace | SpellFlagOmenTrigger,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.21,
			Multiplier: 1 - 0.03*float64(druid.Talents.Moonglow),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		BonusSpellPower:  bonusSpellpower,
		BonusCritRating:  float64(druid.Talents.ImprovedMoonfire) * 5 * core.CritRatingPerCritChance,
		DamageMultiplier: baseDamageMultiplier - malusInitialDamageMultiplier,

		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Moonfire",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					druid.Starfire.BonusCritRating += starfireBonusCrit
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					druid.Starfire.BonusCritRating -= starfireBonusCrit
				},
			},
			NumberOfTicks: druid.moonfireTicks(),
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Spell.DamageMultiplier = baseDamageMultiplier + bonusPeriodicDamageMultiplier
				dot.SnapshotBaseDamage = 150 + 0.13*dot.Spell.SpellPower()
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				dot.Spell.DamageMultiplier = baseDamageMultiplier - malusInitialDamageMultiplier
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(305, 357) + 0.15*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.NumberOfTicks = numTicks
				dot.Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	})
}

func (druid *Druid) moonfireTicks() int32 {
	return 4 +
		core.TernaryInt32(druid.Talents.NaturesSplendor, 1, 0) +
		core.TernaryInt32(druid.HasSetBonus(ItemSetThunderheartRegalia, 2), 1, 0)
}
