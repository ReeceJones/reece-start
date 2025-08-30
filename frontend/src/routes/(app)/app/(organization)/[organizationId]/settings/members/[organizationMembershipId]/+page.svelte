<script lang="ts">
	import { ArrowLeft, CircleCheck, CircleX, Save, Trash } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { API_TYPES } from '$lib/schemas/api';
	import { applyAction, enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';

	const { data, form }: PageProps = $props();

	$inspect(form);

	let submitting = $state(false);
	let role = $state(data.organizationMembership.data.attributes.role);

	const user = $derived(
		data.organizationMembership.included.filter((i) => i.type === API_TYPES.user)[0]
	);
	const canSubmit = $derived(role !== data.organizationMembership.data.attributes.role);
</script>

<div class="space-y-10">
	<div class="space-y-2">
		<button class="btn btn-ghost" onclick={() => history.back()}>
			<ArrowLeft class="size-4" />
			Back
		</button>
		<div class="flex gap-3">
			<img
				src={user.meta.logoDistributionUrl}
				alt={user.attributes.name}
				class="rounded-box size-20"
			/>
			<div class="flex flex-col">
				<h2 class="card-title">{user.attributes.name}</h2>
				<a href={`mailto:${user.attributes.email}`} class="link text-sm text-gray-500">
					{user.attributes.email}
				</a>
			</div>
		</div>

		<form
			class="space-y-2"
			method="post"
			enctype="multipart/form-data"
			use:enhance={() => {
				submitting = true;

				return ({ result }) => {
					applyAction(result);
					invalidateAll();
					submitting = false;
				};
			}}
		>
			<fieldset class="fieldset">
				<legend class="fieldset-legend">Role</legend>
				<select class="select select-bordered" name="role" bind:value={role}>
					<option value="admin">Admin</option>
					<option value="member">Member</option>
				</select>
			</fieldset>

			{#if form?.success}
				<div role="alert" class="alert alert-success mt-3">
					<CircleCheck />
					<span>The member has been updated!</span>
				</div>
			{:else if form?.success === false}
				<div role="alert" class="alert alert-error mt-3">
					<CircleX />
					<span
						>{form.message ??
							'There was an error updating the member. Make sure you have filled out all the fields correctly.'}</span
					>
				</div>
			{/if}

			<button class="btn btn-primary" disabled={!canSubmit || submitting}>
				{#if submitting}
					<span class="loading loading-spinner"></span>
				{:else}
					<Save class="size-4" />
				{/if}
				Save
			</button>
		</form>
	</div>
	<div class="divider"></div>
	<div>
		<h2 class="card-title mb-4">Danger Zone</h2>
		<button class="btn btn-error btn-outline">
			<Trash class="size-4" />
			Remove member
		</button>
	</div>
</div>
