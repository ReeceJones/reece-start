<script lang="ts">
	import CountryOptions from '../CountryOptions.svelte';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import StateOptions from '../StateOptions.svelte';
	import { State } from 'country-state-city';
	import { t } from '$lib/i18n';
	import * as Field from '$lib/components/ui/field';
	import { Input } from '$lib/components/ui/input';

	const {
		hidden,
		onboardingState = $bindable()
	}: { hidden: boolean; onboardingState: CreateOrganizationFormData } = $props();
</script>

<OnboardingStepContainer {hidden}>
	<Field.Field>
		<Field.Label for="addressCountry">{$t('onboarding.addressStep.country')}</Field.Label>
		<select
			id="addressCountry"
			name="addressCountry"
			required
			class="h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-base shadow-xs ring-offset-background transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30"
			bind:value={onboardingState.addressCountry}
		>
			<CountryOptions />
		</select>
		<Field.Description>{$t('onboarding.addressStep.selectCountry')}</Field.Description>
	</Field.Field>

	<Field.Field>
		<Field.Label for="addressLine1">{$t('onboarding.addressStep.address')}</Field.Label>
		<Input
			type="text"
			id="addressLine1"
			name="addressLine1"
			required
			class="input"
			placeholder={$t('onboarding.addressStep.address')}
			bind:value={onboardingState.addressLine1}
		/>
		<Field.Description>{$t('onboarding.addressStep.enterStreetAddress')}</Field.Description>
	</Field.Field>

	<Field.Field>
		<Field.Label for="addressLine2">{$t('onboarding.addressStep.addressLine2')}</Field.Label>
		<Input
			type="text"
			id="addressLine2"
			name="addressLine2"
			class="input"
			placeholder={$t('onboarding.addressStep.addressLine2')}
			bind:value={onboardingState.addressLine2}
		/>
		<Field.Description>
			{$t('onboarding.addressStep.addressLine2Description')}
		</Field.Description>
	</Field.Field>

	<Field.Field>
		<Field.Label for="addressCity">{$t('onboarding.addressStep.city')}</Field.Label>
		<Input
			type="text"
			id="addressCity"
			name="addressCity"
			required
			class="input"
			placeholder={$t('onboarding.addressStep.city')}
			bind:value={onboardingState.addressCity}
		/>
		<Field.Description>{$t('onboarding.addressStep.enterCity')}</Field.Description>
	</Field.Field>

	{#if onboardingState.addressCountry && State.getStatesOfCountry(onboardingState.addressCountry).length > 0}
		<Field.Field>
			<Field.Label for="addressStateOrProvince">{$t('onboarding.addressStep.state')}</Field.Label>
			<select
				id="addressStateOrProvince"
				name="addressStateOrProvince"
				class="h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-base shadow-xs ring-offset-background transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30"
				bind:value={onboardingState.addressStateOrProvince}
			>
				{#if onboardingState.addressCountry}
					<StateOptions countryCode={onboardingState.addressCountry} />
				{/if}
			</select>
			<Field.Description>{$t('onboarding.addressStep.selectStateOrProvince')}</Field.Description>
		</Field.Field>
	{/if}

	<Field.Field>
		<Field.Label for="addressZip">{$t('onboarding.addressStep.zip')}</Field.Label>
		<Input
			type="text"
			id="addressZip"
			name="addressZip"
			required
			class="input"
			placeholder={$t('onboarding.addressStep.zip')}
			bind:value={onboardingState.addressZip}
		/>
		<Field.Description>{$t('onboarding.addressStep.enterZip')}</Field.Description>
	</Field.Field>
</OnboardingStepContainer>
