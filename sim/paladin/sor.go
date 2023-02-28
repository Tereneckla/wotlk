package paladin

import (
	"github.com/Tereneckla/wotlk/sim/core"
)

func (paladin *Paladin) registerSealOfRighteousnessSpellAndAura() {
	/*
	 * Seal of Righteousness is an Spell/Aura that when active makes the paladin capable of procing
	 * 2 different SpellIDs depending on a paladin's casted spell or melee swing.
	 * NOTE:
	 *   Seal of Righteousness is unique in that it is the only seal that can proc off its own judgements.
	 *
	 * SpellID 20187 (Judgement of Righteousness):
	 *   - Procs off of any "Primary" Judgement (JoL, JoW, JoJ).
	 *   - Cannot miss or be dodged/parried.
	 *   - Deals hybrid AP/SP damage.
	 *   - Crits off of a melee modifier.
	 *
	 * SpellID 20154 (Seal of Righteousness):
	 *   - Procs off of any melee special ability, or white hit.
	 *   - Cannot miss or be dodged/parried.
	 *   - Deals hybrid AP/SP damage * current weapon speed.
	 *   - CANNOT CRIT.
	 */
	justicarBattle2 := core.TernaryFloat64(paladin.HasSetBonus(ItemSetJusticarBattlegear, 2), 33, 0) //damage to judge
	justicarArmor2 := core.TernaryFloat64(paladin.HasSetBonus(ItemSetJusticarArmor, 2), 0.1, 0)      //multiplier seal

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20187}, // Judgement of Righteousness.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeOrRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

		BonusCritRating: (6 * float64(paladin.Talents.Fanaticism) * core.CritRatingPerCritChance),

		DamageMultiplier: 1 *
			(1 + paladin.getTalentSealsOfThePureBonus() +
				paladin.getMajorGlyphOfJudgementBonus() + paladin.getTalentTheArtOfWarBonus()) *
			(1 + paladin.getTalentTwoHandedWeaponSpecializationBonus()) *
			(1 + core.TernaryFloat64(paladin.HasSetBonus(ItemSetJusticarBattlegear, 4), 0.1, 0)),
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// i = 1 + 0.2 * AP + 0.32 * HolP
			baseDamage := 1 +
				.20*spell.MeleeAttackPower() +
				.32*spell.SpellPower() +
				justicarBattle2

			// Secondary Judgements cannot miss if the Primary Judgement hit, only roll for crit.
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20154}, // Seal of Righteousness damage bonus.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1 *
			(1 + paladin.getTalentSealsOfThePureBonus()) *
			(1 + paladin.getMajorGlyphSealOfRighteousnessBonus()) *
			(1 + paladin.getTalentTwoHandedWeaponSpecializationBonus()) *
			(1 + justicarArmor2),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// weapon_speed * (0.022* AP + 0.044*HolP)
			speed := 2.0
			if paladin.GetMHWeapon() != nil {
				speed = paladin.GetMHWeapon().SwingSpeed
			}
			baseDamage := speed * (.022*spell.MeleeAttackPower() + .044*spell.SpellPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})

	// Seal of Righteousness aura.
	auraActionID := core.ActionID{SpellID: 21084}
	paladin.SealOfRighteousnessAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Righteousness",
		Tag:      "Seal",
		ActionID: auraActionID,
		Duration: SealDuration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Don't proc on misses or our own procs.
			if !result.Landed() || spell.SpellID == onJudgementProc.SpellID || spell.SpellID == onSpecialOrSwingProc.SpellID {
				return
			}

			// Differ between judgements and other melee abilities.
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				// SoR is the only seal that can proc off its own judgement.
				onSpecialOrSwingProc.Cast(sim, result.Target)
				onJudgementProc.Cast(sim, result.Target)
				if paladin.Talents.JudgementsOfTheJust > 0 {
					// Special JoJ talent behavior, procs swing seal on judgements
					// Yes, for SoR this means it proces TWICE on one judgement.
					onSpecialOrSwingProc.Cast(sim, result.Target)
				}
			} else {
				if spell.IsMelee() {
					onSpecialOrSwingProc.Cast(sim, result.Target)
				}
			}
		},
	})

	aura := paladin.SealOfRighteousnessAura
	paladin.SealOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID, // Seal of Righteousness self buff.
		SpellSchool: core.SpellSchoolHoly,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.14,
			Multiplier: 1 - 0.02*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
