<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import * as Dialog from '$lib/components/ui/dialog';
	import type { InputProps } from '../ui/input/input.svelte';
	import * as m from '$lib/paraglide/messages';
	import LogoCrop from './LogoCrop.svelte';

	export type LogoInputProps = InputProps & {
		logo: FileList | null | undefined;
	};
	let { logo = $bindable(), id, name, disabled, ...restProps }: LogoInputProps = $props();
	let logoInput: HTMLInputElement | null = $state(null);
	let uncroppedLogo = $state<FileList | null | undefined>(undefined);
	let logoCropDialogOpen = $state(false);

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

<Input
	{...restProps}
	type="file"
	{id}
	{name}
	{disabled}
	class="file-input"
	accept="image/*"
	multiple={false}
	bind:ref={logoInput}
	onchange={(e: Event) => {
		if (disabled) return;
		const target = e.target as HTMLInputElement;
		uncroppedLogo = target.files;
		logoCropDialogOpen = true;
	}}
/>

<Dialog.Root bind:open={logoCropDialogOpen}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>
				{m.profile__update_image()}
			</Dialog.Title>
			<Dialog.Description>
				{m.profile__edit_image_description()}
			</Dialog.Description>
		</Dialog.Header>
		{#if uncroppedLogo && uncroppedLogo.length > 0}
			<LogoCrop
				imageFile={uncroppedLogo[0]}
				onCancel={() => {
					resetLogoUpload();
					logoCropDialogOpen = false;
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

					logoCropDialogOpen = false;
				}}
			/>
		{/if}
	</Dialog.Content>
</Dialog.Root>
