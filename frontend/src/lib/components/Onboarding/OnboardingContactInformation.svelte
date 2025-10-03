<script lang="ts">
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';

	const {
		hidden,
		onboardingState = $bindable()
	}: { hidden: boolean; onboardingState: CreateOrganizationFormData } = $props();
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
		<p class="fieldset-label">
			Enter an email address we can contact your organization at, if you have one
		</p>
	</fieldset>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">Contact Phone</legend>
		<input
			type="tel"
			name="contactPhone"
			class="input"
			placeholder="Phone"
			bind:value={onboardingState.contactPhone}
		/>
		<p class="fieldset-label">
			Enter a phone number we can contact your organization at, if you have one
		</p>
	</fieldset>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">Website</legend>
		<input
			name="websiteUrl"
			class="input"
			placeholder="Website"
			bind:value={onboardingState.websiteUrl}
		/>
		<p class="fieldset-label">Enter your website URL, if you have one</p>
	</fieldset>
</OnboardingStepContainer>
