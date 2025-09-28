<script lang="ts">
	import CountryOptions from '../CountryOptions.svelte';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';

	const { hidden }: { hidden: boolean } = $props();

	let addressCountry = $state<string>('');
</script>

<OnboardingStepContainer {hidden}>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">Address</legend>
		<input type="text" name="addressLine1" class="input" placeholder="Address" />
		<p class="fieldset-label">Enter your street address</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Address Line 2</legend>
		<input type="text" name="addressLine2" class="input" placeholder="Address Line 2" />
		<p class="fieldset-label">
			If you have a second line of address (e.g. apartment or suite number), enter it here
		</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">City</legend>
		<input type="text" name="addressCity" class="input" placeholder="City" />
		<p class="fieldset-label">Enter your city</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">State</legend>
		<input
			type="text"
			name="addressStateOrProvince"
			class="input"
			placeholder="State or Province"
		/>
		<p class="fieldset-label">Enter your state or province</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Zip</legend>
		<input type="text" name="addressZip" class="input" placeholder="Zip" />
		<p class="fieldset-label">Enter your zip</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Country</legend>
		<select name="addressCountry" class="select" bind:value={addressCountry}>
			<CountryOptions />
		</select>
		<p class="fieldset-label">Select your country</p>
	</fieldset>
</OnboardingStepContainer>
