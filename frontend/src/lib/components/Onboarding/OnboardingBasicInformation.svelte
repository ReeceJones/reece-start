<script lang="ts">
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import * as m from '$lib/paraglide/messages';
	import * as Field from '$lib/components/ui/field';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Building2 } from 'lucide-svelte';
	import LogoPreview from '$lib/components/Logo/LogoPreview.svelte';
	import LogoInput from '$lib/components/Logo/LogoInput.svelte';

	const {
		hidden,
		onboardingState = $bindable()
	}: {
		hidden: boolean;
		onboardingState: CreateOrganizationFormData;
	} = $props();
</script>

<OnboardingStepContainer {hidden}>
	<Field.Field>
		<Field.Label for="name">{m.onboarding__basic_information_step__name()}</Field.Label>
		<Input
			type="text"
			id="name"
			name="name"
			required
			class="input"
			placeholder={m.onboarding__basic_information_step__name()}
			bind:value={onboardingState.name}
		/>
		<Field.Description>
			{m.onboarding__basic_information_step__name_description()}
		</Field.Description>
	</Field.Field>

	<Field.Field>
		<Field.Label for="description">{m.onboarding__basic_information_step__description()}</Field.Label>
		<Textarea
			id="description"
			name="description"
			class="min-h-[100px]"
			placeholder={m.onboarding__basic_information_step__description()}
			bind:value={onboardingState.description}
		/>
		<Field.Description
			>{m.onboarding__basic_information_step__description_description()}</Field.Description
		>
	</Field.Field>

	<Field.Field>
		<Field.Label for="logo">{m.onboarding__basic_information_step__logo()}</Field.Label>
		<LogoPreview
			logoFile={onboardingState.logo}
			logoUrl={undefined}
			alt={m.onboarding__basic_information_step__logo_preview()}
		>
			{#snippet fallback()}
				<Building2 class="size-32 text-neutral-600" />
			{/snippet}
		</LogoPreview>
		<LogoInput id="logo" name="logo" bind:logo={onboardingState.logo} />
		<Field.Description>{m.onboarding__basic_information_step__upload_logo()}</Field.Description>
	</Field.Field>
</OnboardingStepContainer>
