<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import { Save, Building2 } from 'lucide-svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Field from '$lib/components/ui/field';
	import * as m from '$lib/paraglide/messages';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import LogoInput from '$lib/components/Logo/LogoInput.svelte';
	import LogoPreview from '$lib/components/Logo/LogoPreview.svelte';
	import FormActionStatus from '$lib/components/Form/FormActionStatus.svelte';
	import LoadingIcon from '$lib/components/Icons/LoadingIcon.svelte';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';

	let { data, form }: PageProps = $props();

	let submitting = $state(false);
	let logo = $state<FileList | null | undefined>(undefined);
	let name = $derived(data.organization.data.attributes.name);
	let description = $derived(data.organization.data.attributes.description);

	const isDirty = $derived(
		logo != null ||
			name !== data.organization.data.attributes.name ||
			description !== data.organization.data.attributes.description
	);
	const isValid = $derived(!!name);
	const canUpdate = $derived(hasScope(UserScope.OrganizationUpdate));
	const canSubmit = $derived(isDirty && isValid && canUpdate && !submitting);
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{m.settings__organization__general__title()}</Card.Title>
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
			<input type="hidden" tabindex="-1" name="organizationId" value={data.organization.data.id} />

			<Field.Field>
				<Field.Label for="logo">{m.settings__organization__general__logo__label()}</Field.Label>
				<LogoPreview
					logoFile={logo}
					logoUrl={data.organization.data.meta.logoDistributionUrl}
					alt="Organization logo"
				>
					{#snippet fallback()}
						<Building2 class="size-32 text-neutral-600" />
					{/snippet}
				</LogoPreview>
				<LogoInput id="logo" name="logo" bind:logo disabled={!canUpdate} />
				<Field.Description>{m.settings__organization__general__logo__description()}</Field.Description
				>
			</Field.Field>

			<Field.Field>
				<Field.Label for="name">{m.settings__organization__general__name__label()}</Field.Label>
				<Input
					type="text"
					id="name"
					name="name"
					required
					class="input"
					placeholder={m.settings__organization__general__name__placeholder()}
					bind:value={name}
					disabled={!canUpdate}
				/>
				<Field.Description>{m.settings__organization__general__name__description()}</Field.Description
				>
			</Field.Field>

			<Field.Field>
				<Field.Label for="description"
					>{m.settings__organization__general__description__label()}</Field.Label
				>
				<textarea
					id="description"
					name="description"
					class="textarea min-h-[100px] w-full rounded-md border border-input bg-background px-3 py-2 text-base shadow-xs ring-offset-background transition-[color,box-shadow] outline-none selection:bg-primary selection:text-primary-foreground placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30"
					placeholder={m.settings__organization__general__description__placeholder()}
					bind:value={description}
					maxlength={255}
					disabled={!canUpdate}
				></textarea>
				<Field.Description
					>{m.settings__organization__general__description__description()}</Field.Description
				>
			</Field.Field>

			<FormActionStatus
				{form}
				success={m.settings__organization__general__success__organization_updated()}
				failure={form?.message ??
					m.settings__organization__general__success__organization_update_error()}
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
