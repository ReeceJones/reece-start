<script lang="ts">
	import { localeToLanguageName } from '$lib/locale';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import { formatPhoneNumberWithCountryCode } from '$lib/phone-utils';
	import * as m from '$lib/paraglide/messages';

	const {
		hidden,
		onboardingState
	}: { hidden: boolean; onboardingState: CreateOrganizationFormData } = $props();

	const logoPreview = $derived(
		onboardingState.logo && onboardingState.logo.length > 0
			? URL.createObjectURL(onboardingState.logo[0])
			: undefined
	);

	function display(value: string | undefined | null): string {
		if (value === undefined || value === null || value === '') return '-';
		return value;
	}
</script>

<OnboardingStepContainer {hidden}>
	<div class="rounded-box bg-base-200 text-base-content space-y-6 p-6">
		<div class="space-y-4">
			<h3 class="text-xl font-semibold">{m.onboarding__review_details_step__basic_information()}</h3>
			<div class="space-y-4">
				{#if logoPreview}
					<img
						src={logoPreview}
						alt={m.onboarding__review_details_step__logo_preview()}
						class="rounded-box aspect-square w-20"
					/>
				{/if}
				<div class="space-y-2">
					<div>
						<p class="text-sm opacity-70">{m.onboarding__review_details_step__name()}</p>
						<p class="font-medium">{display(onboardingState.name)}</p>
					</div>
					<div>
						<p class="text-sm opacity-70">{m.onboarding__review_details_step__description()}</p>
						<p class="font-medium">{display(onboardingState.description)}</p>
					</div>
				</div>
			</div>
		</div>

		<div class="space-y-4">
			<h3 class="text-xl font-semibold">{m.onboarding__review_details_step__contact()}</h3>
			<div class="grid grid-cols-1 gap-4">
				<div>
					<p class="text-sm opacity-70">{m.onboarding__review_details_step__email()}</p>
					<p class="font-medium">{display(onboardingState.contactEmail)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{m.onboarding__review_details_step__phone()}</p>
					<p class="font-medium">
						{display(
							formatPhoneNumberWithCountryCode(
								onboardingState.contactPhone || '',
								onboardingState.contactPhoneCountry || ''
							)
						)}
					</p>
				</div>
			</div>
		</div>

		<div class="space-y-4">
			<h3 class="text-xl font-semibold">{m.onboarding__review_details_step__address()}</h3>
			<div class="grid grid-cols-1 gap-4">
				<div>
					<p class="text-sm opacity-70">{m.onboarding__review_details_step__address_line1()}</p>
					<p class="font-medium">{display(onboardingState.addressLine1)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{m.onboarding__review_details_step__address_line2()}</p>
					<p class="font-medium">{display(onboardingState.addressLine2)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{m.onboarding__review_details_step__city()}</p>
					<p class="font-medium">{display(onboardingState.addressCity)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{m.onboarding__review_details_step__state_province()}</p>
					<p class="font-medium">{display(onboardingState.addressStateOrProvince)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{m.onboarding__review_details_step__zip_postal_code()}</p>
					<p class="font-medium">{display(onboardingState.addressZip)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{m.onboarding__review_details_step__country()}</p>
					<p class="font-medium">{display(onboardingState.addressCountry)}</p>
				</div>
			</div>
		</div>

		<div class="space-y-4">
			<h3 class="text-xl font-semibold">{m.onboarding__review_details_step__business_details()}</h3>
			<div class="grid grid-cols-1 gap-4">
				<div>
					<p class="text-sm opacity-70">{m.onboarding__review_details_step__organization_type()}</p>
					<p class="font-medium">{display(onboardingState.entityType)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{m.onboarding__review_details_step__language()}</p>
					<p class="font-medium">
						<!-- TODO: add i8n here -->
						{onboardingState.locale
							? localeToLanguageName('en', onboardingState.locale)
							: display(onboardingState.locale)}
					</p>
				</div>
			</div>
		</div>
	</div>
</OnboardingStepContainer>
