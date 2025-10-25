<script lang="ts">
	import { ArrowLeft, CircleCheck, CircleX, Save, Trash, User } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { API_TYPES } from '$lib/schemas/api';
	import { enhance } from '$app/forms';
	import SettingsCard from '$lib/components/Settings/SettingsCard.svelte';
	import SettingsCardTitle from '$lib/components/Settings/SettingsCardTitle.svelte';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
	import { t } from '$lib/i18n';

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
		{$t('back')}
	</button>

	<SettingsCard>
		<SettingsCardTitle>{$t('settings.organization.members.memberInformation')}</SettingsCardTitle>
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
					<CardTitle>{user.attributes.name}</CardTitle>
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
					<legend class="fieldset-legend">{$t('settings.organization.members.role.label')}</legend>
					<select
						class="select-bordered select"
						name="role"
						bind:value={role}
						disabled={!canUpdateMembership}
					>
						<option value="admin">{$t('settings.organization.members.role.admin')}</option>
						<option value="member">{$t('settings.organization.members.role.member')}</option>
					</select>
				</fieldset>

				{#if form?.success}
					<div role="alert" class="alert alert-success mt-3">
						<CircleCheck />
						<span>{$t('settings.organization.members.success.memberUpdated')}</span>
					</div>
				{:else if form?.success === false}
					<div role="alert" class="alert alert-error mt-3">
						<CircleX />
						<span
							>{form.message ?? $t('settings.organization.members.success.memberUpdateError')}</span
						>
					</div>
				{/if}

				<button class="btn btn-primary" disabled={!canSubmit || submitting || !canUpdateMembership}>
					{#if submittingSave}
						<span class="loading loading-spinner"></span>
					{:else}
						<Save class="size-4" />
					{/if}
					{$t('save')}
				</button>
			</form>
		</div>
	</SettingsCard>

	<SettingsCard>
		<SettingsCardTitle>{$t('settings.organization.members.dangerZone')}</SettingsCardTitle>
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
				<button class="btn btn-outline btn-error" disabled={submitting || !canDeleteMembership}>
					{#if submittingDelete}
						<span class="loading loading-spinner"></span>
					{:else}
						<Trash class="size-4" />
					{/if}
					{$t('settings.organization.members.removeMember')}
				</button>
			</form>
		</div>
	</SettingsCard>
</div>
