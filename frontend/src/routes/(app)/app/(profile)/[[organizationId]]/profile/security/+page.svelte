<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance, applyAction } from '$app/forms';
	import { CircleCheck, CircleX, Save } from 'lucide-svelte';
	import clsx from 'clsx/lite';
	import { invalidateAll } from '$app/navigation';
	import SettingsCard from '$lib/components/Settings/SettingsCard.svelte';
	import SettingsCardTitle from '$lib/components/Settings/SettingsCardTitle.svelte';

	let { data, form }: PageProps = $props();

	let submitting = $state(false);
	let email = $state(data.user.data.attributes.email);
	let userProfile = $state(data.user.data);
	let password = $state('');
	let confirmPassword = $state('');

	let canSubmit = $derived.by(() => {
		if (password !== '' && password.length < 8) {
			return false;
		}

		if (password !== '' && password !== confirmPassword) {
			return false;
		}

		return true;
	});
</script>

<SettingsCard>
	<SettingsCardTitle>Security</SettingsCardTitle>
	<form
		method="post"
		use:enhance={() => {
			submitting = true;

			return ({ result }) => {
				invalidateAll();
				applyAction(result);
				submitting = false;
			};
		}}
		enctype="multipart/form-data"
	>
		<input type="hidden" tabindex="-1" name="userId" value={userProfile.id} />

		<fieldset class="fieldset">
			<legend class="fieldset-legend">Email</legend>
			<input
				type="email"
				name="email"
				required
				class="input"
				placeholder="Email"
				bind:value={email}
			/>
			<p class="fieldset-label">
				The email you use to log into your account and receive notifications
			</p>
		</fieldset>

		<fieldset class="fieldset">
			<legend class="fieldset-legend">Update Password</legend>
			<input
				type="password"
				name="password"
				class="input"
				placeholder="Password"
				bind:value={password}
			/>
			<p class="fieldset-label">Update your password used to sign in to your account</p>
			{#if password.length > 0 && password.length < 8}
				<p class="fieldset-label text-error">Password must be at least 8 characters long.</p>
			{/if}
		</fieldset>

		{#if password !== ''}
			<fieldset class="fieldset">
				<legend class="fieldset-legend">Confirm password</legend>
				<input
					type="password"
					name="confirmPassword"
					class={clsx('input', password !== confirmPassword && 'input-error')}
					placeholder="Confirm password"
					bind:value={confirmPassword}
				/>
				{#if password !== confirmPassword}
					<p class={clsx('fieldset-label', 'text-error')}>Passwords do not match.</p>
				{/if}
			</fieldset>
		{/if}

		{#if form?.success}
			<div role="alert" class="alert alert-success mt-3">
				<CircleCheck />
				<span>Your profile has been updated!</span>
			</div>
		{:else if form?.success === false}
			<div role="alert" class="alert alert-error mt-3">
				<CircleX />
				<span
					>{form.message ??
						'There was an error updating your profile. Make sure you have filled out all the fields correctly.'}</span
				>
			</div>
		{/if}

		<div class="card-actions mt-3 justify-start">
			<button type="submit" class="btn btn-primary" disabled={!canSubmit || submitting}>
				{#if submitting}
					<span class="loading loading-spinner"></span>
				{:else}
					<Save />
				{/if}
				<span>Save</span>
			</button>
		</div>
	</form>
</SettingsCard>
