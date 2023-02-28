package warrior

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/proto"
)

func (warrior *Warrior) registerShockwaveSpell() {
	if !warrior.Talents.Shockwave {
		return
	}

	warrior.Shockwave = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 46968},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRanged, // TODO: Is this correct?
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RageCost: core.RageCostOptions{
			Cost:   15 - float64(warrior.Talents.FocusedRage),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: 20*time.Second - core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfShockwave), 3*time.Second, 0),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance)
		},

		CritMultiplier:   warrior.critMultiplier(none),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.75 * spell.MeleeAttackPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.Targets {
				result := spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				// TODO: AOE spells usually don't give refunds, this is probably wrong
				if !result.Landed() {
					spell.IssueRefund(sim)
				}
			}
		},
	})
}
