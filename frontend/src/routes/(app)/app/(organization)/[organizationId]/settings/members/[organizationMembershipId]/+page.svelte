<script lang="ts">
	import { ArrowLeft, Save, Trash, User } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { API_TYPES } from '$lib/schemas/api';
	import { enhance } from '$app/forms';
	import * as Card from '$lib/components/ui/card';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { t } from '$lib/i18n';
	import { Button } from '$lib/components/ui/button';
	import * as Field from '$lib/components/ui/field';
	import LoadingIcon from '$lib/components/Icons/LoadingIcon.svelte';
	import * as Select from '$lib/components/ui/select';
	import FormActionStatus from '$lib/components/Form/FormActionStatus.svelte';

	const { data, form }: PageProps = $props();

	let submittingSave = $state(false);
	let submittingDelete = $state(false);
	const submitting = $derived(submittingSave || submittingDelete);
	let role = $derived(data.organizationMembership.data.attributes.role);

	const canUpdateMembership = $derived(hasScope(UserScope.OrganizationMembershipsUpdate));
	const canDeleteMembership = $derived(hasScope(UserScope.OrganizationMembershipsDelete));
	const user = $derived(
		data.organizationMembership.included.filter((i) => i.type === API_TYPES.user)[0]
	);
	const canSubmit = $derived(role !== data.organizationMembership.data.attributes.role);
</script>

<div class="space-y-10">
	<Button variant="ghost" onclick={() => history.back()}>
		<ArrowLeft class="size-4" />
		{$t('back')}
	</Button>

	<Card.Root>
		<Card.Header>
			<Card.Title>{$t('settings.organization.members.memberInformation')}</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-2">
			<div class="flex gap-3">
				{#if user.meta.logoDistributionUrl}
					<img
						src={user.meta.logoDistributionUrl}
						alt={user.attributes.name}
						class="size-20 rounded-lg"
					/>
				{:else}
					<User class="size-20 rounded-lg bg-neutral-200" />
				{/if}
				<div class="flex flex-col">
					<Card.Title>{user.attributes.name}</Card.Title>
					<a href={`mailto:${user.attributes.email}`} class="link text-sm text-gray-500">
						{user.attributes.email}
					</a>
				</div>
			</div>

			<form
				class="space-y-2 lg:max-w-sm"
				method="post"
				action="?/update"
				enctype="multipart/form-data"
				use:enhance={() => {
					submittingSave = true;

					return ({ update }) => {
						update({ reset: false });
						submittingSave = false;
					};
				}}
			>
				<Field.Field>
					<Field.Label for="member-role"
						>{$t('settings.organization.members.role.label')}</Field.Label
					>
					<Select.Root name="role" type="single" bind:value={role} disabled={!canUpdateMembership}>
						<Select.Trigger id="member-role">
							{$t(`settings.organization.members.role.${role}`)}
						</Select.Trigger>
						<Select.Content>
							<Select.Item value="admin"
								>{$t('settings.organization.members.role.admin')}</Select.Item
							>
							<Select.Item value="member"
								>{$t('settings.organization.members.role.member')}</Select.Item
							>
						</Select.Content>
					</Select.Root>
				</Field.Field>

				<FormActionStatus
					{form}
					success={$t('settings.organization.members.success.memberUpdated')}
					failure={$t('settings.organization.members.success.memberUpdateError')}
				/>

				<Button type="submit" disabled={!canSubmit || submitting || !canUpdateMembership}>
					<LoadingIcon loading={submitting}>
						{#snippet icon()}
							<Save class="size-4" />
						{/snippet}
					</LoadingIcon>
					{$t('save')}
				</Button>
			</form>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>{$t('settings.organization.members.dangerZone')}</Card.Title>
		</Card.Header>
		<Card.Content>
			<form
				method="post"
				action="?/delete"
				use:enhance={() => {
					submittingDelete = true;

					return ({ update }) => {
						update({ reset: false });
						submittingDelete = false;
					};
				}}
			>
				<Card.Action>
					<Button type="submit" variant="destructive" disabled={submitting || !canDeleteMembership}>
						<LoadingIcon loading={submitting}>
							{#snippet icon()}
								<Trash class="size-4" />
							{/snippet}
						</LoadingIcon>
						{$t('settings.organization.members.removeMember')}
					</Button>
				</Card.Action>
			</form>
		</Card.Content>
	</Card.Root>
</div>
