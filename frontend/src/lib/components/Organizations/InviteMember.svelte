<script lang="ts">
	import { enhance } from '$app/forms';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { CircleX, UserPlus, X } from 'lucide-svelte';
	import { t } from '$lib/i18n';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Field from '$lib/components/ui/field';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import * as Alert from '$lib/components/ui/alert';
	import type { Attachment } from 'svelte/attachments';
	import { Input } from '../ui/input';

	const {
		organizationId,
		onMemberInvited
	}: { organizationId: string; onMemberInvited: (email: string) => void } = $props();

	const canAddMember = $derived(hasScope(UserScope.OrganizationInvitationsCreate));

	let email = $state('');
	let submitting = $state(false);
	let role = $state('member');
	let error = $state('');
	let dialogOpen = $state(false);

	const resetFormAttachment: Attachment = () => {
		email = '';
		role = 'member';
		error = '';
	};
</script>

<Dialog.Root bind:open={dialogOpen}>
	<Dialog.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="default" class="w-fit" disabled={!canAddMember}>
				<UserPlus class="size-5" />
				{$t('addMember')}
			</Button>
		{/snippet}
	</Dialog.Trigger>
	<Dialog.Content>
		<form
			class="modal-box"
			method="POST"
			action={`/app/${organizationId}/settings/members?/invite`}
			enctype="multipart/form-data"
			use:enhance={() => {
				submitting = true;

				return ({ result, update, formData }) => {
					update();
					submitting = false;

					if (result.type === 'failure') {
						error = result.data?.message as string;
					} else {
						dialogOpen = false;
						onMemberInvited(formData.get('email') as string);
					}
				};
			}}
			{@attach resetFormAttachment}
		>
			<Dialog.Header>
				<Dialog.Title>
					{$t('inviteMember')}
				</Dialog.Title>
			</Dialog.Header>
			<div class="my-6 space-y-3">
				<Field.Field>
					<Field.Label for="invite-member-email">{$t('email')}</Field.Label>
					<Input
						class="input w-full"
						type="email"
						bind:value={email}
						placeholder={$t('email')}
						required
						name="email"
						id="invite-member-email"
					/>
				</Field.Field>
				<Field.Set>
					<Field.Label>{$t('role')}</Field.Label>
					<RadioGroup.Root required bind:value={role} name="role">
						<Field.Field orientation="horizontal">
							<RadioGroup.Item value="admin" id="invite-member-role-admin" />
							<div class="ml-3">
								<Field.Label for="invite-member-role-admin">{$t('roles.admin.title')}</Field.Label>
								<Field.Description>{$t('roles.admin.description')}</Field.Description>
							</div>
						</Field.Field>
						<Field.Field orientation="horizontal">
							<RadioGroup.Item value="member" id="invite-member-role-member" />
							<div class="ml-3">
								<Field.Label for="invite-member-role-member">{$t('roles.member.title')}</Field.Label
								>
								<Field.Description>{$t('roles.member.description')}</Field.Description>
							</div>
						</Field.Field>
					</RadioGroup.Root>
				</Field.Set>
			</div>
			{#if error}
				<Alert.Root>
					<CircleX />
					<Alert.Description>{error}</Alert.Description>
				</Alert.Root>
			{/if}
			<Dialog.Footer>
				<Button variant="outline" type="button" onclick={() => (dialogOpen = false)}>
					<X class="size-4" />
					{$t('close')}
				</Button>
				<Button variant="default" type="submit" disabled={submitting}>
					{#if submitting}
						<span class="loading loading-xs loading-spinner"></span>
					{:else}
						<UserPlus class="size-4" />
						{$t('invite')}
					{/if}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
