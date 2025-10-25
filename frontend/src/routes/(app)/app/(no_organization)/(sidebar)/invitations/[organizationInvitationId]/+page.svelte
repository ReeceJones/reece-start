<script lang="ts">
	import { API_TYPES } from '$lib/schemas/api';
	import type { OrganizationData } from '$lib/schemas/organization';
	import type { UserData } from '$lib/schemas/user';
	import { Check, CircleX, X } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
	import { t } from '$lib/i18n';

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

<Card class="mx-auto max-w-[600px]">
	<CardBody>
		{#if status === 'pending'}
			<div class="flex flex-col items-center justify-center gap-6">
				<div class="flex flex-col items-center justify-center gap-1">
					<CardTitle class="text-center">
						{$t('noOrganization.invitation.invitedBy', {
							inviterName: invitingUser?.attributes.name || '',
							organizationName: organization?.attributes.name || ''
						})}
					</CardTitle>
					<p class="text-center text-gray-500">
						{$t('noOrganization.invitation.invitationDescription')}
					</p>
				</div>

				{#if organization?.meta.logoDistributionUrl}
					<img
						src={organization?.meta.logoDistributionUrl}
						alt={organization?.attributes.name}
						class="rounded-box size-32"
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
						<button class="btn btn-md btn-neutral" disabled={submitting}>
							<X class="size-4" />
							{$t('noOrganization.invitation.decline')}
						</button>
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
						<button class="btn btn-md btn-primary" disabled={submitting}>
							<Check class="size-4" />
							{$t('noOrganization.invitation.accept')}
						</button>
					</form>
				</div>
				{#if form && !form?.success}
					<div class="alert alert-error">
						<CircleX class="size-5" />
						<span>{form.message}</span>
					</div>
				{/if}
			</div>
		{:else if status === 'accepted'}
			<div class="flex flex-col items-center justify-center gap-6">
				<CardTitle class="text-center">{$t('noOrganization.invitation.accepted.title')}</CardTitle>
				<p class="text-center text-gray-500">
					{$t('noOrganization.invitation.accepted.description')}
				</p>
			</div>
		{:else if status === 'declined'}
			<div class="flex flex-col items-center justify-center gap-6">
				<CardTitle class="text-center">{$t('noOrganization.invitation.declined.title')}</CardTitle>
				<p class="text-center text-gray-500">
					{$t('noOrganization.invitation.declined.description')}
				</p>
			</div>
		{:else if status === 'expired'}
			<div class="flex flex-col items-center justify-center gap-6">
				<CardTitle class="text-center">{$t('noOrganization.invitation.expired.title')}</CardTitle>
				<p class="text-center text-gray-500">
					{$t('noOrganization.invitation.expired.description')}
				</p>
			</div>
		{:else if status === 'revoked'}
			<div class="flex flex-col items-center justify-center gap-6">
				<CardTitle class="text-center">{$t('noOrganization.invitation.revoked.title')}</CardTitle>
				<p class="text-center text-gray-500">
					{$t('noOrganization.invitation.revoked.description')}
				</p>
			</div>
		{/if}
	</CardBody>
</Card>
