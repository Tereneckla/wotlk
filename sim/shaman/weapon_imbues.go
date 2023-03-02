package shaman

import (
	"time"

	"github.com/Tereneckla/wotlk/sim/core"
	"github.com/Tereneckla/wotlk/sim/core/proto"
	"github.com/Tereneckla/wotlk/sim/core/stats"
)

var TotemOfTheAstralWinds int32 = 27815

func (shaman *Shaman) RegisterOnItemSwapWithImbue(effectID int32, procMask *core.ProcMask, aura *core.Aura) {
	shaman.RegisterOnItemSwap(func(sim *core.Simulation) {
		mh := shaman.Equip[proto.ItemSlot_ItemSlotMainHand].TempEnchant == effectID
		oh := shaman.Equip[proto.ItemSlot_ItemSlotOffHand].TempEnchant == effectID
		*procMask = core.GetMeleeProcMaskForHands(mh, oh)

		if !mh && !oh {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}

func (shaman *Shaman) newWindfuryImbueSpell(isMH bool) *core.Spell {
	apBonus := 445.0
	if shaman.Equip[proto.ItemSlot_ItemSlotRanged].ID == TotemOfTheAstralWinds {
		apBonus += 80
	}

	spellConfig := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 25505}.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMelee,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		DamageMultiplier: []float64{1, 1.13, 1.27, 1.4}[shaman.Talents.ElementalWeapons],
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := spell.BonusWeaponDamage()
			mAP := spell.MeleeAttackPower() + apBonus

			if isMH {
				baseDamage1 := constBaseDamage + spell.Unit.MHWeaponDamage(sim, mAP)
				baseDamage2 := constBaseDamage + spell.Unit.MHWeaponDamage(sim, mAP)
				result1 := spell.CalcDamage(sim, target, baseDamage1, spell.OutcomeMeleeSpecialHitAndCrit)
				result2 := spell.CalcDamage(sim, target, baseDamage2, spell.OutcomeMeleeSpecialHitAndCrit)
				spell.DealDamage(sim, result1)
				spell.DealDamage(sim, result2)
			} else {
				baseDamage1 := constBaseDamage + spell.Unit.OHWeaponDamage(sim, mAP)
				baseDamage2 := constBaseDamage + spell.Unit.OHWeaponDamage(sim, mAP)
				result1 := spell.CalcDamage(sim, target, baseDamage1, spell.OutcomeMeleeSpecialHitAndCrit)
				result2 := spell.CalcDamage(sim, target, baseDamage2, spell.OutcomeMeleeSpecialHitAndCrit)
				spell.DealDamage(sim, result1)
				spell.DealDamage(sim, result2)
			}
		},
	}

	return shaman.RegisterSpell(spellConfig)
}

func (shaman *Shaman) RegisterWindfuryImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	var proc = 0.2
	if mh && oh {
		proc = 0.36
	}
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfWindfuryWeapon) {
		proc += 0.02 //TODO: confirm how this actually works
	}

	mhSpell := shaman.newWindfuryImbueSpell(true)
	ohSpell := shaman.newWindfuryImbueSpell(false)

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Second * 3,
	}

	if mh {
		shaman.Equip[proto.ItemSlot_ItemSlotMainHand].TempEnchant = 2636
	}

	if oh {
		shaman.Equip[proto.ItemSlot_ItemSlotOffHand].TempEnchant = 2636
	}

	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	aura := shaman.RegisterAura(core.Aura{
		Label:    "Windfury Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			isMHHit := spell.IsMH()
			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Windfury Imbue") > proc {
				return
			}
			icd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(2636, &procMask, aura)
}

func (shaman *Shaman) newFlametongueImbueSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 25489},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weapon := core.Ternary(isMH, shaman.GetMHWeapon(), shaman.GetOHWeapon())

			var damage float64 = 0
			if weapon != nil {
				baseDamage := weapon.SwingSpeed * 31.07
				spellCoeff := 0.1 * weapon.SwingSpeed / 2.6
				damage = baseDamage + spellCoeff*spell.SpellPower()
			}

			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) ApplyFlametongueImbueToItem(item *core.Item, isDownranked bool) {
	if item == nil || item.TempEnchant == 2634 || item.TempEnchant == 1666 {
		return
	}

	spBonus := 96.0 + core.TernaryFloat64(shaman.HasSetBonus(ItemSetCycloneRegalia, 2), 20, 0)
	spMod := 1.0 + 0.1*float64(shaman.Talents.ElementalWeapons)
	id := 2634
	if isDownranked {
		spBonus = 45.0 + core.TernaryFloat64(shaman.HasSetBonus(ItemSetCycloneRegalia, 2), 20, 0)
		id = 1666
	}

	newStats := stats.Stats{stats.SpellPower: spBonus * spMod}
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFlametongueWeapon) {
		newStats = newStats.Add(stats.Stats{stats.SpellCrit: 2 * core.CritRatingPerCritChance})
	}

	item.Stats = item.Stats.Add(newStats)
	item.TempEnchant = int32(id)
}

