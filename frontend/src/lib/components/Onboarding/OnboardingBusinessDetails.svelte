<script lang="ts">
	import { localeToLanguageName } from '$lib/locale';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import { t } from '$lib/i18n';

	const {
		hidden,
		onboardingState = $bindable()
	}: {
		hidden: boolean;
		onboardingState: CreateOrganizationFormData;
	} = $props();
</script>

<OnboardingStepContainer {hidden}>
	<fieldset class="fieldset">
		<legend class="fieldset-legend"
			>{$t('onboarding.businessDetailsStep.registeredBusiness')}</legend
		>
		<select name="entityType" class="select" bind:value={onboardingState.entityType}>
			<option value="company">{$t('onboarding.businessDetailsStep.yes')}</option>
			<option value="individual">{$t('onboarding.businessDetailsStep.no')}</option>
		</select>
		<p class="fieldset-label">
			{$t('onboarding.businessDetailsStep.registeredBusinessDescription')}
		</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.businessDetailsStep.language')}</legend>
		<select name="locale" class="select" bind:value={onboardingState.locale}>
			<!-- TODO: add i8n here -->
			<option value="en">{localeToLanguageName('en', 'en')}</option>
			<option value="es">{localeToLanguageName('en', 'es')}</option>
			<option value="fr">{localeToLanguageName('en', 'fr')}</option>
			<option value="de">{localeToLanguageName('en', 'de')}</option>
		</select>
		<p class="fieldset-label">
			{$t('onboarding.businessDetailsStep.languageDescription')}
		</p>
	</fieldset>
</OnboardingStepContainer>
