package mage

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/proto"
	"github.com/Tereneckla/wotlk/sim/core/stats"
)

func (mage *Mage) registerManaGemsCD() {

	actionID := core.ActionID{ItemID: 33312}
	manaMetrics := mage.NewManaMetrics(actionID)
	var gemAura *core.Aura

	var serpentCoilAura *core.Aura
	if mage.HasTrinketEquipped(30720) {
		serpentCoilAura = mage.NewTemporaryStatsAura("Serpent-Coil Braid", core.ActionID{ItemID: 30720}, stats.Stats{stats.SpellPower: 225}, 15*time.Second)
	}

	manaMultiplier := core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfManaGem), 1.4, 1) *
		(1 +
			core.TernaryFloat64(serpentCoilAura != nil, 0.25, 0))

	minManaEmeraldGain := 2340.0 * manaMultiplier
	maxManaEmeraldGain := 2460.0 * manaMultiplier
	//minManaSapphireGain := 3330.0 * manaMultiplier
	//maxManaSapphireGain := 3500.0 * manaMultiplier
	minManaRubyGain := 1073 * manaMultiplier
	maxManaRubyGain := 1127 * manaMultiplier

	var remainingManaGems int
	mage.RegisterResetEffect(func(sim *core.Simulation) {
		remainingManaGems = 6
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return remainingManaGems != 0
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			var manaGain float64
			if remainingManaGems > 3 {
				// Mana Emerald
				manaGain = sim.Roll(minManaEmeraldGain, maxManaEmeraldGain)
			} else {
				// Mana Ruby
				manaGain = sim.Roll(minManaRubyGain, maxManaRubyGain)
			}

			if gemAura != nil {
				gemAura.Activate(sim)
			}
			if serpentCoilAura != nil {
				serpentCoilAura.Activate(sim)
			}

			mage.AddMana(sim, manaGain, manaMetrics)

			remainingManaGems--
			if remainingManaGems == 0 {
				// Disable this cooldown since we're out of emeralds.
				mage.GetMajorCooldown(actionID).Disable()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Only pop if we have less than the max mana provided by the gem minus 1mp5 tick.
			totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
			maxManaGain := maxManaEmeraldGain
			if remainingManaGems <= 3 {
				maxManaGain = maxManaRubyGain
			}
			if character.MaxMana()-(character.CurrentMana()+totalRegen) < maxManaGain {
				return false
			}

			return true
		},
	})
}
