<script lang="ts">
	import { localeToLanguageName } from '$lib/locale';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';

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
			>Is this organization associated with a registered business?</legend
		>
		<select name="entityType" class="select" bind:value={onboardingState.entityType}>
			<option value="company">Yes</option>
			<option value="individual">No</option>
		</select>
		<p class="fieldset-label">
			Select 'Yes' if you are a registered business (e.g. LLC, corporation, etc.). Otherwise, select
			'No'.
		</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Language</legend>
		<select name="locale" class="select" bind:value={onboardingState.locale}>
			<!-- TODO: add i8n here -->
			<option value="en">{localeToLanguageName('en', onboardingState.locale)}</option>
			<option value="es">{localeToLanguageName('en', onboardingState.locale)}</option>
			<option value="fr">{localeToLanguageName('en', onboardingState.locale)}</option>
			<option value="de">{localeToLanguageName('en', onboardingState.locale)}</option>
		</select>
		<p class="fieldset-label">
			Select the language you want to use for your organization. We will use this language for your
			organization's invoices and any communications we send to you.
		</p>
	</fieldset>
</OnboardingStepContainer>
