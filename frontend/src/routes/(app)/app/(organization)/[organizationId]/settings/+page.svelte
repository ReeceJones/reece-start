<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance, applyAction } from '$app/forms';
	import { CircleCheck, CircleX, Save } from 'lucide-svelte';
	import { invalidateAll } from '$app/navigation';
	import LogoCrop from '$lib/components/Logo/LogoCrop.svelte';
	import SettingsCard from '$lib/components/Settings/SettingsCard.svelte';
	import SettingsCardTitle from '$lib/components/Settings/SettingsCardTitle.svelte';
	import {
		getFormDataFromOrganization,
		type OrganizationFormData
	} from '$lib/schemas/organization';
	import deepEqual from 'deep-equal';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import SettingsCardActions from '$lib/components/Settings/SettingsCardActions.svelte';
	import { t } from '$lib/i18n';

	const { data, form }: PageProps = $props();

	let logoInput: HTMLInputElement | null = null;
	let logoCropModal: HTMLDialogElement | null = null;

	let submitting = $state(false);
	let uncroppedLogo = $state<FileList | null | undefined>(undefined);
	let formData: OrganizationFormData = $state(getFormDataFromOrganization(data.organization));

	// Simple derived state that just returns the appropriate URL
	const logoPreview = $derived(
		formData.logo && formData.logo.length > 0
			? URL.createObjectURL(formData.logo[0])
			: data.organization.data.meta.logoDistributionUrl
	);
	const isDirty = $derived(!deepEqual(formData, getFormDataFromOrganization(data.organization)));
	const isValid = $derived(!!formData.name);
	const canUpdate = $derived(hasScope(UserScope.OrganizationUpdate));
	const canSubmit = $derived(isDirty && isValid && canUpdate && !submitting);

	function resetLogoUpload() {
		formData.logo = undefined;
		uncroppedLogo = undefined;
		// reset the logo input
		if (logoInput) {
			// Completely reset the value of the input
			logoInput.value = '';
			logoInput.files = null;
		}
	}
</script>

<SettingsCard>
	<SettingsCardTitle>{$t('settings.organization.general.title')}</SettingsCardTitle>
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
		<fieldset class="fieldset">
			<legend class="fieldset-legend">{$t('settings.organization.general.logo.label')}</legend>
			{#if logoPreview}
				<img src={logoPreview} alt="Organization logo" class="aspect-square w-48 rounded-box" />
			{:else}
				<div class="flex aspect-square w-48 items-center justify-center rounded-box bg-base-300">
					<span class="text-base-content/50"
						>{$t('settings.organization.general.logo.noLogoUploaded')}</span
					>
				</div>
			{/if}
			<input
				type="file"
				name="logo"
				class="file-input"
				accept="image/*"
				multiple={false}
				bind:this={logoInput}
				onchange={(e: Event) => {
					const target = e.target as HTMLInputElement;
					uncroppedLogo = target.files;
					if (logoCropModal) {
						logoCropModal.showModal();
					}
				}}
				disabled={!canUpdate}
			/>
			<p class="fieldset-label">{$t('settings.organization.general.logo.description')}</p>
		</fieldset>

		<fieldset class="fieldset">
			<legend class="fieldset-legend">{$t('settings.organization.general.name.label')}</legend>
			<input
				type="text"
				name="name"
				required
				class="input"
				placeholder={$t('settings.organization.general.name.placeholder')}
				bind:value={formData.name}
				disabled={!canUpdate}
			/>
			<p class="fieldset-label">{$t('settings.organization.general.name.description')}</p>
		</fieldset>

		<fieldset class="fieldset">
			<legend class="fieldset-legend"
				>{$t('settings.organization.general.description.label')}</legend
			>
			<textarea
				name="description"
				class="textarea"
				placeholder={$t('settings.organization.general.description.placeholder')}
				bind:value={formData.description}
				maxlength={255}
				disabled={!canUpdate}
			></textarea>
			<p class="fieldset-label">{$t('settings.organization.general.description.description')}</p>
		</fieldset>

		{#if form?.success}
			<div role="alert" class="mt-3 alert alert-success">
				<CircleCheck />
				<span>{$t('settings.organization.general.success.organizationUpdated')}</span>
			</div>
		{:else if form?.success === false}
			<div role="alert" class="mt-3 alert alert-error">
				<CircleX />
				<span
					>{form.message ??
						$t('settings.organization.general.success.organizationUpdateError')}</span
				>
			</div>
		{/if}

		<SettingsCardActions>
			<button type="submit" class="btn btn-primary" disabled={!canSubmit}>
				{#if submitting}
					<span class="loading loading-spinner"></span>
				{:else}
					<Save />
				{/if}
				<span>{$t('save')}</span>
			</button>
		</SettingsCardActions>
	</form>

	<dialog id="logo-crop-modal" class="modal" bind:this={logoCropModal}>
		<div class="modal-box">
			<h3 class="text-lg font-bold">{$t('settings.organization.general.logo.updateLogo')}</h3>
			<p class="py-4">{$t('settings.organization.general.logo.updateLogoDescription')}</p>
			{#if uncroppedLogo && uncroppedLogo.length > 0}
				<LogoCrop
					imageFile={uncroppedLogo[0]}
					onCancel={() => {
						resetLogoUpload();
						if (logoCropModal) {
							logoCropModal.close();
						}
					}}
					onSave={(file) => {
						console.log('Cropped file received:', file);
						console.log('File size:', file.size);
						console.log('File type:', file.type);

						const dt = new DataTransfer();
						dt.items.add(file);
						formData.logo = dt.files;
						uncroppedLogo = undefined;

						// set the value of the input to the new cropped image
						if (logoInput) {
							logoInput.files = dt.files;
						}

						if (logoCropModal) {
							logoCropModal.close();
						}
					}}
				/>
			{/if}
		</div>
		<form method="dialog" class="modal-backdrop" onsubmit={() => resetLogoUpload()}>
			<button>{$t('close')}</button>
		</form>
	</dialog>
</SettingsCard>
