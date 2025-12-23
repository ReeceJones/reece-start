<script lang="ts">
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import { t } from '$lib/i18n';
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
		<Field.Label for="name">{$t('onboarding.basicInformationStep.name')}</Field.Label>
		<Input
			type="text"
			id="name"
			name="name"
			required
			class="input"
			placeholder={$t('onboarding.basicInformationStep.name')}
			bind:value={onboardingState.name}
		/>
		<Field.Description>
			{$t('onboarding.basicInformationStep.nameDescription')}
		</Field.Description>
	</Field.Field>

	<Field.Field>
		<Field.Label for="description">{$t('onboarding.basicInformationStep.description')}</Field.Label>
		<Textarea
			id="description"
			name="description"
			class="min-h-[100px]"
			placeholder={$t('onboarding.basicInformationStep.description')}
			bind:value={onboardingState.description}
		/>
		<Field.Description
			>{$t('onboarding.basicInformationStep.descriptionDescription')}</Field.Description
		>
	</Field.Field>

	<Field.Field>
		<Field.Label for="logo">{$t('onboarding.basicInformationStep.logo')}</Field.Label>
		<LogoPreview
			logoFile={onboardingState.logo}
			logoUrl={undefined}
			alt={$t('onboarding.basicInformationStep.logoPreview')}
		>
			{#snippet fallback()}
				<Building2 class="size-32 text-neutral-600" />
			{/snippet}
		</LogoPreview>
		<LogoInput id="logo" name="logo" bind:logo={onboardingState.logo} />
		<Field.Description>{$t('onboarding.basicInformationStep.uploadLogo')}</Field.Description>
	</Field.Field>
</OnboardingStepContainer>
