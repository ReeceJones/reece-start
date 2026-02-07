<script lang="ts">
	import { enhance } from '$app/forms';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import type { OrganizationInvitation } from '$lib/schemas/organization-invitation';
	import { Check, Copy, Trash } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
	import { copyToClipboard } from '$lib/clipboard';
	import * as Table from '$lib/components/ui/table';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { toast } from 'svelte-sonner';
	import LoadingIcon from '../Icons/LoadingIcon.svelte';

	const { invitation }: { invitation: OrganizationInvitation } = $props();

	const canDeleteInvitation = $derived(hasScope(UserScope.OrganizationInvitationsDelete));

	let submitting = $state(false);
	let copied = $state(false);
</script>

<Table.Row class="hover:bg-base-300">
	<Table.Cell
		><a href={`mailto:${invitation.attributes.email}`} class={buttonVariants({ variant: 'link' })}>
			{invitation.attributes.email}
		</a>
	</Table.Cell>
	<Table.Cell>
		<Button
			variant="ghost"
			size="sm"
			onclick={async () => {
				const origin = window.location.origin;
				await copyToClipboard(`${origin}/app/invitations/${invitation.id}`);
				copied = true;
				setTimeout(() => {
					copied = false;
				}, 5000);
				toast.success(m.invitation_link_copied());
			}}
		>
			{#if copied}
				<Check class="size-4 transition-all" />
			{:else}
				<Copy class="size-4 transition-all" />
			{/if}
			<span> {m.copy_invitation_link()} </span>
		</Button>
	</Table.Cell>
	<Table.Cell>
		<div class="flex items-center justify-end">
			<form
				method="POST"
				enctype="multipart/form-data"
				action={`/app/${invitation.relationships.organization.data.id}/settings/members?/deleteInvitation`}
				use:enhance={() => {
					submitting = true;
					return ({ update }) => {
						update();
						console.log('update');
						submitting = false;
					};
				}}
			>
				<input type="hidden" name="invitationId" value={invitation.id} />
				<Button
					variant="ghost-destructive"
					size="icon-sm"
					type="submit"
					disabled={submitting || !canDeleteInvitation}
				>
					<LoadingIcon loading={submitting}>
						{#snippet icon()}
							<Trash class="size-4" />
						{/snippet}
					</LoadingIcon>
				</Button>
			</form>
		</div>
	</Table.Cell>
</Table.Row>
