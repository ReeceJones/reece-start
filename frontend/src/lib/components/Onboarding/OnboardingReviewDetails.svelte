<script lang="ts">
	import { localeToLanguageName } from '$lib/locale';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import OnboardingStepContainer from './OnboardingStepContainer.svelte';
	import { formatPhoneNumberWithCountryCode } from '$lib/phone-utils';
	import { t } from '$lib/i18n';

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
			<h3 class="text-xl font-semibold">{$t('onboarding.reviewDetailsStep.basicInformation')}</h3>
			<div class="space-y-4">
				{#if logoPreview}
					<img
						src={logoPreview}
						alt={$t('onboarding.reviewDetailsStep.logoPreview')}
						class="aspect-square w-20 rounded-box"
					/>
				{/if}
				<div class="space-y-2">
					<div>
						<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.name')}</p>
						<p class="font-medium">{display(onboardingState.name)}</p>
					</div>
					<div>
						<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.description')}</p>
						<p class="font-medium">{display(onboardingState.description)}</p>
					</div>
				</div>
			</div>
		</div>

		<div class="space-y-4">
			<h3 class="text-xl font-semibold">{$t('onboarding.reviewDetailsStep.contact')}</h3>
			<div class="grid grid-cols-1 gap-4">
				<div>
					<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.email')}</p>
					<p class="font-medium">{display(onboardingState.contactEmail)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.phone')}</p>
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
			<h3 class="text-xl font-semibold">{$t('onboarding.reviewDetailsStep.address')}</h3>
			<div class="grid grid-cols-1 gap-4">
				<div>
					<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.addressLine1')}</p>
					<p class="font-medium">{display(onboardingState.addressLine1)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.addressLine2')}</p>
					<p class="font-medium">{display(onboardingState.addressLine2)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.city')}</p>
					<p class="font-medium">{display(onboardingState.addressCity)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.stateProvince')}</p>
					<p class="font-medium">{display(onboardingState.addressStateOrProvince)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.zipPostalCode')}</p>
					<p class="font-medium">{display(onboardingState.addressZip)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.country')}</p>
					<p class="font-medium">{display(onboardingState.addressCountry)}</p>
				</div>
			</div>
		</div>

		<div class="space-y-4">
			<h3 class="text-xl font-semibold">{$t('onboarding.reviewDetailsStep.businessDetails')}</h3>
			<div class="grid grid-cols-1 gap-4">
				<div>
					<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.organizationType')}</p>
					<p class="font-medium">{display(onboardingState.entityType)}</p>
				</div>
				<div>
					<p class="text-sm opacity-70">{$t('onboarding.reviewDetailsStep.language')}</p>
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
