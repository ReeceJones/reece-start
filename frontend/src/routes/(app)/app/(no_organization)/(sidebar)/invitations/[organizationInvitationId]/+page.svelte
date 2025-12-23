<script lang="ts">
	import { API_TYPES } from '$lib/schemas/api';
	import type { OrganizationData } from '$lib/schemas/organization';
	import type { UserData } from '$lib/schemas/user';
	import { Check, CircleX, X } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import * as Card from '$lib/components/ui/card';
	import { t } from '$lib/i18n';
	import { Button } from '$lib/components/ui/button';
	import LoadingIcon from '$lib/components/Icons/LoadingIcon.svelte';

	const { data, form }: PageProps = $props();
	const organization = $derived(
		data.invitation.included.find(
			(i): i is OrganizationData =>
				i.type === API_TYPES.organization &&
				i.id === data.invitation.data.relationships.organization.data.id
		)
	);
	const invitingUser = $derived(
		data.invitation.included.find(
			(i): i is UserData =>
				i.type === API_TYPES.user &&
				i.id === data.invitation.data.relationships.invitingUser.data.id
		)
	);
	const status = $derived(data.invitation.data.attributes.status);

	let submitting = $state(false);
</script>

<Card.Root class="mx-auto max-w-[600px]">
	{#if status === 'pending'}
		<Card.Header>
			<div class="flex flex-col items-center justify-center gap-1">
				<Card.Title class="text-center">
					{$t('noOrganization.invitation.invitedBy', {
						inviterName: invitingUser?.attributes.name || '',
						organizationName: organization?.attributes.name || ''
					})}
				</Card.Title>
				<Card.Description class="text-center">
					{$t('noOrganization.invitation.invitationDescription')}
				</Card.Description>
			</div>
		</Card.Header>
		<Card.Content>
			<div class="flex flex-col items-center justify-center gap-6">
				{#if organization?.meta.logoDistributionUrl}
					<img
						src={organization?.meta.logoDistributionUrl}
						alt={organization?.attributes.name}
						class="size-32 rounded-lg"
					/>
				{/if}

				<div class="flex flex-row gap-3">
					<form
						method="POST"
						enctype="multipart/form-data"
						action="?/decline"
						use:enhance={() => {
							submitting = true;
							return ({ update }) => {
								update();
								submitting = false;
							};
						}}
					>
						<Button type="submit" variant="secondary" disabled={submitting}>
							<LoadingIcon loading={submitting}>
								{#snippet icon()}
									<X class="size-4" />
								{/snippet}
							</LoadingIcon>
							{$t('noOrganization.invitation.decline')}
						</Button>
					</form>
					<form
						method="POST"
						enctype="multipart/form-data"
						action="?/accept"
						use:enhance={() => {
							submitting = true;
							return ({ update }) => {
								update();
								submitting = false;
							};
						}}
					>
						<Button type="submit" disabled={submitting}>
							<LoadingIcon loading={submitting}>
								{#snippet icon()}
									<Check class="size-4" />
								{/snippet}
							</LoadingIcon>
							{$t('noOrganization.invitation.accept')}
						</Button>
					</form>
				</div>
				{#if form && !form?.success}
					<div class="alert alert-error">
						<CircleX class="size-5" />
						<span>{form.message}</span>
					</div>
				{/if}
			</div>
		</Card.Content>
	{:else if status === 'accepted'}
		<Card.Header>
			<div class="flex flex-col items-center justify-center gap-6">
				<Card.Title class="text-center">{$t('noOrganization.invitation.accepted.title')}</Card.Title
				>
				<Card.Description class="text-center">
					{$t('noOrganization.invitation.accepted.description')}
				</Card.Description>
			</div>
		</Card.Header>
	{:else if status === 'declined'}
		<Card.Header>
			<div class="flex flex-col items-center justify-center gap-6">
				<Card.Title class="text-center">{$t('noOrganization.invitation.declined.title')}</Card.Title
				>
				<Card.Description class="text-center">
					{$t('noOrganization.invitation.declined.description')}
				</Card.Description>
			</div>
		</Card.Header>
	{:else if status === 'expired'}
		<Card.Header>
			<div class="flex flex-col items-center justify-center gap-6">
				<Card.Title class="text-center">{$t('noOrganization.invitation.expired.title')}</Card.Title>
				<Card.Description class="text-center">
					{$t('noOrganization.invitation.expired.description')}
				</Card.Description>
			</div>
		</Card.Header>
	{:else if status === 'revoked'}
		<Card.Header>
			<div class="flex flex-col items-center justify-center gap-6">
				<Card.Title class="text-center">{$t('noOrganization.invitation.revoked.title')}</Card.Title>
				<Card.Description class="text-center">
					{$t('noOrganization.invitation.revoked.description')}
				</Card.Description>
			</div>
		</Card.Header>
	{/if}
</Card.Root>
