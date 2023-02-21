package hunter

import (
	"math"
	"time"

	"github.com/Tereneckla/wotlk70/sim/core"
	"github.com/Tereneckla/wotlk70/sim/core/stats"
)

type PetAbilityType int

// Pet AI doesn't use abilities immediately, so model this with a 1.6s GCD.
const PetGCD = time.Millisecond * 1600

const (
	Unknown PetAbilityType = iota
	AcidSpit
	Bite
	Claw
	DemoralizingScreech
	FireBreath
	FuriousHowl
	FroststormBreath
	Gore
	LavaBreath
	LightningBreath
	MonstrousBite
	NetherShock
	Pin
	PoisonSpit
	Rake
	Ravage
	SavageRend
	ScorpidPoison
	Smack
	Snatch
	SonicBlast
	SpiritStrike
	SporeCloud
	Stampede
	Sting
	Swipe
	TendonRip
	VenomWebSpray
)

// These IDs are needed for certain talents.
const BiteSpellID = 52474
const ClawSpellID = 52472
const SmackSpellID = 52476

type PetAbility struct {
	Type PetAbilityType

	// Focus cost
	Cost float64

	*core.Spell
}

func (ability *PetAbility) IsEmpty() bool {
	return ability.Spell == nil
}

// Returns whether the ability was successfully cast.
func (ability *PetAbility) TryCast(sim *core.Simulation, target *core.Unit, hp *HunterPet) bool {
	if ability.IsEmpty() {
		return false
	}
	if hp.currentFocus < ability.Cost {
		return false
	}
	if !ability.IsReady(sim) {
		return false
	}

	hp.SpendFocus(sim, ability.Cost*hp.PseudoStats.CostMultiplier, ability.ActionID)
	ability.Cast(sim, target)
	return true
}

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) PetAbility {
	switch abilityType {
	case AcidSpit:
		return hp.newAcidSpit()
	case Bite:
		return hp.newBite()
	case Claw:
		return hp.newClaw()
	case DemoralizingScreech:
		return hp.newDemoralizingScreech()
	case FireBreath:
		return hp.newFireBreath()
	case FroststormBreath:
		return hp.newFroststormBreath()
	case FuriousHowl:
		return hp.newFuriousHowl()
	case Gore:
		return hp.newGore()
	case LavaBreath:
		return hp.newLavaBreath()
	case LightningBreath:
		return hp.newLightningBreath()
	case MonstrousBite:
		return hp.newMonstrousBite()
	case NetherShock:
		return hp.newNetherShock()
	case Pin:
		return hp.newPin()
	case PoisonSpit:
		return hp.newPoisonSpit()
	case Rake:
		return hp.newRake()
	case Ravage:
		return hp.newRavage()
	case SavageRend:
		return hp.newSavageRend()
	case ScorpidPoison:
		return hp.newScorpidPoison()
	case Smack:
		return hp.newSmack()
	case Snatch:
		return hp.newSnatch()
	case SonicBlast:
		return hp.newSonicBlast()
	case SpiritStrike:
		return hp.newSpiritStrike()
	case SporeCloud:
		return hp.newSporeCloud()
	case Stampede:
		return hp.newStampede()
	case Sting:
		return hp.newSting()
	case Swipe:
		return hp.newSwipe()
	case TendonRip:
		return hp.newTendonRip()
	case VenomWebSpray:
		return hp.newVenomWebSpray()
	case Unknown:
		return PetAbility{}
	default:
		panic("Invalid pet ability type")
	}
}

func (hp *HunterPet) newFocusDump(pat PetAbilityType, spellID int32) PetAbility {
	return PetAbility{
		Type: pat,
		Cost: 25,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellID},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: PetGCD,
				},
				IgnoreHaste: true,
			},

			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			CritMultiplier:   2,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(118, 168) + 0.07*spell.MeleeAttackPower()
				baseDamage *= hp.killCommandMult()
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			},
		}),
	}
}

func (hp *HunterPet) newBite() PetAbility {
	return hp.newFocusDump(Bite, BiteSpellID)
}
func (hp *HunterPet) newClaw() PetAbility {
	return hp.newFocusDump(Claw, ClawSpellID)
}
func (hp *HunterPet) newSmack() PetAbility {
	return hp.newFocusDump(Smack, SmackSpellID)
}

