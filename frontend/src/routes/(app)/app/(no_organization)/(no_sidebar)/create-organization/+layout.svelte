<script lang="ts">
	import { enhance } from '$app/forms';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import OnboardingAddress from '$lib/components/Onboarding/OnboardingAddress.svelte';
	import OnboardingBackButton from '$lib/components/Onboarding/OnboardingBackButton.svelte';
	import OnboardingBasicInformation from '$lib/components/Onboarding/OnboardingBasicInformation.svelte';
	import OnboardingBusinessDetails from '$lib/components/Onboarding/OnboardingBusinessDetails.svelte';
	import OnboardingCompleteButton from '$lib/components/Onboarding/OnboardingCompleteButton.svelte';
	import OnboardingContactInformation from '$lib/components/Onboarding/OnboardingContactInformation.svelte';
	import OnboardingNextButton from '$lib/components/Onboarding/OnboardingNextButton.svelte';
	import { ArrowLeft, CircleX } from 'lucide-svelte';
	import type { LayoutProps } from './$types';
	import OnboardingReviewDetails from '$lib/components/Onboarding/OnboardingReviewDetails.svelte';
	import type { CreateOrganizationFormData } from '$lib/schemas/organization';
	import {
		isAddressValid,
		isBasicInformationValid,
		isBusinessDetailsValid,
		isContactInformationValid
	} from '$lib/organization-onboarding';
	import * as Card from '$lib/components/ui/card';
	import { t, locale } from '$lib/i18n';
	import { Progress } from '$lib/components/ui/progress';
	import { buttonVariants } from '$lib/components/ui/button';

	const { children, data }: LayoutProps = $props();

	const steps = [
		{
			path: '/app/create-organization/basic-information',
			label: $t('createOrganizationPages.steps.basicInformation'),
			index: 0
		},
		{
			path: '/app/create-organization/contact-information',
			label: $t('createOrganizationPages.steps.contactInformation'),
			index: 1
		},
		{
			path: '/app/create-organization/address',
			label: $t('createOrganizationPages.steps.address'),
			index: 2
		},
		{
			path: '/app/create-organization/business-details',
			label: $t('createOrganizationPages.steps.businessDetails'),
			index: 3
		},
		{
			path: '/app/create-organization/review',
			label: $t('createOrganizationPages.steps.review'),
			index: 4
		}
	];

	let submitting = $state(false);
	let onboardingState = $derived<CreateOrganizationFormData>({
		name: '',
		description: '',
		logo: undefined,
		locale: $locale,
		entityType: 'individual',
		addressCity: '',
		addressStateOrProvince: '',
		addressZip: '',
		addressCountry: 'US',
		addressLine1: '',
		addressLine2: '',
		contactEmail: data.user.data.attributes.email,
		contactPhone: '',
		contactPhoneCountry: 'US'
	});
	let error = $state('');

	const activeStep = $derived(steps.find((step) => step.path === page.url.pathname));
	const activeStepValid = $derived.by(() => {
		if (!activeStep) {
			return false;
		}

		switch (activeStep.path) {
			case '/app/create-organization/basic-information':
				return isBasicInformationValid(onboardingState);
			case '/app/create-organization/contact-information':
				return isContactInformationValid(onboardingState);
			case '/app/create-organization/address':
				return isAddressValid(onboardingState);
			case '/app/create-organization/business-details':
				return isBusinessDetailsValid(onboardingState);
			default:
				return true;
		}
	});

	$effect(() => {
		if (!activeStep) {
			goto(steps[0].path);
		}
	});
</script>

<a
	class={buttonVariants({
		variant: 'ghost',
		class: 'm-1'
	})}
	href="/app"
>
	<ArrowLeft class="size-4" />
	{$t('onboarding.back')}
</a>

<Card.Root class="mx-auto mb-10 max-w-2xl space-y-6">
	<Card.Header class="mb-0 space-y-3">
		<Card.Title
			>{activeStep?.label} ({activeStep?.index ? activeStep.index + 1 : 1} of {steps.length})</Card.Title
		>
		<Progress value={activeStep?.index ?? 0} max={steps.length - 1}></Progress>
		{@render children?.()}
	</Card.Header>
	<Card.Content>
		<form
			id="create-organization-form"
			class="space-y-6"
			enctype="multipart/form-data"
			method="post"
			action="/app/create-organization"
			use:enhance={() => {
				submitting = true;
				return ({ result, update }) => {
					update();
					submitting = false;

					if (result.type === 'failure') {
						error = result.data?.message as string;
					} else {
						error = '';
					}
				};
			}}
		>
			<OnboardingBasicInformation
				hidden={page.url.pathname !== '/app/create-organization/basic-information'}
				bind:onboardingState
			/>
			<OnboardingContactInformation
				hidden={page.url.pathname !== '/app/create-organization/contact-information'}
				bind:onboardingState
			/>
			<OnboardingAddress
				hidden={page.url.pathname !== '/app/create-organization/address'}
				bind:onboardingState
			/>
			<OnboardingBusinessDetails
				hidden={page.url.pathname !== '/app/create-organization/business-details'}
				bind:onboardingState
			/>
			<OnboardingReviewDetails
				hidden={page.url.pathname !== '/app/create-organization/review'}
				{onboardingState}
			/>

			{#if error}
				<div role="alert" class="alert alert-error">
					<CircleX />
					<span>{error}</span>
				</div>
			{/if}

			<div class="flex items-center gap-2">
				{#if activeStep && activeStep.index > 0}
					<OnboardingBackButton disabled={submitting} step={steps[activeStep.index - 1]} />
				{/if}
				{#if activeStep && activeStepValid && activeStep.index < steps.length - 1}
					<OnboardingNextButton disabled={submitting} step={steps[activeStep.index + 1]} />
				{/if}
				{#if activeStep && activeStep.index === steps.length - 1}
					<OnboardingCompleteButton loading={submitting} />
				{/if}
			</div>
		</form>
	</Card.Content>
</Card.Root>
