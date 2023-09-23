import { UnitReference, UnitReference_Type as UnitType } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import {
	HealingPriest,
	HealingPriest_Rotation as PriestRotation,
	HealingPriest_Rotation_RotationType as RotationType,
	HealingPriest_Rotation_SpellOption as SpellOption,
} from '../core/proto/priest.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfPowerInfusion = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecHealingPriest>({
	fieldName: 'powerInfusionTarget',
	id: ActionId.fromSpellId(10060),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecHealingPriest>) => player.getSpecOptions().powerInfusionTarget?.type == UnitType.Player,
	setValue: (eventID: EventID, player: Player<Spec.SpecHealingPriest>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.powerInfusionTarget = UnitReference.create({
			type: newValue ? UnitType.Player : UnitType.Unknown,
			index: 0,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const InnerFire = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecHealingPriest>({
	fieldName: 'useInnerFire',
	id: ActionId.fromSpellId(25431),
});

export const Shadowfiend = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecHealingPriest>({
	fieldName: 'useShadowfiend',
	id: ActionId.fromSpellId(34433),
});

export const RapturesPerMinute = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHealingPriest>({
	fieldName: 'rapturesPerMinute',
	label: 'Raptures / Min',
	labelTooltip: 'Number of times to proc Rapture each minute (due to a PWS being fully absorbed).',
	showWhen: (player: Player<Spec.SpecHealingPriest>) => player.getTalents().rapture > 0,
	changeEmitter: (player: Player<Spec.SpecHealingPriest>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
});

export const HealingPriestRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecHealingPriest, RotationType>({
			fieldName: 'type',
			label: 'Type',
			values: [
				{ name: 'Cycle', value: RotationType.Cycle },
				{ name: 'Custom', value: RotationType.Custom },
			],
		}),
		InputHelpers.makeCustomRotationInput<Spec.SpecHealingPriest, SpellOption>({
			fieldName: 'customRotation',
			numColumns: 2,
			showCastsPerMinute: true,
			values: [
				{ actionId: ActionId.fromSpellId(25213), value: SpellOption.GreaterHeal },
				{ actionId: ActionId.fromSpellId(25235), value: SpellOption.FlashHeal },
				{ actionId: ActionId.fromSpellId(25222), value: SpellOption.Renew },
				{ actionId: ActionId.fromSpellId(25218), value: SpellOption.PowerWordShield },
				{ actionId: ActionId.fromSpellId(34866), value: SpellOption.CircleOfHealing },
				{ actionId: ActionId.fromSpellId(25316), value: SpellOption.PrayerOfHealing },
				{ actionId: ActionId.fromSpellId(33076), value: SpellOption.PrayerOfMending },
				{ actionId: ActionId.fromSpellId(53005), value: SpellOption.Penance },
				{ actionId: ActionId.fromSpellId(32546), value: SpellOption.BindingHeal },
			],
			showWhen: (player: Player<Spec.SpecHealingPriest>) => player.getRotation().type == RotationType.Custom,
		}),
	],
};