type PetSpecialAbilityConfig struct {
	Type    PetAbilityType
	Cost    float64
	SpellID int32
	School  core.SpellSchool
	GCD     time.Duration
	CD      time.Duration
	MinDmg  float64
	MaxDmg  float64
	APRatio float64

	Dot core.DotConfig

	OnSpellHitDealt func(*core.Simulation, *core.Spell, *core.SpellResult)
}

func (hp *HunterPet) newSpecialAbility(config PetSpecialAbilityConfig) PetAbility {
	var flags core.SpellFlag
	var applyEffects core.ApplySpellResults
	var procMask core.ProcMask
	onSpellHitDealt := config.OnSpellHitDealt
	if config.School == core.SpellSchoolPhysical {
		flags = core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage
		procMask = core.ProcMaskSpellDamage
		applyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(config.MinDmg, config.MaxDmg) + config.APRatio*spell.MeleeAttackPower()
			baseDamage *= hp.killCommandMult()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			if onSpellHitDealt != nil {
				onSpellHitDealt(sim, spell, result)
			}
		}
	} else {
		procMask = core.ProcMaskMeleeMHSpecial
		applyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(config.MinDmg, config.MaxDmg) + config.APRatio*spell.MeleeAttackPower()
			baseDamage *= 1 + 0.2*float64(hp.KillCommandAura.GetStacks())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if onSpellHitDealt != nil {
				onSpellHitDealt(sim, spell, result)
			}
		}
	}

	return PetAbility{
		Type: config.Type,
		Cost: config.Cost,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: config.SpellID},
			SpellSchool: config.School,
			ProcMask:    procMask,
			Flags:       flags,

			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			CritMultiplier:   2,
			ThreatMultiplier: 1,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: config.GCD,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(config.CD),
				},
			},
			Dot:          config.Dot,
			ApplyEffects: applyEffects,
		}),
	}
}

func (hp *HunterPet) newAcidSpit() PetAbility {
	acidSpitAura := core.AcidSpitAura(hp.CurrentTarget)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    AcidSpit,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 55754,
		School:  core.SpellSchoolNature,
		MinDmg:  58,
		MaxDmg:  82,
		APRatio: 0.049,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				acidSpitAura.Activate(sim)
				if acidSpitAura.IsActive() {
					acidSpitAura.AddStack(sim)
				}
			}
		},
	})
}

func (hp *HunterPet) newDemoralizingScreech() PetAbility {
	var debuffs []*core.Aura
	for _, target := range hp.Env.Encounter.Targets {
		debuffs = append(debuffs, core.DemoralizingScreechAura(&target.Unit))
	}

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    DemoralizingScreech,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 55487,
		School:  core.SpellSchoolPhysical,
		MinDmg:  42,
		MaxDmg:  58,
		APRatio: 0.07,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				for _, debuff := range debuffs {
					debuff.Activate(sim)
				}
			}
		},
	})
}

func (hp *HunterPet) newFireBreath() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    FireBreath,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 55485,
		School:  core.SpellSchoolFire,
		MinDmg:  20,
		MaxDmg:  26,
		APRatio: 0.049,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Fire Breath",
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(22/2, 26/2) * hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newFroststormBreath() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    FroststormBreath,
		Cost:    20,
		GCD:     0,
		CD:      time.Second * 10,
		SpellID: 55492,
		School:  core.SpellSchoolFrost,
		MinDmg:  59,
		MaxDmg:  81,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newFuriousHowl() PetAbility {
	actionID := core.ActionID{SpellID: 64495}

	petAura := hp.NewTemporaryStatsAura("FuriousHowl", actionID, stats.Stats{stats.AttackPower: 204, stats.RangedAttackPower: 204}, time.Second*20)
	ownerAura := hp.hunterOwner.NewTemporaryStatsAura("FuriousHowl", actionID, stats.Stats{stats.AttackPower: 204, stats.RangedAttackPower: 204}, time.Second*20)
	const cost = 20.0

	howlSpell := hp.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hp.SpendFocus(sim, cost, actionID)
			petAura.Activate(sim)
			ownerAura.Activate(sim)
		},
	})

	hp.hunterOwner.AddMajorCooldown(core.MajorCooldown{
		Spell: howlSpell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return hp.IsEnabled() && hp.CurrentFocus() >= cost
		},
	})

	return PetAbility{}
}

