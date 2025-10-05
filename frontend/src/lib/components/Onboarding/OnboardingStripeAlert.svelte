<script lang="ts">
	import type { Organization } from '$lib/schemas/organization';
	import { ExternalLink } from 'lucide-svelte';
	import { enhance } from '$app/forms';

	const { organization }: { organization: Organization } = $props();

	let loading = $state(false);
	let error = $state('');
</script>

{#if organization.data.meta.stripe.onboardingStatus === 'missing_requirements' || organization.data.meta.stripe.onboardingStatus === 'missing_capabilities'}
	<div class="mb-6 alert alert-warning">
		<div class="flex flex-col gap-2">
			<p class="font-semibold">
				To accept payments from your customers, Stripe needs more information about your business.
			</p>
			<form
				method="post"
				action={`/app/${organization.data.id}/stripe-onboarding`}
				use:enhance={() => {
					loading = true;
					error = '';
					return ({ result, update }) => {
						update();
						if (result.type === 'failure') {
							loading = false;
							error = (result.data?.message as string) ?? 'Something went wrong. Please try again.';
						} else if (result.type === 'error') {
							loading = false;
							error = result.error ?? 'Something went wrong. Please try again.';
						}
					};
				}}
			>
				<button class="btn btn-sm" disabled={loading}>
					{#if loading}
						<span class="loading loading-xs loading-spinner"></span>
					{:else}
						<ExternalLink class="size-4" />
					{/if}
					<span class="ml-1">Open Stripe</span>
				</button>
			</form>
			{#if error}
				<div class="mt-3 alert alert-error">
					<span>{error}</span>
				</div>
			{/if}
		</div>
	</div>
{/if}

{#if organization.data.meta.stripe.onboardingStatus === 'pending'}
	<div class="mb-6 alert alert-info">
		<div class="flex flex-col gap-2">
			<p class="font-semibold">
				We are setting up your Stripe account so that you can accept payments from your customers.
				Please check back later.
			</p>
		</div>
	</div>
{/if}
