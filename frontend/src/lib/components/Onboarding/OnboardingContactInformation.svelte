<script lang="ts">
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import { getPhoneCodeOptions, type PhoneCodeOption } from '$lib/phone-utils';
	import { t } from '$lib/i18n';

	const {
		hidden,
		onboardingState = $bindable()
	}: { hidden: boolean; onboardingState: CreateOrganizationFormData } = $props();

	const phoneCodeOptions: PhoneCodeOption[] = getPhoneCodeOptions();
</script>

<OnboardingStepContainer {hidden}>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.contactInformationStep.contactEmail')}</legend>
		<input
			type="email"
			name="contactEmail"
			class="input"
			placeholder={$t('onboarding.contactInformationStep.email')}
			bind:value={onboardingState.contactEmail}
		/>
		<p class="fieldset-label">{$t('onboarding.contactInformationStep.emailDescription')}</p>
	</fieldset>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.contactInformationStep.contactPhone')}</legend>

		<div class="flex gap-2">
			<select
				name="contactPhoneCountry"
				class="select w-48"
				placeholder={$t('onboarding.contactInformationStep.selectCountry')}
				bind:value={onboardingState.contactPhoneCountry}
			>
				{#each phoneCodeOptions as option (option.countryCode)}
					<option value={option.countryCode}>
						+{option.code} ({option.countryName}
						{option.flag})
					</option>
				{/each}
			</select>
			<input
				type="tel"
				name="contactPhone"
				class="input flex-1"
				placeholder={$t('onboarding.contactInformationStep.phoneNumber')}
				bind:value={onboardingState.contactPhone}
			/>
		</div>
		<p class="fieldset-label">{$t('onboarding.contactInformationStep.phoneDescription')}</p>
	</fieldset>
</OnboardingStepContainer>
