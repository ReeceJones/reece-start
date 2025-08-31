<script lang="ts">
	import { API_TYPES } from '$lib/schemas/api';
	import type { OrganizationData } from '$lib/schemas/organization';
	import type { UserData } from '$lib/schemas/user';
	import { Check, CircleX, X } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';

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

<div class="card bg-base-200 border-base-300 mx-auto max-w-[600px] shadow-sm">
	<div class="card-body">
		{#if status === 'pending'}
			<div class="flex flex-col items-center justify-center gap-6">
				<div class="flex flex-col items-center justify-center gap-1">
					<h1 class="card-title text-center">
						{invitingUser?.attributes.name} invited you to join "{organization?.attributes.name}"
					</h1>
					<p class="text-center text-gray-500">
						By accepting, you will be added to the organization, and you will be able to collaborate
						with your team.
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
						<button class="btn btn-neutral btn-md" disabled={submitting}>
							<X class="size-4" />
							Decline
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
						<button class="btn btn-primary btn-md" disabled={submitting}>
							<Check class="size-4" />
							Accept
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
				<h1 class="card-title text-center">This invitation has already been accepted.</h1>
				<p class="text-center text-gray-500">
					If you did not accept this invitation, please contact the organization owner for a new
					invitation.
				</p>
			</div>
		{:else if status === 'declined'}
			<div class="flex flex-col items-center justify-center gap-6">
				<h1 class="card-title text-center">This invitation has already been declined.</h1>
				<p class="text-center text-gray-500">
					If you would like to join this organization, please contact the organization owner for a
					new invitation.
				</p>
			</div>
		{:else if status === 'expired'}
			<div class="flex flex-col items-center justify-center gap-6">
				<h1 class="card-title text-center">This invitation has expired.</h1>
				<p class="text-center text-gray-500">
					If you would like to join this organization, please contact the organization owner for a
					new invitation.
				</p>
			</div>
		{:else if status === 'revoked'}
			<div class="flex flex-col items-center justify-center gap-6">
				<h1 class="card-title text-center">This invitation has already been revoked.</h1>
				<p>
					If you would like to join this organization, please contact the organization owner for a
					new invitation.
				</p>
			</div>
		{/if}
	</div>
</div>
