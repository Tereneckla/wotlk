package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) RotationActionCallback_Auto(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {

	if !dk.GCD.IsReady(sim) {
		return dk.NextGCDAt()
	}

	// 1. If you have desolation - Maintain Desolation
	// 2. Maintain Frost Fever
	// 3. Maintain Blood Plague
	// 4. If you have reaping - Spend Blood Runes (Blood Strike)
	// 5. Use extra runes on Scourge Strike
	// 6. If you don't have reaping, spend extra blood runes (Blood Strike)
	// 7. Use Runic Power on Death Coil

	// If we have spent all our runes and ERW is ready, lets use it!
	if dk.EmpowerRuneWeapon.IsReady(sim) && dk.AllRunesSpent() {
		if dk.Presence == deathknight.UnholyPresence && dk.BloodTap.IsReady(sim) {
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
			dk.BloodPresence.Cast(sim, dk.CurrentTarget)
		}
		dk.EmpowerRuneWeapon.Cast(sim, dk.CurrentTarget)
	}

	// TODO: should we make this the default or somehow configurable?

	useBTForGF := true
	// If we need GF and we can't cast it right now, but BT is ready, lets use it!
	canGFWithBT := !dk.GhoulFrenzyAura.IsActive() && dk.BloodTap.CanCast(sim)

	if dk.Talents.Desolation > 0 && !dk.DesolationAura.IsActive() && dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.FrostFeverDisease[dk.CurrentTarget.Index].RemainingDuration(sim) < time.Second*4 && dk.IcyTouch.CanCast(sim) {
		dk.IcyTouch.Cast(sim, dk.CurrentTarget)
	} else if dk.BloodPlagueDisease[dk.CurrentTarget.Index].RemainingDuration(sim) < time.Second*4 && dk.PlagueStrike.CanCast(sim) {
		dk.PlagueStrike.Cast(sim, dk.CurrentTarget)
	} else if (!dk.GhoulFrenzyAura.IsActive() && dk.GhoulFrenzy.CanCast(sim)) || (useBTForGF && canGFWithBT) {
		if !dk.GhoulFrenzy.CanCast(sim) && dk.BloodTap.CanCast(sim) {
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
		}
		dk.GhoulFrenzy.Cast(sim, dk.CurrentTarget)
	} else if dk.SummonGargoyle.CanCast(sim) {
		dk.SummonGargoyle.Cast(sim, dk.CurrentTarget)
	} else if dk.Talents.Reaping > 0 && dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.ScourgeStrike.CanCast(sim) {
		dk.ScourgeStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.DeathCoil.CanCast(sim) {
		dk.DeathCoil.Cast(sim, dk.CurrentTarget)
	} else {
		if dk.HornOfWinter.CanCast(sim) {
			dk.HornOfWinter.Cast(sim, dk.CurrentTarget)
		} else {
			// This means we dont have the resources to do anything.
			dk.WaitUntil(sim, dk.RunicPowerBar.AnySpentRuneReadyAt())
			return 0
		}
	}

	return dk.NextGCDAt()
}

func (dk *DpsDeathknight) RotationActionCallback_AutoDW(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {

	if !dk.GCD.IsReady(sim) {
		return dk.NextGCDAt()
	}

	// 1. If you have desolation - Maintain Desolation
	// 2. Maintain Frost Fever
	// 3. Maintain Blood Plague
	// 4. Maintain DnD
	// 5. Use Gary
	// 7. Use Runic Power on Death Coil

	// If we need to use DnD and we don't have the runes, pop ERW
	if dk.EmpowerRuneWeapon.IsReady(sim) && dk.DeathAndDecay.IsReady(sim) && !dk.DeathAndDecay.CanCast(sim) {
		dk.EmpowerRuneWeapon.Cast(sim, dk.CurrentTarget)
	}

	if dk.FrostFeverDisease[dk.CurrentTarget.Index].IsActive() && dk.BloodPlagueDisease[dk.CurrentTarget.Index].IsActive() && dk.GhoulFrenzyAura.IsActive() && dk.DeathAndDecayDot.IsActive() {
		if dk.Presence == deathknight.UnholyPresence && dk.BloodTap.IsReady(sim) {
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
			dk.BloodPresence.Cast(sim, dk.CurrentTarget)
		}
	}

	// TODO: should we make this the default or somehow configurable?

	useBTForGF := true
	// If we need GF and we can't cast it right now, but BT is ready, lets use it!
	canGFWithBT := !dk.GhoulFrenzyAura.IsActive() && dk.BloodTap.CanCast(sim)

	if dk.Talents.Desolation > 0 && !dk.DesolationAura.IsActive() && dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.FrostFeverDisease[dk.CurrentTarget.Index].RemainingDuration(sim) < time.Second*4 && dk.IcyTouch.CanCast(sim) {
		dk.IcyTouch.Cast(sim, dk.CurrentTarget)
	} else if dk.BloodPlagueDisease[dk.CurrentTarget.Index].RemainingDuration(sim) < time.Second*4 && dk.PlagueStrike.CanCast(sim) {
		dk.PlagueStrike.Cast(sim, dk.CurrentTarget)
	} else if (!dk.GhoulFrenzyAura.IsActive() && dk.GhoulFrenzy.CanCast(sim)) || (useBTForGF && canGFWithBT) {
		if !dk.GhoulFrenzy.CanCast(sim) && dk.BloodTap.CanCast(sim) {
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
		}
		dk.GhoulFrenzy.Cast(sim, dk.CurrentTarget)
	} else if !dk.DeathAndDecayDot.IsActive() && dk.DeathAndDecay.CanCast(sim) {
		dk.DeathAndDecay.Cast(sim, dk.CurrentTarget)
	} else if dk.SummonGargoyle.CanCast(sim) {
		dk.SummonGargoyle.Cast(sim, dk.CurrentTarget)
	} else if dk.DeathCoil.CanCast(sim) {
		dk.DeathCoil.Cast(sim, dk.CurrentTarget)
	} else {
		waitUntil := dk.RunicPowerBar.AnySpentRuneReadyAt()

		if !dk.DeathAndDecay.IsReady(sim) {
			numDeath := dk.CurrentDeathRunes()
			numBlood := dk.CurrentBloodRunes()
			numFrost := dk.CurrentFrostRunes()
			numUnholy := dk.CurrentUnholyRunes()
			if numDeath+numBlood >= 2 {
				// This means we have the runes for DnD but the dot is already up..
				//  might as well cast blood boil!
				dk.BloodBoil.Cast(sim, dk.CurrentTarget)
			} else if numDeath+numFrost >= 2 {
				dk.IcyTouch.Cast(sim, dk.CurrentTarget)
			} else if numDeath+numUnholy >= 2 {
				dk.PlagueStrike.Cast(sim, dk.CurrentTarget)
			}
		}

		if dk.GCD.IsReady(sim) {
			if dk.HornOfWinter.CanCast(sim) && dk.CurrentRunicPower() < 80 {
				dk.HornOfWinter.Cast(sim, dk.CurrentTarget)
			} else {
				if waitUntil == sim.CurrentTime {
					waitUntil = dk.AutoAttacks.NextAttackAt()
				}
				dk.WaitUntil(sim, waitUntil)
				return 0
			}
		}
	}

	return dk.NextGCDAt()
}
