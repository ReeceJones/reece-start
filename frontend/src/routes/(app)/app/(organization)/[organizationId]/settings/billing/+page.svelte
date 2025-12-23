<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
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
	import { t } from '$lib/i18n';

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
			error = $t('billing.failedToStartCheckout');
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
			error = $t('billing.failedToOpenBillingPortal');
			loading = false;
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{$t('billing.title')}</Card.Title>
	</Card.Header>
	<Card.Content>
		<div class="space-y-6">
			<!-- Current Plan Display -->
			<div
				class={`rounded-box p-6 ${isProPlan ? 'bg-base-200' : 'border-2 border-primary/30 bg-gradient-to-br from-primary/10 to-secondary/10'}`}
			>
				<div class="flex items-start justify-between">
					<div class="flex-1">
						<div class="flex items-center gap-2">
							<h3 class="text-2xl font-bold">
								{#if isProPlan}
									<span class="flex items-center gap-2">
										<Sparkles class="size-6 text-primary" />
										{$t('billing.proPlan')}
									</span>
								{:else}
									{$t('billing.freePlan')}
								{/if}
							</h3>
						</div>

						{#if isProPlan}
							<p class="text-base-content/70 mt-2">
								{$t('billing.proDescription')}
							</p>
						{:else}
							<p class="text-base-content/70 mt-2">
								{$t('billing.freeDescription')}
							</p>
						{/if}

						{#if isFreePlan && canManageBilling}
							<div class="mt-4 flex flex-wrap items-center gap-3">
								<button
									class="btn btn-lg btn-primary gap-2 shadow-lg"
									onclick={upgradeToPro}
									disabled={loading}
								>
									{#if loading}
										<span class="loading loading-spinner"></span>
									{:else}
										<Sparkles class="size-5" />
									{/if}
									{$t('billing.upgradeToPro')}
								</button>
								<span class="text-base-content/60 text-sm font-medium"
									>{$t('billing.getStartedInMinutes')}</span
								>
							</div>
						{/if}
					</div>

					{#if isProPlan}
						<div class="badge badge-lg badge-primary gap-2">
							<CheckCircle2 class="size-4" />
							{$t('billing.active')}
						</div>
					{/if}
				</div>

				{#if isProPlan && data.subscription.data.attributes.billingPeriodEnd}
					<div class="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
						<div>
							<p class="text-base-content/60 text-sm">{$t('billing.billingAmount')}</p>
							<p class="text-lg font-semibold">
								{formattedBillingAmount()}{$t('billing.perMonth')}
							</p>
						</div>
						<div>
							<p class="text-base-content/60 text-sm">{$t('billing.nextBillingDate')}</p>
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
					<h4 class="text-lg font-semibold">{$t('billing.freePlan')}</h4>
					<p class="mt-2 text-3xl font-bold">
						$0<span class="text-base font-normal">{$t('billing.perMonth')}</span>
					</p>
					<ul class="mt-4 space-y-2">
						<li class="flex items-start gap-2">
							<CheckCircle2 class="mt-0.5 size-5 flex-shrink-0 text-success" />
							<span class="text-sm">{$t('billing.basicFeatures')}</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="mt-0.5 size-5 flex-shrink-0 text-success" />
							<span class="text-sm">{$t('billing.standardSupport')}</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="mt-0.5 size-5 flex-shrink-0 text-success" />
							<span class="text-sm">{$t('billing.communityAccess')}</span>
						</li>
					</ul>
				</div>

				<!-- Pro Plan Card -->
				<div
					class={`rounded-box border-2 border-primary p-6 shadow-lg transition-all ${isFreePlan ? 'bg-gradient-to-br from-primary/5 to-secondary/5' : 'bg-base-100 ring-2 ring-primary'}`}
				>
					<div class="flex items-center justify-between">
						<h4 class="text-lg font-semibold">{$t('billing.proPlan')}</h4>
						<div class="flex items-center gap-2">
							{#if isProPlan}
								<span class="badge badge-sm badge-primary">{$t('billing.current')}</span>
							{:else}
								<span class="badge badge-sm badge-accent">{$t('billing.recommended')}</span>
							{/if}
							<Sparkles class="size-5 text-primary" />
						</div>
					</div>
					<p class="mt-2 text-3xl font-bold">
						$29<span class="text-base font-normal">{$t('billing.perMonth')}</span>
					</p>
					<ul class="mt-4 space-y-2">
						<li class="flex items-start gap-2">
							<CheckCircle2 class="mt-0.5 size-5 flex-shrink-0 text-success" />
							<span class="text-sm">{$t('billing.allFreeFeatures')}</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="mt-0.5 size-5 flex-shrink-0 text-success" />
							<span class="text-sm">{$t('billing.advancedFeatures')}</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="mt-0.5 size-5 flex-shrink-0 text-success" />
							<span class="text-sm">{$t('billing.prioritySupport')}</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="mt-0.5 size-5 flex-shrink-0 text-success" />
							<span class="text-sm">{$t('billing.customIntegrations')}</span>
						</li>
					</ul>

					{#if isFreePlan && canManageBilling}
						<button
							class="btn btn-block btn-primary mt-4 gap-2"
							onclick={upgradeToPro}
							disabled={loading}
						>
							{#if loading}
								<span class="loading loading-spinner"></span>
							{:else}
								<Sparkles class="size-4" />
							{/if}
							{$t('billing.getPro')}
						</button>
					{/if}
				</div>
			</div>

			{#if error}
				<Alert.Root variant="destructive">
					<span>{error}</span>
				</Alert.Root>
			{/if}
		</div>
	</Card.Content>
	{#if isProPlan}
		<Card.Action>
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
				{$t('billing.manageSubscription')}
			</button>
		</Card.Action>
	{/if}
</Card.Root>
