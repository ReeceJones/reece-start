<script lang="ts">
	import SettingsCard from '$lib/components/Settings/SettingsCard.svelte';
	import SettingsCardTitle from '$lib/components/Settings/SettingsCardTitle.svelte';
	import SettingsCardActions from '$lib/components/Settings/SettingsCardActions.svelte';
	import { post } from '$lib';
	import {
		createCheckoutSessionRequestSchema,
		createCheckoutSessionResponseSchema,
		createBillingPortalSessionRequestSchema,
		createBillingPortalSessionResponseSchema,
		type CreateCheckoutSessionRequest,
		type CreateBillingPortalSessionRequest
	} from '$lib/schemas/stripe';
	import { CreditCard, CheckCircle2, Sparkles } from 'lucide-svelte';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import type { PageProps } from './$types';

	const { data }: PageProps = $props();

	let loading = $state(false);
	let error = $state<string | null>(null);

	const canManageBilling = $derived(hasScope(UserScope.OrganizationBillingUpdate));
	const isProPlan = $derived(data.subscription.data.attributes.plan === 'pro');
	const isFreePlan = $derived(data.subscription.data.attributes.plan === 'free');

	// Format the billing amount
	const formattedBillingAmount = $derived(() => {
		const amount = data.subscription.data.attributes.billingAmount;
		if (!amount) return '$0';
		return `$${(amount / 100).toFixed(2)}`;
	});

	// Format dates
	const formatDate = (dateString: string | null) => {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	};

	async function upgradeToPro() {
		if (!canManageBilling) return;

		loading = true;
		error = null;

		try {
			const currentUrl = window.location.href;
			const requestBody: CreateCheckoutSessionRequest = {
				successUrl: currentUrl,
				cancelUrl: currentUrl
			};

			const response = await post(
				`/api/organizations/${data.organization.data.id}/checkout-session`,
				requestBody,
				{
					fetch: window.fetch,
					requestSchema: createCheckoutSessionRequestSchema,
					responseSchema: createCheckoutSessionResponseSchema
				}
			);

			// Redirect to Stripe checkout
			window.location.href = response.data.attributes.url;
		} catch (err) {
			console.error('Failed to create checkout session:', err);
			error = 'Failed to start checkout. Please try again.';
			loading = false;
		}
	}

	async function manageBilling() {
		if (!canManageBilling) return;

		loading = true;
		error = null;

		try {
			const currentUrl = window.location.href;
			const requestBody: CreateBillingPortalSessionRequest = {
				returnUrl: currentUrl
			};

			const response = await post(
				`/api/organizations/${data.organization.data.id}/billing-portal-session`,
				requestBody,
				{
					fetch: window.fetch,
					requestSchema: createBillingPortalSessionRequestSchema,
					responseSchema: createBillingPortalSessionResponseSchema
				}
			);

			// Redirect to Stripe billing portal
			window.location.href = response.data.attributes.url;
		} catch (err) {
			console.error('Failed to create billing portal session:', err);
			error = 'Failed to open billing portal. Please try again.';
			loading = false;
		}
	}
</script>

