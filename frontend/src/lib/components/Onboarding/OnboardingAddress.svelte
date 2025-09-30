<script lang="ts">
	import CountryOptions from '../CountryOptions.svelte';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import StateOptions from '../StateOptions.svelte';
	import { State } from 'country-state-city';

	const {
		hidden,
		onboardingState = $bindable()
	}: { hidden: boolean; onboardingState: CreateOrganizationFormData } = $props();
</script>

<OnboardingStepContainer {hidden}>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">Country</legend>
		<select
			name="addressCountry"
			class="select"
			bind:value={onboardingState.addressCountry}
			placeholder="Select your country"
		>
			<CountryOptions />
		</select>
		<p class="fieldset-label">Select your country</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Address</legend>
		<input
			type="text"
			name="addressLine1"
			class="input"
			placeholder="Address"
			bind:value={onboardingState.addressLine1}
		/>
		<p class="fieldset-label">Enter your street address</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Address Line 2</legend>
		<input
			type="text"
			name="addressLine2"
			class="input"
			placeholder="Address Line 2"
			bind:value={onboardingState.addressLine2}
		/>
		<p class="fieldset-label">
			If you have a second line of address (e.g. apartment or suite number), enter it here
		</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">City</legend>
		<input
			type="text"
			name="addressCity"
			class="input"
			placeholder="City"
			bind:value={onboardingState.addressCity}
		/>
		<p class="fieldset-label">Enter your city</p>
	</fieldset>

	{#if onboardingState.addressCountry && State.getStatesOfCountry(onboardingState.addressCountry).length > 0}
		<fieldset class="fieldset">
			<legend class="fieldset-legend">State</legend>
			<select
				name="addressStateOrProvince"
				class="select"
				bind:value={onboardingState.addressStateOrProvince}
				placeholder="Select your state or province"
			>
				{#if onboardingState.addressCountry}
					<StateOptions countryCode={onboardingState.addressCountry} />
				{/if}
			</select>
			<p class="fieldset-label">Select your state or province</p>
		</fieldset>
	{/if}

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Zip</legend>
		<input
			type="text"
			name="addressZip"
			class="input"
			placeholder="Zip"
			bind:value={onboardingState.addressZip}
		/>
		<p class="fieldset-label">Enter your zip</p>
	</fieldset>
</OnboardingStepContainer>
