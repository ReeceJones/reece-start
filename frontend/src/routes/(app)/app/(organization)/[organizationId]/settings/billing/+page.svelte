<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { CreditCard, CheckCircle2, Sparkles } from 'lucide-svelte';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import type { PageProps } from './$types';
	import { t } from '$lib/i18n';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import LoadingIcon from '$lib/components/Icons/LoadingIcon.svelte';
	import { enhance } from '$app/forms';
	import FormActionStatus from '$lib/components/Form/FormActionStatus.svelte';
	import { getFormResult, type FormResult, type FormResultCallback } from '$lib/schemas/form';

	const { data }: PageProps = $props();

	let submitting = $state(false);
	let checkoutFormHeaderResult: FormResult | undefined = $state(undefined);
	let checkoutFormCardResult: FormResult | undefined = $state(undefined);
	let portalFormResult: FormResult | undefined = $state(undefined);

	const canManageBilling = $derived(hasScope(UserScope.OrganizationBillingUpdate));
	const isProPlan = $derived(data.subscription.data.attributes.plan === 'pro');
	const isFreePlan = $derived(data.subscription.data.attributes.plan === 'free');

	// Format the billing amount
	const formattedBillingAmount = $derived.by(() => {
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
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{$t('billing.title')}</Card.Title>
	</Card.Header>
	<Card.Content>
		<div class="space-y-6">
			<!-- Current Plan Display -->
			<div class="p-6">
				<div class="flex items-start justify-between">
					<div class="flex-1">
						<div class="flex items-center gap-2">
							<h3 class="text-2xl font-bold">
								{#if isProPlan}
									<span class="flex items-center gap-2">
										<Sparkles class="text-primary size-6" />
										{$t('billing.proPlan')}
									</span>
								{:else}
									{$t('billing.freePlan')}
								{/if}
							</h3>
						</div>

						{#if isProPlan}
							<p class="text-muted-foreground mt-2">
								{$t('billing.proDescription')}
							</p>
						{:else}
							<p class="text-muted-foreground mt-2">
								{$t('billing.freeDescription')}
							</p>
						{/if}

						{#if isFreePlan && canManageBilling}
							<form
								method="post"
								action="?/checkout"
								use:enhance={(): FormResultCallback => {
									submitting = true;

									return ({ update, result }) => {
										update();
										checkoutFormHeaderResult = getFormResult(result);
										submitting = false;
									};
								}}
								enctype="multipart/form-data"
								class="mt-4 flex flex-wrap items-center gap-3"
							>
								<Button
									type="submit"
									variant="default"
									size="lg"
									class="shadow-lg"
									disabled={submitting}
								>
									<LoadingIcon loading={submitting}>
										{#snippet icon()}
											<Sparkles class="size-5" />
										{/snippet}
									</LoadingIcon>
									{$t('billing.upgradeToPro')}
								</Button>
								<span class="text-muted-foreground/80 text-sm font-medium"
									>{$t('billing.getStartedInMinutes')}</span
								>
							</form>
							<div class="mt-4">
								<FormActionStatus
									form={checkoutFormHeaderResult}
									success={$t('billing.successfullyStartedCheckout')}
									failure={$t('billing.failedToStartCheckout')}
								/>
							</div>
						{/if}
					</div>

					{#if isProPlan}
						<Badge variant="default" class="gap-2 px-3 py-1.5 text-sm">
							<CheckCircle2 class="size-4" />
							{$t('billing.active')}
						</Badge>
					{/if}
				</div>

				{#if isProPlan && data.subscription.data.attributes.billingPeriodEnd}
					<div class="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
						<div>
							<p class="text-muted-foreground text-sm">{$t('billing.billingAmount')}</p>
							<p class="text-lg font-semibold">
								{formattedBillingAmount}{$t('billing.perMonth')}
							</p>
						</div>
						<div>
							<p class="text-muted-foreground text-sm">{$t('billing.nextBillingDate')}</p>
							<p class="text-lg font-semibold">
								{formatDate(data.subscription.data.attributes.billingPeriodEnd)}
							</p>
						</div>
					</div>
				{/if}
			</div>

			{#if isFreePlan}
				<!-- Plan Features Comparison -->
				<div class="grid gap-4 md:grid-cols-2">
					<!-- Free Plan Card -->
					<div class="border-border bg-background rounded-lg border p-6">
						<h4 class="text-lg font-semibold">{$t('billing.freePlan')}</h4>
						<p class="mt-2 text-3xl font-bold">
							$0<span class="text-base font-normal">{$t('billing.perMonth')}</span>
						</p>
						<ul class="mt-4 space-y-2">
							<li class="flex items-start gap-2">
								<CheckCircle2
									class="mt-0.5 size-5 flex-shrink-0 text-green-600 dark:text-green-400"
								/>
								<span class="text-sm">{$t('billing.basicFeatures')}</span>
							</li>
							<li class="flex items-start gap-2">
								<CheckCircle2
									class="mt-0.5 size-5 flex-shrink-0 text-green-600 dark:text-green-400"
								/>
								<span class="text-sm">{$t('billing.standardSupport')}</span>
							</li>
							<li class="flex items-start gap-2">
								<CheckCircle2
									class="mt-0.5 size-5 flex-shrink-0 text-green-600 dark:text-green-400"
								/>
								<span class="text-sm">{$t('billing.communityAccess')}</span>
							</li>
						</ul>
					</div>

					<!-- Pro Plan Card -->
					<div class="border-primary rounded-lg border-2 p-6 shadow-lg transition-all">
						<div class="flex items-center justify-between">
							<h4 class="text-lg font-semibold">{$t('billing.proPlan')}</h4>
							<div class="flex items-center gap-2">
								{#if isProPlan}
									<Badge variant="default" class="text-xs">{$t('billing.current')}</Badge>
								{:else}
									<Badge variant="accent" class="text-xs">{$t('billing.recommended')}</Badge>
								{/if}
								<Sparkles class="text-primary size-5" />
							</div>
						</div>
						<p class="mt-2 text-3xl font-bold">
							$29<span class="text-base font-normal">{$t('billing.perMonth')}</span>
						</p>
						<ul class="mt-4 space-y-2">
							<li class="flex items-start gap-2">
								<CheckCircle2
									class="mt-0.5 size-5 flex-shrink-0 text-green-600 dark:text-green-400"
								/>
								<span class="text-sm">{$t('billing.allFreeFeatures')}</span>
							</li>
							<li class="flex items-start gap-2">
								<CheckCircle2
									class="mt-0.5 size-5 flex-shrink-0 text-green-600 dark:text-green-400"
								/>
								<span class="text-sm">{$t('billing.advancedFeatures')}</span>
							</li>
							<li class="flex items-start gap-2">
								<CheckCircle2
									class="mt-0.5 size-5 flex-shrink-0 text-green-600 dark:text-green-400"
								/>
								<span class="text-sm">{$t('billing.prioritySupport')}</span>
							</li>
							<li class="flex items-start gap-2">
								<CheckCircle2
									class="mt-0.5 size-5 flex-shrink-0 text-green-600 dark:text-green-400"
								/>
								<span class="text-sm">{$t('billing.customIntegrations')}</span>
							</li>
						</ul>

						{#if isFreePlan && canManageBilling}
							<form
								method="post"
								action="?/checkout"
								use:enhance={(): FormResultCallback => {
									submitting = true;

									return ({ update, result }) => {
										update();
										checkoutFormCardResult = getFormResult(result);
										submitting = false;
									};
								}}
								enctype="multipart/form-data"
								class="mt-4"
							>
								<Button type="submit" variant="accent" class="w-full" disabled={submitting}>
									<LoadingIcon loading={submitting}>
										{#snippet icon()}
											<Sparkles class="size-4" />
										{/snippet}
									</LoadingIcon>
									{$t('billing.getPro')}
								</Button>
							</form>
							<div class="mt-4">
								<FormActionStatus
									form={checkoutFormCardResult}
									success={$t('billing.successfullyStartedCheckout')}
									failure={$t('billing.failedToStartCheckout')}
								/>
							</div>
						{/if}
					</div>
				</div>
			{/if}
		</div>
	</Card.Content>
	{#if isProPlan}
		<Card.Action>
			<form
				method="post"
				action="?/portal"
				use:enhance={(): FormResultCallback => {
					submitting = true;

					return ({ update, result }) => {
						update();
						portalFormResult = getFormResult(result);
						submitting = false;
					};
				}}
				enctype="multipart/form-data"
			>
				<Button type="submit" variant="secondary" disabled={!canManageBilling || submitting}>
					<LoadingIcon loading={submitting}>
						{#snippet icon()}
							<CreditCard class="size-4" />
						{/snippet}
					</LoadingIcon>
					{$t('billing.manageSubscription')}
				</Button>
			</form>
			<div class="mt-4">
				<FormActionStatus
					form={portalFormResult}
					success={$t('billing.successfullyOpenedBillingPortal')}
					failure={$t('billing.failedToOpenBillingPortal')}
				/>
			</div>
		</Card.Action>
	{/if}
</Card.Root>