func (shaman *Shaman) RegisterFlametongueImbue(mh bool, oh bool) {
	if !mh && !oh && !shaman.ItemSwap.IsEnabled() {
		return
	}

	if mh {
		shaman.ApplyFlametongueImbueToItem(shaman.GetMHWeapon(), false)
	}
	if oh {
		shaman.ApplyFlametongueImbueToItem(shaman.GetOHWeapon(), false)
	}

	ftIcd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond,
	}

	mhSpell := shaman.newFlametongueImbueSpell(true)
	ohSpell := shaman.newFlametongueImbueSpell(false)

	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	aura := shaman.RegisterAura(core.Aura{
		Label:    "Flametongue Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			isMHHit := spell.IsMH()
			if !ftIcd.IsReady(sim) {
				return
			}
			ftIcd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(2634, &procMask, aura)
}

func (shaman *Shaman) newFlametongueDownrankImbueSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 16342},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weapon := core.Ternary(isMH, shaman.GetMHWeapon(), shaman.GetOHWeapon())

			var damage float64 = 0
			if weapon != nil {
				baseDamage := weapon.SwingSpeed * 24.98
				spellCoeff := 0.1 * weapon.SwingSpeed / 2.6
				damage = baseDamage + spellCoeff*spell.SpellPower()
			}

			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) RegisterFlametongueDownrankImbue(mh bool, oh bool) {
	if !mh && !oh && !shaman.ItemSwap.IsEnabled() {
		return
	}

	if mh {
		shaman.ApplyFlametongueImbueToItem(shaman.GetMHWeapon(), true)
	}
	if oh {
		shaman.ApplyFlametongueImbueToItem(shaman.GetOHWeapon(), true)
	}

	ftDownrankIcd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond,
	}

	mhSpell := shaman.newFlametongueDownrankImbueSpell(true)
	ohSpell := shaman.newFlametongueDownrankImbueSpell(false)
	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	aura := shaman.RegisterAura(core.Aura{
		Label:    "Flametongue Imbue (downranked)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			isMHHit := spell.IsMH()
			if !ftDownrankIcd.IsReady(sim) {
				return
			}
			ftDownrankIcd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(1666, &procMask, aura)
}

func (shaman *Shaman) FrostbrandDebuffAura(target *core.Unit) *core.Aura {
	multiplier := 1 + 0.05*float64(shaman.Talents.FrozenPower)
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Frostbrand Attack-" + shaman.Label,
		ActionID: core.ActionID{SpellID: 25501},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.LightningBolt.DamageMultiplier *= multiplier
			shaman.ChainLightning.DamageMultiplier *= multiplier
			shaman.LavaLash.DamageMultiplier *= multiplier
			shaman.EarthShock.DamageMultiplier *= multiplier
			shaman.FlameShock.DamageMultiplier *= multiplier
			shaman.FrostShock.DamageMultiplier *= multiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.LightningBolt.DamageMultiplier /= multiplier
			shaman.ChainLightning.DamageMultiplier /= multiplier
			shaman.LavaLash.DamageMultiplier /= multiplier
			shaman.EarthShock.DamageMultiplier /= multiplier
			shaman.FlameShock.DamageMultiplier /= multiplier
			shaman.FrostShock.DamageMultiplier /= multiplier
		},
	})
}

func (shaman *Shaman) newFrostbrandImbueSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 25500},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskEmpty,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 259 + 0.1*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) RegisterFrostbrandImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	mhSpell := shaman.newFrostbrandImbueSpell(true)
	ohSpell := shaman.newFrostbrandImbueSpell(false)
	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	ppmm := shaman.AutoAttacks.NewPPMManager(9.0, procMask)

	if mh {
		shaman.Equip[proto.ItemSlot_ItemSlotMainHand].TempEnchant = 2635
	} else {
		shaman.Equip[proto.ItemSlot_ItemSlotOffHand].TempEnchant = 2635
	}

	fbDebuffAuras := shaman.NewEnemyAuraArray(shaman.FrostbrandDebuffAura)

	aura := shaman.RegisterAura(core.Aura{
		Label:    "Frostbrand Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			if !ppmm.Proc(sim, spell.ProcMask, "Frostbrand Weapon") {
				return
			}

			if spell.IsMH() {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
			fbDebuffAuras.Get(result.Target).Activate(sim)
		},
	})

	shaman.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(2635, 9.0, &ppmm, aura)
}

//earthliving? not important for dps sims though
