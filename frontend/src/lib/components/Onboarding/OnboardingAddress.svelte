<script lang="ts">
	import CountryOptions from '../CountryOptions.svelte';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import StateOptions from '../StateOptions.svelte';
	import { State } from 'country-state-city';
	import { t } from '$lib/i18n';

	const {
		hidden,
		onboardingState = $bindable()
	}: { hidden: boolean; onboardingState: CreateOrganizationFormData } = $props();
</script>

<OnboardingStepContainer {hidden}>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.addressStep.country')}</legend>
		<select
			name="addressCountry"
			class="select"
			bind:value={onboardingState.addressCountry}
			placeholder={$t('onboarding.addressStep.selectCountry')}
		>
			<CountryOptions />
		</select>
		<p class="fieldset-label">{$t('onboarding.addressStep.selectCountry')}</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.addressStep.address')}</legend>
		<input
			type="text"
			name="addressLine1"
			class="input"
			placeholder={$t('onboarding.addressStep.address')}
			bind:value={onboardingState.addressLine1}
		/>
		<p class="fieldset-label">{$t('onboarding.addressStep.enterStreetAddress')}</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.addressStep.addressLine2')}</legend>
		<input
			type="text"
			name="addressLine2"
			class="input"
			placeholder={$t('onboarding.addressStep.addressLine2')}
			bind:value={onboardingState.addressLine2}
		/>
		<p class="fieldset-label">
			{$t('onboarding.addressStep.addressLine2Description')}
		</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.addressStep.city')}</legend>
		<input
			type="text"
			name="addressCity"
			class="input"
			placeholder={$t('onboarding.addressStep.city')}
			bind:value={onboardingState.addressCity}
		/>
		<p class="fieldset-label">{$t('onboarding.addressStep.enterCity')}</p>
	</fieldset>

	{#if onboardingState.addressCountry && State.getStatesOfCountry(onboardingState.addressCountry).length > 0}
		<fieldset class="fieldset">
			<legend class="fieldset-legend">{$t('onboarding.addressStep.state')}</legend>
			<select
				name="addressStateOrProvince"
				class="select"
				bind:value={onboardingState.addressStateOrProvince}
				placeholder={$t('onboarding.addressStep.selectStateOrProvince')}
			>
				{#if onboardingState.addressCountry}
					<StateOptions countryCode={onboardingState.addressCountry} />
				{/if}
			</select>
			<p class="fieldset-label">{$t('onboarding.addressStep.selectStateOrProvince')}</p>
		</fieldset>
	{/if}

	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.addressStep.zip')}</legend>
		<input
			type="text"
			name="addressZip"
			class="input"
			placeholder={$t('onboarding.addressStep.zip')}
			bind:value={onboardingState.addressZip}
		/>
		<p class="fieldset-label">{$t('onboarding.addressStep.enterZip')}</p>
	</fieldset>
</OnboardingStepContainer>
