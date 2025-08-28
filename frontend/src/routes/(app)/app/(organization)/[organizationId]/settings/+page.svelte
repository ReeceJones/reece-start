<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance, applyAction } from '$app/forms';
	import { CircleCheck, CircleX, Save } from 'lucide-svelte';
	import { invalidate } from '$app/navigation';
	import LogoCrop from '$lib/components/LogoCrop.svelte';

	let { data, form }: PageProps = $props();

	let submitting = $state(false);
	let name = $state(data.organization.data.attributes.name);
	let description = $state(data.organization.data.attributes.description || '');
	let organizationData = $state(data.organization.data);
	let logoInput: HTMLInputElement | null = null;
	let logo = $state<FileList | null | undefined>(undefined);
	let uncroppedLogo = $state<FileList | null | undefined>(undefined);
	let logoCropModal: HTMLDialogElement | null = null;

	// Simple derived state that just returns the appropriate URL
	const logoPreview = $derived(
		logo && logo.length > 0
			? URL.createObjectURL(logo[0])
			: data.organization.data.meta.logoDistributionUrl
	);

	let canSubmit = $derived.by(() => {
		if (!name) {
			return false;
		}

		return true;
	});

	function resetLogoUpload() {
		logo = undefined;
		uncroppedLogo = undefined;
		// reset the logo input
		if (logoInput) {
			// Completely reset the value of the input
			logoInput.value = '';
			logoInput.files = null;
		}
	}
</script>

<form
	method="post"
	use:enhance={() => {
		submitting = true;

		return ({ result }) => {
			invalidate(`/api/organizations/${organizationData.id}`);
			applyAction(result);
			submitting = false;
		};
	}}
	enctype="multipart/form-data"
>
	<input type="hidden" tabindex="-1" name="organizationId" value={organizationData.id} />

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Organization logo</legend>
		{#if logoPreview}
			<img src={logoPreview} alt="Organization logo" class="rounded-box aspect-square w-48" />
		{:else}
			<div class="rounded-box bg-base-300 flex aspect-square w-48 items-center justify-center">
				<span class="text-base-content/50">No logo uploaded</span>
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
		/>
		<p class="fieldset-label">Upload your organization logo</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Name</legend>
		<input
			type="text"
			name="name"
			required
			class="input"
			placeholder="Organization name"
			bind:value={name}
		/>
		<p class="fieldset-label">What should we call your organization?</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">Description</legend>
		<textarea
			name="description"
			class="textarea"
			placeholder="Organization description"
			bind:value={description}
			maxlength={255}
		></textarea>
		<p class="fieldset-label">A brief description of your organization</p>
	</fieldset>

	{#if form?.success}
		<div role="alert" class="alert alert-success mt-3">
			<CircleCheck />
			<span>Your organization has been updated!</span>
		</div>
	{:else if form?.success === false}
		<div role="alert" class="alert alert-error mt-3">
			<CircleX />
			<span
				>{form.message ??
					'There was an error updating your organization. Make sure you have filled out all the fields correctly.'}</span
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

<dialog id="logo-crop-modal" class="modal" bind:this={logoCropModal}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">Update logo</h3>
		<p class="py-4">Edit the logo to your liking and click save.</p>
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
					logo = dt.files;
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
		<button>close</button>
	</form>
</dialog>
