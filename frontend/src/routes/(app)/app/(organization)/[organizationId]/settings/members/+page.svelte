<script lang="ts">
	import { CircleCheck, SquarePen, User } from 'lucide-svelte';
	import clsx, { type ClassValue } from 'clsx/lite';
	import type { PageProps } from './$types';
	import InviteMember from '$lib/components/Organizations/InviteMember.svelte';
	import type { OrganizationMembershipRole } from '$lib/schemas/organization-membership';
	import SettingsCard from '$lib/components/Settings/SettingsCard.svelte';
	import SettingsCardTitle from '$lib/components/Settings/SettingsCardTitle.svelte';
	import InvitationRow from '$lib/components/Organizations/InvitationRow.svelte';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';

	const { data, params }: PageProps = $props();
	let invitedMemberEmail = $state('');

	const canUpdateMembership = $derived(hasScope(UserScope.OrganizationMembershipsUpdate));
	const memberships = $derived.by(() => {
		return data.memberships.data.map((membership) => ({
			membership,
			user: data.memberships.included.find(
				(user) => user.id === membership.relationships.user.data.id
			)
		}));
	});

	function getBadgeColorForRole(role: OrganizationMembershipRole): ClassValue {
		switch (role) {
			case 'admin':
				return 'badge-primary';
			case 'member':
				return 'badge-neutral';
			default:
				return 'badge-neutral';
		}
	}
</script>

<SettingsCard>
	<SettingsCardTitle>Members</SettingsCardTitle>
	<InviteMember
		organizationId={params.organizationId}
		onMemberInvited={(email) => (invitedMemberEmail = email)}
	/>

	{#if invitedMemberEmail}
		<div class="my-1 alert alert-success">
			<CircleCheck class="size-4" />
			<span>
				We've sent an email to
				<strong
					><a href={`mailto:${invitedMemberEmail}`} class="link">{invitedMemberEmail}</a></strong
				> with instructions to join your organization.
			</span>
		</div>
	{/if}

	<div class="overflow-auto">
		<table class="table">
			<thead>
				<tr>
					<th>Name</th>
					<th>Role</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#if memberships.length === 0}
					<tr>
						<td colspan="3" class="text-center">No memberships found</td>
					</tr>
				{/if}
				{#each memberships as membership (membership.membership.id)}
					<tr class="hover:bg-base-300">
						<td class="flex items-center gap-3">
							{#if membership.user?.meta.logoDistributionUrl}
								<img
									src={membership.user?.meta.logoDistributionUrl}
									alt={membership.user?.attributes.name}
									class="size-10 rounded-box"
								/>
							{:else}
								<User class="size-10 rounded-box bg-base-300" />
							{/if}
							<div class="flex flex-col">
								<div class="font-semibold">{membership.user?.attributes.name}</div>
								<a
									href={`mailto:${membership.user?.attributes.email}`}
									class="link text-sm text-gray-500"
								>
									{membership.user?.attributes.email}
								</a>
							</div>
						</td>
						<td>
							<div
								class={clsx('badge', getBadgeColorForRole(membership.membership.attributes.role))}
							>
								{membership.membership.attributes.role.charAt(0).toUpperCase() +
									membership.membership.attributes.role.slice(1)}
							</div>
						</td>
						<td>
							<div class="flex items-center justify-end">
								<a
									class={clsx(
										'btn btn-square btn-ghost btn-sm',
										!canUpdateMembership &&
											'pointer-events-none cursor-default text-base-content/50'
									)}
									href={canUpdateMembership
										? `/app/${params.organizationId}/settings/members/${membership.membership.id}`
										: undefined}
									aria-disabled={!canUpdateMembership}
								>
									<SquarePen class="size-4" />
								</a>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</SettingsCard>

<SettingsCard>
	<SettingsCardTitle>Pending Invitations</SettingsCardTitle>
	<div class="overflow-auto">
		<table class="table">
			<thead>
				<tr>
					<th>Email</th>
					<th></th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#if data.invitations.data.length === 0}
					<tr>
						<td colspan="2" class="text-center">No invitations found</td>
					</tr>
				{/if}
				{#each data.invitations.data as invitation (invitation.id)}
					<InvitationRow {invitation} />
				{/each}
			</tbody>
		</table>
	</div>
</SettingsCard>
