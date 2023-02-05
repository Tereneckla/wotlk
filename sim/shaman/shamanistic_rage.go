package shaman

import (
	"time"

	"github.com/Tereneckla/wotlk70/sim/core"
	"github.com/Tereneckla/wotlk70/sim/core/stats"
)

func (shaman *Shaman) registerShamanisticRageCD() {
	if !shaman.Talents.ShamanisticRage {
		return
	}

	t10Bonus := false
	if shaman.HasSetBonus(ItemSetFrostWitchBattlegear, 2) {
		t10Bonus = true
	}

	actionID := core.ActionID{SpellID: 30823}
	ppmm := shaman.AutoAttacks.NewPPMManager(15, core.ProcMaskMelee)
	manaMetrics := shaman.NewManaMetrics(actionID)
	srAura := shaman.RegisterAura(core.Aura{
		Label:    "Shamanistic Rage",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.7
			if t10Bonus {
				aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.12
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.7
			if t10Bonus {
				aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.12
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// proc mask: 20
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}
			if !ppmm.Proc(sim, spell.ProcMask, "shamanistic rage") {
				return
			}
			mana := aura.Unit.GetStat(stats.AttackPower) * 0.15
			aura.Unit.AddMana(sim, mana, manaMetrics, true)
		},
	})

	spell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 1,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			srAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			manaReserve := shaman.ShamanisticRageManaThreshold / 100 * shaman.MaxMana()
			return character.CurrentMana() <= manaReserve
		},
	})
}
