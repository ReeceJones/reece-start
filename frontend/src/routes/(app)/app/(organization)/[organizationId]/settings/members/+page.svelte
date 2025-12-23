<script lang="ts">
	import { CircleCheck, SquarePen, User } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import InviteMember from '$lib/components/Organizations/InviteMember.svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import * as Table from '$lib/components/ui/table';
	import { Button } from '$lib/components/ui/button';
	import { Avatar, AvatarImage, AvatarFallback } from '$lib/components/ui/avatar';
	import InvitationRow from '$lib/components/Organizations/InvitationRow.svelte';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { t } from '$lib/i18n';
	import RoleBadge from '$lib/components/Organizations/RoleBadge.svelte';

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
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{$t('settings.organization.members.title')}</Card.Title>
	</Card.Header>
	<Card.Content>
		<InviteMember
			organizationId={params.organizationId}
			onMemberInvited={(email) => (invitedMemberEmail = email)}
		/>

		{#if invitedMemberEmail}
			<Alert.Root variant="success" class="my-1">
				<CircleCheck class="size-4" />
				<span>
					{$t('members.invitationSent')}
					<strong
						><a href={`mailto:${invitedMemberEmail}`} class="link">{invitedMemberEmail}</a></strong
					>
					{$t('members.withInstructionsToJoin')}
				</span>
			</Alert.Root>
		{/if}

		<Table.Root class="table">
			<Table.Header>
				<Table.Row>
					<Table.Head>{$t('members.name')}</Table.Head>
					<Table.Head>{$t('members.role')}</Table.Head>
					<Table.Head></Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#if memberships.length === 0}
					<Table.Row>
						<Table.Cell colspan={3} class="text-center"
							>{$t('members.noMembershipsFound')}</Table.Cell
						>
					</Table.Row>
				{/if}
				{#each memberships as membership (membership.membership.id)}
					<Table.Row class="hover:bg-base-300">
						<Table.Cell class="flex items-center gap-3">
							<Avatar class="size-10">
								{#if membership.user?.meta.logoDistributionUrl}
									<AvatarImage
										src={membership.user?.meta.logoDistributionUrl}
										alt={membership.user?.attributes.name}
									/>
								{/if}
								<AvatarFallback>
									<User class="size-6" />
								</AvatarFallback>
							</Avatar>
							<div class="flex flex-col">
								<div class="font-semibold">{membership.user?.attributes.name}</div>
								<a
									href={`mailto:${membership.user?.attributes.email}`}
									class="text-sm text-gray-500 hover:underline"
								>
									{membership.user?.attributes.email}
								</a>
							</div>
						</Table.Cell>
						<Table.Cell>
							<RoleBadge role={membership.membership.attributes.role} />
						</Table.Cell>
						<Table.Cell>
							<div class="flex items-center justify-end">
								<Button
									variant="ghost"
									size="icon-sm"
									href={canUpdateMembership
										? `/app/${params.organizationId}/settings/members/${membership.membership.id}`
										: undefined}
									disabled={!canUpdateMembership}
								>
									<SquarePen class="size-4" />
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card.Content>
</Card.Root>

<Card.Root>
	<Card.Header>
		<Card.Title>{$t('members.pendingInvitations')}</Card.Title>
	</Card.Header>
	<Card.Content>
		<Table.Root class="table">
			<Table.Header>
				<Table.Row>
					<Table.Head>{$t('members.email')}</Table.Head>
					<Table.Head></Table.Head>
					<Table.Head></Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#if data.invitations.data.length === 0}
					<Table.Row>
						<Table.Cell colspan={3} class="text-center"
							>{$t('members.noInvitationsFound')}</Table.Cell
						>
					</Table.Row>
				{/if}
				{#each data.invitations.data as invitation (invitation.id)}
					<InvitationRow {invitation} />
				{/each}
			</Table.Body>
		</Table.Root>
	</Card.Content>
</Card.Root>
