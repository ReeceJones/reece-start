<script lang="ts">
	import type { Organization } from '$lib/schemas/organization';
	import { ExternalLink } from 'lucide-svelte';
	import * as Alert from '$lib/components/ui/alert';
	import { enhance } from '$app/forms';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { t } from '$lib/i18n';
	import { get } from 'svelte/store';
	import { Button } from '../ui/button';
	import { Spinner } from '../ui/spinner';

	const { organization }: { organization: Organization } = $props();

	let loading = $state(false);
	let error = $state('');

	const canAccessStripe = $derived(hasScope(UserScope.OrganizationStripeUpdate));
</script>

{#if organization.data.meta.stripe.onboardingStatus === 'missing_requirements' || organization.data.meta.stripe.onboardingStatus === 'missing_capabilities'}
	<Alert.Root variant="warning" class="mb-6 space-y-1">
		<Alert.Title>Complete stripe setup</Alert.Title>
		<Alert.Description>
			<p>
				{$t('onboarding.stripeAlert.missingRequirements')}
			</p>
			{#if !canAccessStripe}
				<p>
					{$t('onboarding.stripeAlert.adminPermissionsRequired')}
				</p>
			{/if}
		</Alert.Description>
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
			<Button size="sm" type="submit" disabled={loading || !canAccessStripe}>
				{#if loading}
					<Spinner />
				{:else}
					<ExternalLink class="size-4" />
				{/if}
				<span class="ml-1">{$t('onboarding.stripeAlert.openStripe')}</span>
			</Button>
		</form>
		{#if error}
			<Alert.Root>
				<Alert.Title>Error</Alert.Title>
				<Alert.Description>
					{error}
				</Alert.Description>
			</Alert.Root>
		{/if}
	</Alert.Root>
{/if}

{#if organization.data.meta.stripe.onboardingStatus === 'pending'}
	<Alert.Root class="mb-6">
		<Alert.Description>
			<p class="font-semibold">
				{$t('onboarding.stripeAlert.settingUp')}
			</p>
		</Alert.Description>
	</Alert.Root>
{/if}
