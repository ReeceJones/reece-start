<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import { Save, User } from 'lucide-svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Field from '$lib/components/ui/field';
	import * as m from '$lib/paraglide/messages';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import LogoInput from '$lib/components/Logo/LogoInput.svelte';
	import LogoPreview from '$lib/components/Logo/LogoPreview.svelte';
	import FormActionStatus from '$lib/components/Form/FormActionStatus.svelte';
	import LoadingIcon from '$lib/components/Icons/LoadingIcon.svelte';

	let { data, form }: PageProps = $props();

	let submitting = $state(false);
	let name = $derived(data.user.data.attributes.name);
	let email = $derived(data.user.data.attributes.email);
	let userProfile = $derived(data.user.data);
	let password = $state('');
	let confirmPassword = $state('');
	let logo = $state<FileList | null | undefined>(undefined);

	let canSubmit = $derived.by(() => {
		if (!name || !email) {
			return false;
		}

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
		<Card.Title>{m.profile__title()}</Card.Title>
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
				<Field.Label for="logo">{m.profile__profile_picture()}</Field.Label>
				<LogoPreview
					logoFile={logo}
					logoUrl={data.user.data.meta.logoDistributionUrl}
					alt="User logo"
				>
					{#snippet fallback()}
						<User class="size-32 text-neutral-600" />
					{/snippet}
				</LogoPreview>
				<LogoInput id="logo" name="logo" bind:logo />
				<Field.Description>{m.profile__upload_profile_picture()}</Field.Description>
			</Field.Field>

			<Field.Field>
				<Field.Label for="name">{m.profile__name()}</Field.Label>
				<Input
					type="text"
					id="name"
					name="name"
					required
					class="input"
					placeholder={m.profile__name_placeholder()}
					bind:value={name}
				/>
				<Field.Description>{m.profile__name_description()}</Field.Description>
			</Field.Field>

			<FormActionStatus
				{form}
				success={m.profile__profile_updated()}
				failure={m.profile__profile_update_error()}
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
