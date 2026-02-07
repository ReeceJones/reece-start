<script lang="ts">
	import CountryOptions from '../CountryOptions.svelte';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import StateOptions from '../StateOptions.svelte';
	import { State } from 'country-state-city';
	import * as m from '$lib/paraglide/messages';
	import * as Field from '$lib/components/ui/field';
	import { Input } from '$lib/components/ui/input';

	const {
		hidden,
		onboardingState = $bindable()
	}: { hidden: boolean; onboardingState: CreateOrganizationFormData } = $props();
</script>

<OnboardingStepContainer {hidden}>
	<Field.Field>
		<Field.Label for="addressCountry">{m.onboarding__address_step__country()}</Field.Label>
		<select
			id="addressCountry"
			name="addressCountry"
			required
			class="h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-base shadow-xs ring-offset-background transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30"
			bind:value={onboardingState.addressCountry}
		>
			<CountryOptions />
		</select>
		<Field.Description>{m.onboarding__address_step__select_country()}</Field.Description>
	</Field.Field>

	<Field.Field>
		<Field.Label for="addressLine1">{m.onboarding__address_step__address()}</Field.Label>
		<Input
			type="text"
			id="addressLine1"
			name="addressLine1"
			required
			class="input"
			placeholder={m.onboarding__address_step__address()}
			bind:value={onboardingState.addressLine1}
		/>
		<Field.Description>{m.onboarding__address_step__enter_street_address()}</Field.Description>
	</Field.Field>

	<Field.Field>
		<Field.Label for="addressLine2">{m.onboarding__address_step__address_line2()}</Field.Label>
		<Input
			type="text"
			id="addressLine2"
			name="addressLine2"
			class="input"
			placeholder={m.onboarding__address_step__address_line2()}
			bind:value={onboardingState.addressLine2}
		/>
		<Field.Description>
			{m.onboarding__address_step__address_line2_description()}
		</Field.Description>
	</Field.Field>

	<Field.Field>
		<Field.Label for="addressCity">{m.onboarding__address_step__city()}</Field.Label>
		<Input
			type="text"
			id="addressCity"
			name="addressCity"
			required
			class="input"
			placeholder={m.onboarding__address_step__city()}
			bind:value={onboardingState.addressCity}
		/>
		<Field.Description>{m.onboarding__address_step__enter_city()}</Field.Description>
	</Field.Field>

	{#if onboardingState.addressCountry && State.getStatesOfCountry(onboardingState.addressCountry).length > 0}
		<Field.Field>
			<Field.Label for="addressStateOrProvince">{m.onboarding__address_step__state()}</Field.Label>
			<select
				id="addressStateOrProvince"
				name="addressStateOrProvince"
				class="h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-base shadow-xs ring-offset-background transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30"
				bind:value={onboardingState.addressStateOrProvince}
			>
				{#if onboardingState.addressCountry}
					<StateOptions countryCode={onboardingState.addressCountry} />
				{/if}
			</select>
			<Field.Description>{m.onboarding__address_step__select_state_or_province()}</Field.Description>
		</Field.Field>
	{/if}

	<Field.Field>
		<Field.Label for="addressZip">{m.onboarding__address_step__zip()}</Field.Label>
		<Input
			type="text"
			id="addressZip"
			name="addressZip"
			required
			class="input"
			placeholder={m.onboarding__address_step__zip()}
			bind:value={onboardingState.addressZip}
		/>
		<Field.Description>{m.onboarding__address_step__enter_zip()}</Field.Description>
	</Field.Field>
</OnboardingStepContainer>
