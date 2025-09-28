<script lang="ts">
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';

	const { hidden }: { hidden: boolean } = $props();
</script>

<OnboardingStepContainer {hidden}>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">Contact Email</legend>
		<input type="email" name="contactEmail" class="input" placeholder="Email" />
		<p class="fieldset-label">
			Enter an email address we can contact your organization at, if you have one
		</p>
	</fieldset>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">Contact Phone</legend>
		<input type="tel" name="contactPhone" class="input" placeholder="Phone" />
		<p class="fieldset-label">
			Enter a phone number we can contact your organization at, if you have one
		</p>
	</fieldset>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">Website</legend>
		<input type="url" name="websiteUrl" class="input" placeholder="Website" />
		<p class="fieldset-label">Enter your website URL, if you have one</p>
	</fieldset>
</OnboardingStepContainer>
