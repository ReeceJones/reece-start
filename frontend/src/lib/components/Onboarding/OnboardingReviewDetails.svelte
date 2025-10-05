<script lang="ts">
	import { localeToLanguageName } from '$lib/locale';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import { formatPhoneNumberWithCountryCode } from '$lib/phone-utils';

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
	<div class="space-y-6 rounded-box bg-base-200 p-6 text-base-content">
		<div class="space-y-4">
			<h3 class="text-xl font-semibold">Basic information</h3>
			<div class="space-y-4">
				{#if logoPreview}
					<img src={logoPreview} alt="Logo preview" class="aspect-square w-20 rounded-box" />
				{/if}
				<div class="space-y-2">
					<div>
						<p class="text-sm opacity-70">Name</p>
						<p class="font-medium">{display(onboardingState.name)}</p>
					</div>
					<div>
						<p class="text-sm opacity-70">Description</p>
						<p class="font-medium">{display(onboardingState.description)}</p>
					</div>
				</div>
			</div>
		</div>

		<div class="space-y-4">
			<h3 class="text-xl font-semibold">Contact</h3>
			<div class="grid grid-cols-1 gap-4">
				<div>
					<p class="text-sm opacity-70">Email</p>
					<p class="font-medium">{display(onboardingState.contactEmail)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">Phone</p>
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
			<h3 class="text-xl font-semibold">Address</h3>
			<div class="grid grid-cols-1 gap-4">
				<div>
					<p class="text-sm opacity-70">Address line 1</p>
					<p class="font-medium">{display(onboardingState.addressLine1)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">Address line 2</p>
					<p class="font-medium">{display(onboardingState.addressLine2)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">City</p>
					<p class="font-medium">{display(onboardingState.addressCity)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">State/Province</p>
					<p class="font-medium">{display(onboardingState.addressStateOrProvince)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">ZIP/Postal code</p>
					<p class="font-medium">{display(onboardingState.addressZip)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">Country</p>
					<p class="font-medium">{display(onboardingState.addressCountry)}</p>
				</div>
			</div>
		</div>

		<div class="space-y-4">
			<h3 class="text-xl font-semibold">Business details</h3>
			<div class="grid grid-cols-1 gap-4">
				<div>
					<p class="text-sm opacity-70">Organization type</p>
					<p class="font-medium">{display(onboardingState.entityType)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">Language</p>
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
