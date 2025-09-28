<script lang="ts">
	import CountryOptions from '../CountryOptions.svelte';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';

	const {
		hidden,
		userName,
		canSubmit = $bindable()
	}: { hidden: boolean; userName: string; canSubmit: Record<number, boolean> } = $props();

	let residingCountry = $state<string>('US');
	let businessType = $state<'company' | 'individual'>('individual');
	let registeredBusinessName = $state<string>('');
	let firstName = $state<string>('');
	let lastName = $state<string>('');
	let locale = $state<string>('en');
	let currency = $state<string>('usd');

	$effect(() => {
		const identityFieldsValid =
			businessType === 'company' ? !!registeredBusinessName : !!firstName && !!lastName;
		canSubmit[3] =
			!!businessType && !!residingCountry && !!locale && !!currency && identityFieldsValid;
	});
</script>

<OnboardingStepContainer {hidden}>
	<fieldset class="fieldset">
		<legend class="fieldset-legend"
			>Is this organization associated with a registered business?</legend
		>
		<select name="entityType" class="select" bind:value={businessType}>
			<option value="company">Yes</option>
			<option value="individual">No</option>
		</select>
		<p class="fieldset-label">
			Select 'Yes' if you are a registered business (e.g. LLC, corporation, etc.). Otherwise, select
			'No'.
		</p>
	</fieldset>

	<!-- Common -->
	<fieldset class="fieldset">
		<legend class="fieldset-legend">Residing Country</legend>
		<select name="residingCountry" class="select" bind:value={residingCountry}>
			<CountryOptions />
		</select>
		<p class="fieldset-label">
			{#if businessType === 'company'}
				Enter the country where the business is registered.
			{:else}
				Enter the country where you reside.
			{/if}
			<br />
			This will impact the currency and payment methods available.
		</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Language</legend>
		<select name="locale" class="select" bind:value={locale}>
			<option value="en">English</option>
			<option value="es">Spanish</option>
			<option value="fr">French</option>
			<option value="de">German</option>
		</select>
		<p class="fieldset-label">Select the language you want to use for your organization</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Currency</legend>
		<select name="currency" class="select" bind:value={currency}>
			{#if residingCountry === 'US'}
				<option value="usd">$ USD</option>
			{:else}
				<option value="eur">€ EUR</option>
			{/if}
		</select>
		<p class="fieldset-label">Select the currency you want to use for your organization</p>
	</fieldset>

	<!-- Individual -->
	{#if businessType === 'individual'}
		<fieldset class="fieldset">
			<legend class="fieldset-legend">First Name</legend>
			<input
				type="text"
				name="firstName"
				class="input"
				placeholder="First Name"
				defaultValue={userName.split(' ')[0] ?? ''}
				bind:value={firstName}
			/>
			<p class="fieldset-label">Enter your first name</p>
		</fieldset>

		<fieldset class="fieldset">
			<legend class="fieldset-legend">Last Name</legend>
			<input
				type="text"
				name="lastName"
				class="input"
				placeholder="Last Name"
				defaultValue={userName.split(' ')[1] ?? ''}
				bind:value={lastName}
			/>
			<p class="fieldset-label">Enter your last name</p>
		</fieldset>
	{/if}

	<!-- Company -->
	{#if businessType === 'company'}
		<fieldset class="fieldset">
			<legend class="fieldset-legend">Business Name</legend>
			<input
				type="text"
				name="registeredBusinessName"
				class="input"
				placeholder="Business Name"
				bind:value={registeredBusinessName}
			/>
			<p class="fieldset-label">Enter the registered name of the business</p>
		</fieldset>
	{/if}
</OnboardingStepContainer>
