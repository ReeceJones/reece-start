<script lang="ts">
	import { enhance } from '$app/forms';
	import { CircleX, UserPlus, X } from 'lucide-svelte';

	const {
		organizationId,
		onMemberInvited
	}: { organizationId: string; onMemberInvited: (email: string) => void } = $props();

	let email = $state('');
	let role = $state('member');
	let submitting = $state(false);
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
>
	<UserPlus class="size-5" /> Add Member
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
		<h3 class="text-lg font-bold">Invite member</h3>
		<div class="my-6 space-y-3">
			<fieldset class="fieldset">
				<legend class="fieldset-legend">Email</legend>
				<input
					class="input w-full"
					type="email"
					bind:value={email}
					placeholder="Email"
					required
					name="email"
				/>
			</fieldset>
			<fieldset class="fieldset">
				<legend class="fieldset-legend">Role</legend>
				<div class="flex flex-col gap-2">
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
							<p class="text-neutral-content text-sm font-bold">Admin</p>
							<p>Manage organization settings and manage members</p>
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
							<p class="text-neutral-content text-sm font-bold">Member</p>
							<p>Manage XYZ</p>
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
				Close
			</button>
			<button class="btn btn-primary" disabled={submitting}>
				{#if submitting}
					<span class="loading loading-spinner loading-xs"></span>
				{:else}
					<UserPlus class="size-4" />
					Invite
				{/if}
			</button>
		</div>
	</form>
</dialog>
