<script lang="ts">
	import { enhance } from '$app/forms';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { CircleX, UserPlus, X } from 'lucide-svelte';
	import { t } from '$lib/i18n';

	const {
		organizationId,
		onMemberInvited
	}: { organizationId: string; onMemberInvited: (email: string) => void } = $props();

	const canAddMember = $derived(hasScope(UserScope.OrganizationInvitationsCreate));

	let email = $state('');
	let submitting = $state(false);
	let role = $state('member');
	let error = $state('');

	let inviteMemberModal: HTMLDialogElement;
</script>

<button
	class="btn btn-primary w-fit"
	onclick={() => {
		inviteMemberModal.showModal();
		email = '';
		role = 'member';
		error = '';
	}}
	disabled={!canAddMember}
>
	<UserPlus class="size-5" />
	{$t('addMember')}
</button>

<dialog bind:this={inviteMemberModal} class="modal">
	<form
		class="modal-box"
		method="POST"
		action={`/app/${organizationId}/settings/members?/invite`}
		enctype="multipart/form-data"
		use:enhance={() => {
			submitting = true;

			return ({ result, update, formData }) => {
				console.log(result);
				update();
				submitting = false;

				if (result.type === 'failure') {
					error = result.data?.message as string;
				} else {
					inviteMemberModal.close();
					onMemberInvited(formData.get('email') as string);
				}
			};
		}}
	>
		<h3 class="text-lg font-bold">{$t('inviteMember')}</h3>
		<div class="my-6 space-y-3">
			<fieldset class="fieldset">
				<legend class="fieldset-legend">{$t('email')}</legend>
				<input
					class="input w-full"
					type="email"
					bind:value={email}
					placeholder={$t('email')}
					required
					name="email"
				/>
			</fieldset>
			<fieldset class="fieldset">
				<legend class="fieldset-legend">{$t('role')}</legend>
				<div class="flex flex-col gap-3">
					<label class="label">
						<input
							type="radio"
							required
							class="radio radio-sm bg-transparent"
							name="role"
							value="admin"
							bind:group={role}
						/>
						<div class="ml-3">
							<p class="text-base-content text-sm font-bold">{$t('roles.admin.title')}</p>
							<p>{$t('roles.admin.description')}</p>
						</div>
					</label>
					<label class="label">
						<input
							type="radio"
							required
							class="radio radio-sm bg-transparent"
							name="role"
							value="member"
							bind:group={role}
						/>
						<div class="ml-3">
							<p class="text-base-content text-sm font-bold">{$t('roles.member.title')}</p>
							<p>{$t('roles.member.description')}</p>
						</div>
					</label>
				</div>
			</fieldset>
		</div>
		{#if error}
			<div class="alert alert-error">
				<CircleX />
				<span>{error}</span>
			</div>
		{/if}
		<div class="modal-action">
			<button class="btn" type="button" onclick={() => inviteMemberModal.close()}>
				<X class="size-4" />
				{$t('close')}
			</button>
			<button class="btn btn-primary" disabled={submitting}>
				{#if submitting}
					<span class="loading loading-xs loading-spinner"></span>
				{:else}
					<UserPlus class="size-4" />
					{$t('invite')}
				{/if}
			</button>
		</div>
	</form>
</dialog>
