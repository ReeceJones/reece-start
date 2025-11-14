<script lang="ts">
	import type { Organization } from '$lib/schemas/organization';
	import { ExternalLink } from 'lucide-svelte';
	import { enhance } from '$app/forms';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { t } from '$lib/i18n';
	import { get } from 'svelte/store';

	const { organization }: { organization: Organization } = $props();

	let loading = $state(false);
	let error = $state('');

	const canAccessStripe = $derived(hasScope(UserScope.OrganizationStripeUpdate));
</script>

{#if organization.data.meta.stripe.onboardingStatus === 'missing_requirements' || organization.data.meta.stripe.onboardingStatus === 'missing_capabilities'}
	<div class="alert alert-warning mb-6">
		<div class="flex flex-col gap-2">
			<p class="font-semibold">
				{$t('onboarding.stripeAlert.missingRequirements')}
			</p>
			{#if !canAccessStripe}
				<p class="text-base-content/70 text-sm">
					{$t('onboarding.stripeAlert.adminPermissionsRequired')}
				</p>
			{/if}
			<form
				method="post"
				action={`/app/${organization.data.id}/stripe-onboarding`}
				use:enhance={() => {
					loading = true;
					error = '';
					return ({ result, update }) => {
						update();
						const translate = get(t);
						if (result.type === 'failure') {
							loading = false;
							error =
								(result.data?.message as string) ??
								translate('onboarding.stripeAlert.somethingWentWrong');
						} else if (result.type === 'error') {
							loading = false;
							error = result.error ?? translate('onboarding.stripeAlert.somethingWentWrong');
						}
					};
				}}
			>
				<button class="btn btn-sm" disabled={loading || !canAccessStripe}>
					{#if loading}
						<span class="loading loading-xs loading-spinner"></span>
					{:else}
						<ExternalLink class="size-4" />
					{/if}
					<span class="ml-1">{$t('onboarding.stripeAlert.openStripe')}</span>
				</button>
			</form>
			{#if error}
				<div class="alert alert-error mt-3">
					<span>{error}</span>
				</div>
			{/if}
		</div>
	</div>
{/if}

{#if organization.data.meta.stripe.onboardingStatus === 'pending'}
	<div class="alert alert-info mb-6">
		<div class="flex flex-col gap-2">
			<p class="font-semibold">
				{$t('onboarding.stripeAlert.settingUp')}
			</p>
		</div>
	</div>
{/if}
