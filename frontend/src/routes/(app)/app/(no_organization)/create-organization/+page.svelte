<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance, applyAction } from '$app/forms';
	import { CircleCheck, CircleX, Plus, Save } from 'lucide-svelte';

	let { form }: PageProps = $props();

	let submitting = $state(false);
	let name = $state('');
	let description = $state('');
	let dirty = $derived(name !== '');
</script>

<div class="card bg-base-200 shadow-sm">
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
				<input type="file" name="logo" class="file-input" accept="image/*" />
				<p class="fieldset-label">Upload your organization's logo</p>
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
