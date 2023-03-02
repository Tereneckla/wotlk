package rogue

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/proto"
)

const RuptureEnergyCost = 25.0
const RuptureSpellID = 26867

func (rogue *Rogue) registerRupture() {
	glyphTicks := core.TernaryInt32(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfRupture), 2, 0)
	rogue.Rupture = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: RuptureSpellID},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | rogue.finisherFlags(),
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:          RuptureEnergyCost,
			Refund:        0.4 * float64(rogue.Talents.QuickRecovery),
			RefundMetrics: rogue.QuickRecoveryMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		DamageMultiplier: 1 +
			0.15*float64(rogue.Talents.BloodSpatter) +
			0.02*float64(rogue.Talents.FindWeakness) +
			0.1*float64(rogue.Talents.SerratedBlades),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rupture",
				Tag:   RogueBleedTag,
			},
			NumberOfTicks: 0, // Set dynamically
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				comboPoints := rogue.ComboPoints()
				dot.SnapshotBaseDamage = 127 +
					18*float64(comboPoints) +
					[]float64{0, 0.06 / 4, 0.12 / 5, 0.18 / 6, 0.24 / 7, 0.30 / 8}[comboPoints]*dot.Spell.MeleeAttackPower()

				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(target, attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				comboPoints := rogue.ComboPoints()
				dot := spell.Dot(target)
				dot.Spell = spell
				dot.NumberOfTicks = 3 + comboPoints + glyphTicks
				dot.RecomputeAuraDuration()
				dot.Apply(sim)
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

func (rogue *Rogue) RuptureDuration(comboPoints int32) time.Duration {
	return time.Second*6 +
		time.Second*2*time.Duration(comboPoints) +
		core.TernaryDuration(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfRupture), time.Second*4, 0)
}
