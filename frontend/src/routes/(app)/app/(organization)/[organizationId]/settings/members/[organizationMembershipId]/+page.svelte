<script lang="ts">
	import { ArrowLeft, CircleCheck, CircleX, Save, Trash, User } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { API_TYPES } from '$lib/schemas/api';
	import { enhance } from '$app/forms';
	import SettingsCard from '$lib/components/Settings/SettingsCard.svelte';
	import SettingsCardTitle from '$lib/components/Settings/SettingsCardTitle.svelte';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';

	const { data, form }: PageProps = $props();

	let submittingSave = $state(false);
	let submittingDelete = $state(false);
	const submitting = $derived(submittingSave || submittingDelete);
	let role = $state(data.organizationMembership.data.attributes.role);

	const canUpdateMembership = $derived(hasScope(UserScope.OrganizationMembershipsUpdate));
	const canDeleteMembership = $derived(hasScope(UserScope.OrganizationMembershipsDelete));
	const user = $derived(
		data.organizationMembership.included.filter((i) => i.type === API_TYPES.user)[0]
	);
	const canSubmit = $derived(role !== data.organizationMembership.data.attributes.role);
</script>

<div class="space-y-10">
	<button class="btn btn-ghost" onclick={() => history.back()}>
		<ArrowLeft class="size-4" />
		Back
	</button>

	<SettingsCard>
		<SettingsCardTitle>Member Information</SettingsCardTitle>
		<div class="space-y-2">
			<div class="flex gap-3">
				{#if user.meta.logoDistributionUrl}
					<img
						src={user.meta.logoDistributionUrl}
						alt={user.attributes.name}
						class="rounded-box size-20"
					/>
				{:else}
					<User class="rounded-box bg-base-300 size-20" />
				{/if}
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
				action="?/update"
				enctype="multipart/form-data"
				use:enhance={() => {
					submittingSave = true;

					return ({ update }) => {
						update();
						submittingSave = false;
					};
				}}
			>
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Role</legend>
					<select
						class="select select-bordered"
						name="role"
						bind:value={role}
						disabled={!canUpdateMembership}
					>
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

				<button class="btn btn-primary" disabled={!canSubmit || submitting || !canUpdateMembership}>
					{#if submittingSave}
						<span class="loading loading-spinner"></span>
					{:else}
						<Save class="size-4" />
					{/if}
					Save
				</button>
			</form>
		</div>
	</SettingsCard>

	<SettingsCard>
		<SettingsCardTitle>Danger Zone</SettingsCardTitle>
		<div>
			<form
				method="post"
				action="?/delete"
				use:enhance={() => {
					submittingDelete = true;

					return ({ update }) => {
						update();
						submittingDelete = false;
					};
				}}
			>
				<button class="btn btn-error btn-outline" disabled={submitting || !canDeleteMembership}>
					{#if submittingDelete}
						<span class="loading loading-spinner"></span>
					{:else}
						<Trash class="size-4" />
					{/if}
					Remove member
				</button>
			</form>
		</div>
	</SettingsCard>
</div>
