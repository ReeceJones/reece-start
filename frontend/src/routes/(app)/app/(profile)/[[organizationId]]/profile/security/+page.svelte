<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import { Save } from 'lucide-svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Field from '$lib/components/ui/field';
	import * as m from '$lib/paraglide/messages';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import FormActionStatus from '$lib/components/Form/FormActionStatus.svelte';
	import LoadingIcon from '$lib/components/Icons/LoadingIcon.svelte';

	let { data, form }: PageProps = $props();

	let submitting = $state(false);
	let email = $derived(data.user.data.attributes.email);
	let userProfile = $derived(data.user.data);
	let password = $state('');
	let confirmPassword = $state('');

	let canSubmit = $derived.by(() => {
		if (password !== '' && password.length < 8) {
			return false;
		}

		if (password !== '' && password !== confirmPassword) {
			return false;
		}

		return true;
	});
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{m.settings__security()}</Card.Title>
	</Card.Header>
	<Card.Content>
		<form
			method="post"
			use:enhance={() => {
				submitting = true;

				return ({ update }) => {
					update({ reset: false });
					submitting = false;
				};
			}}
			enctype="multipart/form-data"
			class="space-y-4 lg:max-w-sm"
		>
			<input type="hidden" tabindex="-1" name="userId" value={userProfile.id} />

			<Field.Field>
				<Field.Label for="email">{m.settings__fields__email__label()}</Field.Label>
				<Input
					type="email"
					id="email"
					name="email"
					required
					class="input"
					placeholder={m.settings__fields__email__placeholder()}
					bind:value={email}
				/>
				<Field.Description>{m.settings__fields__email__description()}</Field.Description>
			</Field.Field>

			<Field.Field>
				<Field.Label for="password">{m.settings__fields__update_password__label()}</Field.Label>
				<Input
					type="password"
					id="password"
					name="password"
					class="input"
					placeholder={m.settings__fields__update_password__placeholder()}
					bind:value={password}
				/>
				<Field.Description>{m.settings__fields__update_password__description()}</Field.Description>
				{#if password.length > 0 && password.length < 8}
					<Field.Description class="text-error">
						{m.settings__fields__update_password__password_too_short()}
					</Field.Description>
				{/if}
			</Field.Field>

			{#if password !== ''}
				<Field.Field>
					<Field.Label for="confirmPassword"
						>{m.settings__fields__confirm_password__label()}</Field.Label
					>
					<Input
						type="password"
						id="confirmPassword"
						name="confirmPassword"
						class="input"
						aria-invalid={password !== confirmPassword && confirmPassword !== ''}
						placeholder={m.settings__fields__confirm_password__placeholder()}
						bind:value={confirmPassword}
					/>
					{#if password !== confirmPassword && confirmPassword !== ''}
						<Field.Description class="text-error">
							{m.settings__fields__confirm_password__password_does_not_match()}
						</Field.Description>
					{/if}
				</Field.Field>
			{/if}

			<FormActionStatus
				{form}
				success={m.settings__success__profile_updated()}
				failure={form?.message ?? m.settings__success__profile_update_error()}
			/>

			<Card.Action>
				<Button type="submit" disabled={!canSubmit || submitting}>
					<LoadingIcon loading={submitting}>
						{#snippet icon()}
							<Save />
						{/snippet}
					</LoadingIcon>
					<span>{m.save()}</span>
				</Button>
			</Card.Action>
		</form>
	</Card.Content>
</Card.Root>
