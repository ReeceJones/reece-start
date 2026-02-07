<script lang="ts">
	import { localeToLanguageName } from '$lib/locale';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import * as m from '$lib/paraglide/messages';
	import * as Field from '$lib/components/ui/field';

	const {
		hidden,
		onboardingState = $bindable()
	}: {
		hidden: boolean;
		onboardingState: CreateOrganizationFormData;
	} = $props();
</script>

<OnboardingStepContainer {hidden}>
	<Field.Field>
		<Field.Label for="entityType"
			>{m.onboarding__business_details_step__registered_business()}</Field.Label
		>
		<select
			id="entityType"
			name="entityType"
			required
			class="h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-base shadow-xs ring-offset-background transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30"
			bind:value={onboardingState.entityType}
		>
			<option value="company">{m.onboarding__business_details_step__yes()}</option>
			<option value="individual">{m.onboarding__business_details_step__no()}</option>
		</select>
		<Field.Description>
			{m.onboarding__business_details_step__registered_business_description()}
		</Field.Description>
	</Field.Field>

	<Field.Field>
		<Field.Label for="locale">{m.onboarding__business_details_step__language()}</Field.Label>
		<select
			id="locale"
			name="locale"
			required
			class="h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-base shadow-xs ring-offset-background transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30"
			bind:value={onboardingState.locale}
		>
			<!-- TODO: add i8n here -->
			<option value="en">{localeToLanguageName('en', 'en')}</option>
			<option value="es">{localeToLanguageName('en', 'es')}</option>
			<option value="fr">{localeToLanguageName('en', 'fr')}</option>
			<option value="de">{localeToLanguageName('en', 'de')}</option>
		</select>
		<Field.Description>
			{m.onboarding__business_details_step__language_description()}
		</Field.Description>
	</Field.Field>
</OnboardingStepContainer>
