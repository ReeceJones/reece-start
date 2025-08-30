<script lang="ts">
	import { enhance } from '$app/forms';
	import type { OrganizationInvitation } from '$lib/schemas/organization-invitation';
	import { Trash } from 'lucide-svelte';

	const { invitation }: { invitation: OrganizationInvitation } = $props();

	let submitting = $state(false);
</script>

<tr class="hover:bg-base-300">
	<td
		><a href={`mailto:${invitation.attributes.email}`} class="link">
			{invitation.attributes.email}
		</a>
	</td>
	<td>
		<div class="flex items-center justify-end">
			<form
				method="POST"
				enctype="multipart/form-data"
				action={`/app/${invitation.relationships.organization.data.id}/settings/members?/deleteInvitation`}
				use:enhance={() => {
					submitting = true;
					return ({ update }) => {
						update();
						submitting = false;
					};
				}}
			>
				<input type="hidden" name="invitationId" value={invitation.id} />
				<button class="btn btn-ghost btn-error btn-sm btn-square" disabled={submitting}>
					{#if submitting}
						<span class="loading loading-spinner size-4"></span>
					{:else}
						<Trash class="size-4" />
					{/if}
				</button>
			</form>
		</div>
	</td>
</tr>
