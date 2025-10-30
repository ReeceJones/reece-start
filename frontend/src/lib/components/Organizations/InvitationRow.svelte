<script lang="ts">
	import { enhance } from '$app/forms';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import type { OrganizationInvitation } from '$lib/schemas/organization-invitation';
	import { Check, Copy, Trash } from 'lucide-svelte';
	import { slide } from 'svelte/transition';
	import { t } from '$lib/i18n';

	const { invitation }: { invitation: OrganizationInvitation } = $props();

	const canDeleteInvitation = $derived(hasScope(UserScope.OrganizationInvitationsDelete));

	let submitting = $state(false);
	let copied = $state(false);
</script>

{#if copied}
	<div class="toast" transition:slide>
		<div class="alert">
			<Check class="size-4" />
			<span>{$t('invitationLinkCopied')}</span>
		</div>
	</div>
{/if}

<tr class="hover:bg-base-300">
	<td
		><a href={`mailto:${invitation.attributes.email}`} class="link">
			{invitation.attributes.email}
		</a>
	</td>
	<td>
		<button
			class="btn btn-ghost btn-sm"
			onclick={() => {
				const origin = window.location.origin;
				navigator.clipboard.writeText(`${origin}/app/invitations/${invitation.id}`);
				copied = true;
				setTimeout(() => {
					copied = false;
				}, 5000);
			}}
		>
			{#if copied}
				<Check class="size-4 transition-all" />
			{:else}
				<Copy class="size-4 transition-all" />
			{/if}
			<span> {$t('copyInvitationLink')} </span>
		</button>
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
				<button
					class="btn btn-square btn-ghost btn-sm btn-error"
					disabled={submitting || !canDeleteInvitation}
				>
					{#if submitting}
						<span class="loading size-4 loading-spinner"></span>
					{:else}
						<Trash class="size-4" />
					{/if}
				</button>
			</form>
		</div>
	</td>
</tr>