func (hp *HunterPet) newGore() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Gore,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 35295,
		School:  core.SpellSchoolPhysical,
		MinDmg:  57,
		MaxDmg:  75,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newLavaBreath() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    LavaBreath,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 58611,
		School:  core.SpellSchoolFire,
		MinDmg:  60,
		MaxDmg:  80,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newLightningBreath() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    LightningBreath,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 25012,
		School:  core.SpellSchoolNature,
		MinDmg:  40,
		MaxDmg:  52,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newMonstrousBite() PetAbility {
	procAura := hp.RegisterAura(core.Aura{
		Label:     "Monstrous Bite",
		ActionID:  core.ActionID{SpellID: 55499},
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= math.Pow(1.03, float64(oldStacks))
			aura.Unit.PseudoStats.DamageDealtMultiplier *= math.Pow(1.03, float64(newStacks))
		},
	})

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    MonstrousBite,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 55499,
		School:  core.SpellSchoolPhysical,
		MinDmg:  43,
		MaxDmg:  57,
		APRatio: 0.07,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			}
		},
	})
}

func (hp *HunterPet) newNetherShock() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    NetherShock,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 53589,
		School:  core.SpellSchoolShadow,
		MinDmg:  30,
		MaxDmg:  40,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newPin() PetAbility {
	return PetAbility{
		Type: Pin,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 53548},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD:         PetGCD,
					ChannelTime: time.Second * 4,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
				},
			},

			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Pin",
				},
				NumberOfTicks: 4,
				TickLength:    time.Second * 1,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = sim.Roll(56/4, 72/4) + 0.07*dot.Spell.MeleeAttackPower()
					dot.SnapshotBaseDamage *= hp.killCommandMult()
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
				if result.Landed() {
					spell.Dot(result.Target).Apply(sim)
				}
			},
		}),
	}
}

func (hp *HunterPet) newPoisonSpit() PetAbility {
	return PetAbility{
		Type: PoisonSpit,
		Cost: 20,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 55557},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: PetGCD,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(time.Second * 10),
				},
			},

			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "PoisonSpit",
				},
				NumberOfTicks: 4,
				TickLength:    time.Second * 2,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = sim.Roll(48/4, 64/4) + (0.049/4)*dot.Spell.MeleeAttackPower()
					dot.SnapshotBaseDamage *= hp.killCommandMult()
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
				if result.Landed() {
					spell.Dot(result.Target).Apply(sim)
				}
			},
		}),
	}
}

func (hp *HunterPet) newRake() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Rake,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 59886,
		School:  core.SpellSchoolPhysical,
		MinDmg:  22,
		MaxDmg:  30,
		APRatio: 0.0175,
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rake",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(7, 13) + 0.0175*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newRavage() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Ravage,
		Cost:    0,
		CD:      time.Second * 40,
		SpellID: 53562,
		School:  core.SpellSchoolPhysical,
		MinDmg:  50,
		MaxDmg:  70,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newSavageRend() PetAbility {
	actionID := core.ActionID{SpellID: 53582}
	const cost = 20.0

	procAura := hp.RegisterAura(core.Aura{
		Label:    "Savage Rend",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
		},
	})

	srSpell := hp.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagApplyArmorReduction,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 60),
			},
		},

		DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SavageRend",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 5,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(10, 12) + 0.07*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(28, 38) + 0.07*spell.MeleeAttackPower()
			baseDamage *= hp.killCommandMult()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			hp.SpendFocus(sim, cost, actionID)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				if result.DidCrit() {
					procAura.Activate(sim)
				}
			}
		},
	})

	hp.hunterOwner.AddMajorCooldown(core.MajorCooldown{
		Spell: srSpell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return hp.IsEnabled() && hp.CurrentFocus() >= cost
		},
	})

	return PetAbility{}
}

