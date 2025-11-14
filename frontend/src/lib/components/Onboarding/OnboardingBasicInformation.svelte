<script lang="ts">
	import LogoCrop from '$lib/components/Logo/LogoCrop.svelte';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import { t } from '$lib/i18n';

	const {
		hidden,
		onboardingState = $bindable()
	}: {
		hidden: boolean;
		onboardingState: CreateOrganizationFormData;
	} = $props();

	let logoInput: HTMLInputElement | null = null;
	let uncroppedLogo = $state<FileList | null | undefined>(undefined);
	let logoCropModal: HTMLDialogElement | null = null;

	// Simple derived state that returns a preview URL or placeholder
	const logoPreview = $derived(
		onboardingState.logo && onboardingState.logo.length > 0
			? URL.createObjectURL(onboardingState.logo[0])
			: undefined
	);

	function resetLogoUpload() {
		onboardingState.logo = undefined;
		uncroppedLogo = undefined;
		// reset the logo input
		if (logoInput) {
			// Completely reset the value of the input
			logoInput.value = '';
			logoInput.files = null;
		}
	}
</script>

<OnboardingStepContainer {hidden}>
	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.basicInformationStep.name')}</legend>
		<input
			type="text"
			name="name"
			class="input"
			placeholder={$t('onboarding.basicInformationStep.name')}
			bind:value={onboardingState.name}
		/>
		<p class="fieldset-label">
			{$t('onboarding.basicInformationStep.nameDescription')}
		</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.basicInformationStep.description')}</legend>
		<textarea
			name="description"
			class="textarea"
			placeholder={$t('onboarding.basicInformationStep.description')}
			bind:value={onboardingState.description}
		></textarea>
		<p class="fieldset-label">{$t('onboarding.basicInformationStep.descriptionDescription')}</p>
	</fieldset>

	<fieldset class="fieldset">
		<legend class="fieldset-legend">{$t('onboarding.basicInformationStep.logo')}</legend>
		{#if logoPreview}
			<img
				src={logoPreview}
				alt={$t('onboarding.basicInformationStep.logoPreview')}
				class="rounded-box mb-4 aspect-square w-48"
			/>
		{:else}
			<div
				class="rounded-box bg-base-300 text-base-content/50 mb-4 flex aspect-square w-48 items-center justify-center"
			>
				<span>{$t('onboarding.basicInformationStep.noLogoSelected')}</span>
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
		<p class="fieldset-label">{$t('onboarding.basicInformationStep.uploadLogo')}</p>
	</fieldset>
</OnboardingStepContainer>

<dialog id="logo-crop-modal" class="modal" bind:this={logoCropModal}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">{$t('onboarding.basicInformationStep.updateLogo')}</h3>
		<p class="py-4">{$t('onboarding.basicInformationStep.editLogoDescription')}</p>
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
					onboardingState.logo = dt.files;
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
		<button>{$t('onboarding.basicInformationStep.close')}</button>
	</form>
</dialog>