<SettingsCard>
	<SettingsCardTitle>Billing & Subscription</SettingsCardTitle>

	<div class="space-y-6">
		<!-- Current Plan Display -->
		<div
			class={`rounded-box p-6 ${isProPlan ? 'bg-base-200' : 'from-primary/10 to-secondary/10 border-primary/30 border-2 bg-gradient-to-br'}`}
		>
			<div class="flex items-start justify-between">
				<div class="flex-1">
					<div class="flex items-center gap-2">
						<h3 class="text-2xl font-bold">
							{#if isProPlan}
								<span class="flex items-center gap-2">
									<Sparkles class="text-primary size-6" />
									Pro Plan
								</span>
							{:else}
								Free Plan
							{/if}
						</h3>
					</div>

					{#if isProPlan}
						<p class="text-base-content/70 mt-2">
							You're subscribed to the Pro plan with all premium features.
						</p>
					{:else}
						<p class="text-base-content/70 mt-2">
							You're currently on the Free plan. Upgrade to Pro to unlock advanced features and grow
							your business.
						</p>
					{/if}

					{#if isFreePlan && canManageBilling}
						<div class="mt-4 flex flex-wrap items-center gap-3">
							<button
								class="btn btn-primary btn-lg gap-2 shadow-lg"
								onclick={upgradeToPro}
								disabled={loading}
							>
								{#if loading}
									<span class="loading loading-spinner"></span>
								{:else}
									<Sparkles class="size-5" />
								{/if}
								Upgrade to Pro Now
							</button>
							<span class="text-base-content/60 text-sm font-medium"> Get started in minutes </span>
						</div>
					{/if}
				</div>

				{#if isProPlan}
					<div class="badge badge-primary badge-lg gap-2">
						<CheckCircle2 class="size-4" />
						Active
					</div>
				{/if}
			</div>

			{#if isProPlan && data.subscription.data.attributes.billingPeriodEnd}
				<div class="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div>
						<p class="text-base-content/60 text-sm">Billing Amount</p>
						<p class="text-lg font-semibold">{formattedBillingAmount()}/month</p>
					</div>
					<div>
						<p class="text-base-content/60 text-sm">Next Billing Date</p>
						<p class="text-lg font-semibold">
							{formatDate(data.subscription.data.attributes.billingPeriodEnd)}
						</p>
					</div>
				</div>
			{/if}
		</div>

		<!-- Plan Features Comparison -->
		<div class="grid gap-4 md:grid-cols-2">
			<!-- Free Plan Card -->
			<div class="rounded-box border-base-300 bg-base-100 border p-6">
				<h4 class="text-lg font-semibold">Free Plan</h4>
				<p class="mt-2 text-3xl font-bold">$0<span class="text-base font-normal">/month</span></p>
				<ul class="mt-4 space-y-2">
					<li class="flex items-start gap-2">
						<CheckCircle2 class="text-success mt-0.5 size-5 flex-shrink-0" />
						<span class="text-sm">Basic features</span>
					</li>
					<li class="flex items-start gap-2">
						<CheckCircle2 class="text-success mt-0.5 size-5 flex-shrink-0" />
						<span class="text-sm">Standard support</span>
					</li>
					<li class="flex items-start gap-2">
						<CheckCircle2 class="text-success mt-0.5 size-5 flex-shrink-0" />
						<span class="text-sm">Community access</span>
					</li>
				</ul>
			</div>

			<!-- Pro Plan Card -->
			<div
				class={`rounded-box border-primary border-2 p-6 shadow-lg transition-all ${isFreePlan ? 'from-primary/5 to-secondary/5 bg-gradient-to-br' : 'bg-base-100 ring-primary ring-2'}`}
			>
				<div class="flex items-center justify-between">
					<h4 class="text-lg font-semibold">Pro Plan</h4>
					<div class="flex items-center gap-2">
						{#if isProPlan}
							<span class="badge badge-sm badge-primary">Current</span>
						{:else}
							<span class="badge badge-sm badge-accent">Recommended</span>
						{/if}
						<Sparkles class="text-primary size-5" />
					</div>
				</div>
				<p class="mt-2 text-3xl font-bold">$29<span class="text-base font-normal">/month</span></p>
				<ul class="mt-4 space-y-2">
					<li class="flex items-start gap-2">
						<CheckCircle2 class="text-success mt-0.5 size-5 flex-shrink-0" />
						<span class="text-sm">All Free features</span>
					</li>
					<li class="flex items-start gap-2">
						<CheckCircle2 class="text-success mt-0.5 size-5 flex-shrink-0" />
						<span class="text-sm">Advanced features</span>
					</li>
					<li class="flex items-start gap-2">
						<CheckCircle2 class="text-success mt-0.5 size-5 flex-shrink-0" />
						<span class="text-sm">Priority support</span>
					</li>
					<li class="flex items-start gap-2">
						<CheckCircle2 class="text-success mt-0.5 size-5 flex-shrink-0" />
						<span class="text-sm">Custom integrations</span>
					</li>
				</ul>

				{#if isFreePlan && canManageBilling}
					<button
						class="btn btn-primary btn-block mt-4 gap-2"
						onclick={upgradeToPro}
						disabled={loading}
					>
						{#if loading}
							<span class="loading loading-spinner"></span>
						{:else}
							<Sparkles class="size-4" />
						{/if}
						Get Pro
					</button>
				{/if}
			</div>
		</div>

		{#if error}
			<div role="alert" class="alert alert-error">
				<span>{error}</span>
			</div>
		{/if}
	</div>

	<SettingsCardActions>
		{#if isProPlan}
			<button
				class="btn btn-neutral"
				onclick={manageBilling}
				disabled={!canManageBilling || loading}
			>
				{#if loading}
					<span class="loading loading-spinner"></span>
				{:else}
					<CreditCard class="size-4" />
				{/if}
				Manage Subscription
			</button>
		{/if}
	</SettingsCardActions>
</SettingsCard>
