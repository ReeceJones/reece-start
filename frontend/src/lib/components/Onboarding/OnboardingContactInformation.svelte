<script lang="ts">
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import { getPhoneCodeOptions, type PhoneCodeOption } from '$lib/phone-utils';

	const {
		hidden,
		onboardingState = $bindable()
	}: { hidden: boolean; onboardingState: CreateOrganizationFormData } = $props();

	const phoneCodeOptions: PhoneCodeOption[] = getPhoneCodeOptions();
</script>

<OnboardingStepContainer {hidden}>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">Contact Email</legend>
		<input
			type="email"
			name="contactEmail"
			class="input"
			placeholder="Email"
			bind:value={onboardingState.contactEmail}
		/>
		<p class="fieldset-label">Enter an email address we can contact your organization at.</p>
	</fieldset>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">Contact Phone</legend>

		<select
			name="contactPhoneCountry"
			class="select w-48"
			bind:value={onboardingState.contactPhoneCountry}
		>
			{#each phoneCodeOptions as option (option.countryCode)}
				<option value={option.countryCode}>
					{option.flag} +{option.code} ({option.countryName})
				</option>
			{/each}
		</select>
		<input
			type="tel"
			name="contactPhone"
			class="input flex-1"
			placeholder="Phone number"
			bind:value={onboardingState.contactPhone}
		/>
		<p class="fieldset-label">Enter a phone number we can contact your organization at.</p>
	</fieldset>
</OnboardingStepContainer>