func (hp *HunterPet) newScorpidPoison() PetAbility {
	return PetAbility{
		Type: ScorpidPoison,
		Cost: 20,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 55728},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: PetGCD,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(time.Second * 10),
				},
			},

			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "ScorpidPoison",
				},
				NumberOfTicks: 5,
				TickLength:    time.Second * 2,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = sim.Roll(35/5, 65/5) + (0.07/5)*dot.Spell.MeleeAttackPower()
					dot.SnapshotBaseDamage *= hp.killCommandMult()
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
		}),
	}
}

func (hp *HunterPet) newSnatch() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Snatch,
		Cost:    20,
		CD:      time.Second * 60,
		SpellID: 53543,
		School:  core.SpellSchoolPhysical,
		MinDmg:  42,
		MaxDmg:  58,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newSonicBlast() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    SonicBlast,
		Cost:    80,
		CD:      time.Second * 60,
		SpellID: 53568,
		School:  core.SpellSchoolNature,
		MinDmg:  29,
		MaxDmg:  41,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newSpiritStrike() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    SpiritStrike,
		Cost:    20,
		GCD:     0,
		CD:      time.Second * 10,
		SpellID: 61198,
		School:  core.SpellSchoolArcane,
		MinDmg:  23,
		MaxDmg:  29,
		APRatio: 0.04,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SpiritStrike",
			},
			NumberOfTicks: 1,
			TickLength:    time.Second * 6,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(23, 29) + 0.04*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newSporeCloud() PetAbility {
	var debuffs []*core.Aura
	for _, target := range hp.Env.Encounter.Targets {
		debuffs = append(debuffs, core.SporeCloudAura(&target.Unit))
	}

	return PetAbility{
		Type: SporeCloud,
		Cost: 20,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 53598},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellDamage,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: PetGCD,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(time.Second * 10),
				},
			},

			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				IsAOE: true,
				Aura: core.Aura{
					Label: "SporeCloud",
				},
				NumberOfTicks: 3,
				TickLength:    time.Second * 3,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = sim.Roll(11, 13) + (0.049/3)*dot.Spell.MeleeAttackPower()
					dot.SnapshotBaseDamage *= hp.killCommandMult()
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
					for _, aoeTarget := range sim.Encounter.Targets {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, &aoeTarget.Unit, dot.OutcomeTick)
					}
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				spell.AOEDot().Apply(sim)
				for _, debuff := range debuffs {
					debuff.Activate(sim)
				}
			},
		}),
	}
}

func (hp *HunterPet) newStampede() PetAbility {
	debuff := core.StampedeAura(hp.CurrentTarget)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Stampede,
		Cost:    0,
		CD:      time.Second * 60,
		SpellID: 57393,
		School:  core.SpellSchoolPhysical,
		MinDmg:  85,
		MaxDmg:  115,
		APRatio: 0.07,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				debuff.Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) newSting() PetAbility {
	debuff := core.StingAura(hp.CurrentTarget)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Sting,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 6,
		SpellID: 56631,
		School:  core.SpellSchoolNature,
		MinDmg:  30,
		MaxDmg:  40,
		APRatio: 0.049,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				debuff.Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) newSwipe() PetAbility {
	// TODO: This is frontal cone, but might be more realistic as single-target
	// since pets are hard to control.
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Swipe,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 5,
		SpellID: 53533,
		School:  core.SpellSchoolPhysical,
		MinDmg:  42,
		MaxDmg:  58,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newTendonRip() PetAbility {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    TendonRip,
		Cost:    20,
		CD:      time.Second * 20,
		SpellID: 53575,
		School:  core.SpellSchoolPhysical,
		MinDmg:  33,
		MaxDmg:  45,
		APRatio: 0,
	})
}

func (hp *HunterPet) newVenomWebSpray() PetAbility {
	return PetAbility{
		Type: VenomWebSpray,
		Cost: 0,

		Spell: hp.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 55509},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    hp.NewTimer(),
					Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
				},
			},

			DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "VenomWebSpray",
				},
				NumberOfTicks: 4,
				TickLength:    time.Second * 1,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = 39 + 0.07*dot.Spell.MeleeAttackPower()
					dot.SnapshotBaseDamage *= hp.killCommandMult()
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
		}),
	}
}
