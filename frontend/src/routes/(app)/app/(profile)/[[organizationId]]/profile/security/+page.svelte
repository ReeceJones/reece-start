<script lang="ts">
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';
	import { enhance, applyAction } from '$app/forms';
	import { CircleCheck, CircleX, Save } from 'lucide-svelte';
	import clsx from 'clsx/lite';
	import { invalidateAll } from '$app/navigation';
	import SettingsCard from '$lib/components/Settings/SettingsCard.svelte';
	import SettingsCardTitle from '$lib/components/Settings/SettingsCardTitle.svelte';
	import SettingsCardActions from '$lib/components/Settings/SettingsCardActions.svelte';

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
	<SettingsCardTitle>{$t('settings.security')}</SettingsCardTitle>
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
			<legend class="fieldset-legend">{$t('settings.fields.email.label')}</legend>
			<input
				type="email"
				name="email"
				required
				class="input"
				placeholder={$t('settings.fields.email.placeholder')}
				bind:value={email}
			/>
			<p class="fieldset-label">
				{$t('settings.fields.email.description')}
			</p>
		</fieldset>

		<fieldset class="fieldset">
			<legend class="fieldset-legend">{$t('settings.fields.updatePassword.label')}</legend>
			<input
				type="password"
				name="password"
				class="input"
				placeholder={$t('settings.fields.updatePassword.placeholder')}
				bind:value={password}
			/>
			<p class="fieldset-label">{$t('settings.fields.updatePassword.description')}</p>
			{#if password.length > 0 && password.length < 8}
				<p class="fieldset-label text-error">
					{$t('settings.fields.updatePassword.passwordTooShort')}
				</p>
			{/if}
		</fieldset>

		{#if password !== ''}
			<fieldset class="fieldset">
				<legend class="fieldset-legend">{$t('settings.fields.confirmPassword.label')}</legend>
				<input
					type="password"
					name="confirmPassword"
					class={clsx('input', password !== confirmPassword && 'input-error')}
					placeholder={$t('settings.fields.confirmPassword.placeholder')}
					bind:value={confirmPassword}
				/>
				{#if password !== confirmPassword}
					<p class={clsx('fieldset-label', 'text-error')}>
						{$t('settings.fields.confirmPassword.passwordDoesNotMatch')}
					</p>
				{/if}
			</fieldset>
		{/if}

		{#if form?.success}
			<div role="alert" class="mt-3 alert alert-success">
				<CircleCheck />
				<span>{$t('settings.success.profileUpdated')}</span>
			</div>
		{:else if form?.success === false}
			<div role="alert" class="mt-3 alert alert-error">
				<CircleX />
				<span>{form.message ?? $t('settings.success.profileUpdateError')}</span>
			</div>
		{/if}

		<SettingsCardActions>
			<button type="submit" class="btn btn-primary" disabled={!canSubmit || submitting}>
				{#if submitting}
					<span class="loading loading-spinner"></span>
				{:else}
					<Save />
				{/if}
				<span>{$t('save')}</span>
			</button>
		</SettingsCardActions>
	</form>
</SettingsCard>
