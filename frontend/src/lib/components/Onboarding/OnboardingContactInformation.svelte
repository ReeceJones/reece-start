<script lang="ts">
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import { getPhoneCodeOptions, type PhoneCodeOption } from '$lib/phone-utils';
	import * as m from '$lib/paraglide/messages';
	import * as Field from '$lib/components/ui/field';
	import { Input } from '$lib/components/ui/input';

	const {
		hidden,
		onboardingState = $bindable()
	}: { hidden: boolean; onboardingState: CreateOrganizationFormData } = $props();

	const phoneCodeOptions: PhoneCodeOption[] = getPhoneCodeOptions();
</script>

<OnboardingStepContainer {hidden}>
	<Field.Field>
		<Field.Label for="contactEmail"
			>{m.onboarding__contact_information_step__contact_email()}</Field.Label
		>
		<Input
			type="email"
			id="contactEmail"
			name="contactEmail"
			required
			class="input"
			placeholder={m.onboarding__contact_information_step__email()}
			bind:value={onboardingState.contactEmail}
		/>
		<Field.Description>{m.onboarding__contact_information_step__email_description()}</Field.Description
		>
	</Field.Field>
	<Field.Field>
		<Field.Label for="contactPhone"
			>{m.onboarding__contact_information_step__contact_phone()}</Field.Label
		>
		<div class="flex gap-2">
			<select
				name="contactPhoneCountry"
				id="contactPhoneCountry"
				class="select h-9 w-48 rounded-md border border-input bg-background px-3 py-1 text-base shadow-xs ring-offset-background transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30"
				bind:value={onboardingState.contactPhoneCountry}
			>
				{#each phoneCodeOptions as option (option.countryCode)}
					<option value={option.countryCode}>
						+{option.code} ({option.countryName}
						{option.flag})
					</option>
				{/each}
			</select>
			<Input
				type="tel"
				id="contactPhone"
				name="contactPhone"
				class="input flex-1"
				placeholder={m.onboarding__contact_information_step__phone_number()}
				bind:value={onboardingState.contactPhone}
			/>
		</div>
		<Field.Description>{m.onboarding__contact_information_step__phone_description()}</Field.Description
		>
	</Field.Field>
</OnboardingStepContainer>
