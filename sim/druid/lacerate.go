package druid

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
)

func (druid *Druid) registerLacerateSpell() {
	tickDamage := 155.0 / 5

	if druid.HasSetBonus(ItemSetNordrassilHarness, 4) {
		tickDamage += 15
	}
	initialDamage := 88.0
	if druid.Equip[core.ItemSlotRanged].ID == 27744 { // Idol of Ursoc
		tickDamage += 8
		initialDamage += 8
	}

	druid.Lacerate = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 33745},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		RageCost: core.RageCostOptions{
			Cost:   15 - float64(druid.Talents.ShreddingAttacks),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		DamageMultiplier: 1,
		CritMultiplier:   druid.MeleeCritMultiplier(Bear),
		ThreatMultiplier: 0.5,
		// FlatThreatBonus:  515.5, // Handled below

		Dot: core.DotConfig{
			Aura: druid.applyRendAndTear(core.Aura{
				Label:     "Lacerate",
				MaxStacks: 5,
				Duration:  time.Second * 15,
			}),
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = tickDamage + 0.01*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= float64(dot.Aura.GetStacks())

				if !isRollover {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.Spell.DamageMultiplier = 1
					dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(target, attackTable)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if druid.Talents.PrimalGore {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := initialDamage + 0.01*spell.MeleeAttackPower()
			if druid.BleedCategories.Get(target).AnyActive() {
				baseDamage *= 1.3
			}

			// Hack so that FlatThreatBonus only applies to the initial portion.
			spell.FlatThreatBonus = 515.5
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.FlatThreatBonus = 0
			spell.DamageMultiplier = 1

			if result.Landed() {
				dot := spell.Dot(target)
				if dot.IsActive() {
					dot.Refresh(sim)
					dot.AddStack(sim)
					dot.TakeSnapshot(sim, true)
				} else {
					dot.Activate(sim)
					dot.SetStacks(sim, 1)
					dot.TakeSnapshot(sim, true)
				}
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
