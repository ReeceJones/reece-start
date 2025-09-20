<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance, applyAction } from '$app/forms';
	import { CircleCheck, CircleX, Plus } from 'lucide-svelte';
	import LogoCrop from '$lib/components/Logo/LogoCrop.svelte';

	let { form }: PageProps = $props();

	let submitting = $state(false);
	let name = $state('');
	let description = $state('');
	let logoInput: HTMLInputElement | null = null;
	let logo = $state<FileList | null | undefined>(undefined);
	let uncroppedLogo = $state<FileList | null | undefined>(undefined);
	let logoCropModal: HTMLDialogElement | null = null;
	let dirty = $derived(name !== '');

	// Simple derived state that returns a preview URL or placeholder
	const logoPreview = $derived(logo && logo.length > 0 ? URL.createObjectURL(logo[0]) : undefined);

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

<div class="card bg-base-200 border-base-300 border shadow-sm">
	<div class="card-body">
		<p class="card-title">Create Organization</p>
		<p>Create a new organization to do XYZ.</p>

		<form
			method="post"
			use:enhance={() => {
				submitting = true;

				return ({ result }) => {
					applyAction(result);
					submitting = false;
				};
			}}
			enctype="multipart/form-data"
		>
			<fieldset class="fieldset">
				<legend class="fieldset-legend">Logo</legend>
				{#if logoPreview}
					<img
						src={logoPreview}
						alt="Organization logo preview"
						class="rounded-box mb-4 aspect-square w-48"
					/>
				{:else}
					<div
						class="rounded-box bg-base-300 text-base-content/50 mb-4 flex aspect-square w-48 items-center justify-center"
					>
						<span>No logo selected</span>
					</div>
				{/if}
				<input
					type="file"
					name="logo"
					class="file-input"
					accept="image/*"
					multiple={false}
					bind:this={logoInput}
					onchange={(e) => {
						const target = e.target as HTMLInputElement;
						uncroppedLogo = target.files;
						if (logoCropModal) {
							logoCropModal.showModal();
						}
					}}
				/>
				<p class="fieldset-label">Upload your organization's logo (optional)</p>
			</fieldset>

			<fieldset class="fieldset">
				<legend class="fieldset-legend">Name</legend>
				<input type="text" name="name" bind:value={name} class="input" />
				<p class="fieldset-label">Enter a name for your organization</p>
			</fieldset>

			<fieldset class="fieldset">
				<legend class="fieldset-legend">Description</legend>
				<input type="text" name="description" bind:value={description} class="input" />
				<p class="fieldset-label">Enter a description for your organization</p>
			</fieldset>

			{#if form?.success}
				<div role="alert" class="alert alert-success">
					<CircleCheck />
					<span>Organization created successfully!</span>
				</div>
			{:else if form?.success === false}
				<div role="alert" class="alert alert-error">
					<CircleX />
					<span
						>{form.message ??
							'There was an error creating your organization. Make sure you have filled out all the fields correctly.'}</span
					>
				</div>
			{/if}

			<div class="card-actions mt-3 justify-start">
				<button type="submit" class="btn btn-primary" disabled={!dirty || submitting}>
					{#if submitting}
						<span class="loading loading-spinner"></span>
					{:else}
						<Plus />
					{/if}
					<span>Create Organization</span>
				</button>
			</div>
		</form>
	</div>
</div>

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
