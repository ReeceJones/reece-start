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
	import { CircleAlert, CircleX } from 'lucide-svelte';
	import type { LayoutProps } from './$types';

	const { children, data }: LayoutProps = $props();

	const steps = [
		{
			path: '/app/create-organization/basic-information',
			label: 'Organization Information',
			index: 0
		},
		{
			path: '/app/create-organization/contact-information',
			label: 'Contact Information',
			index: 1
		},
		{
			path: '/app/create-organization/address',
			label: 'Address',
			index: 2
		},
		{
			path: '/app/create-organization/business-details',
			label: 'Business Details',
			index: 3
		}
	];

	let submitting = $state(false);
	let canSubmit = $state<Record<number, boolean>>({
		0: false,
		1: true,
		2: true,
		3: false
	});
	let error = $state('');

	const activeStep = $derived(steps.find((step) => step.path === page.url.pathname));
	const userName = $derived(data.user.data.attributes.name);

	$effect(() => {
		if (!activeStep) {
			goto(steps[0].path);
		}
	});
</script>

<div class="space-y-6">
	<div class="space-y-2">
		<progress
			class="progress w-full transition-all duration-500"
			value={activeStep?.index ?? 0}
			max={steps.length - 1}
		></progress>
		<h2 class="text-xl font-semibold">{activeStep?.label}</h2>
	</div>
	<form
		id="create-organization-form"
		class="space-y-6"
		enctype="multipart/form-data"
		method="post"
		action="/app/create-organization/business-details"
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
		{@render children?.()}

		<OnboardingBasicInformation hidden={activeStep?.index !== 0} bind:canSubmit />
		<OnboardingContactInformation hidden={activeStep?.index !== 1} />
		<OnboardingAddress hidden={activeStep?.index !== 2} />
		<OnboardingBusinessDetails hidden={activeStep?.index !== 3} {userName} bind:canSubmit />

		{#if activeStep && activeStep.index === steps.length - 1 && Object.values(canSubmit).some((value) => !value)}
			<div role="alert" class="alert alert-info">
				<CircleAlert class="size-4" />
				<span
					>You must fill out the {Object.entries(canSubmit)
						.filter(([, value]) => !value)
						.map(([key]) => steps.find((step) => step.index === Number(key))?.label)
						.join(', ')} fields before creating your organization.</span
				>
			</div>
		{/if}

		{#if error}
			<div role="alert" class="alert alert-error">
				<CircleX />
				<span>{error}</span>
			</div>
		{/if}

		<div class="flex gap-3">
			{#if activeStep && activeStep.index > 0}
				<OnboardingBackButton disabled={submitting} step={steps[activeStep.index - 1]} />
			{/if}
			{#if activeStep && activeStep.index < steps.length - 1}
				<OnboardingNextButton disabled={submitting} step={steps[activeStep.index + 1]} />
			{/if}
			{#if activeStep && activeStep.index === steps.length - 1}
				<OnboardingCompleteButton
					disabled={!Object.values(canSubmit).every((value) => value)}
					loading={submitting}
				/>
			{/if}
		</div>
	</form>
</div>
