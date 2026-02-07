<script lang="ts">
	import type { Organization } from '$lib/schemas/organization';
	import { ExternalLink } from 'lucide-svelte';
	import * as Alert from '$lib/components/ui/alert';
	import { enhance } from '$app/forms';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import * as m from '$lib/paraglide/messages';
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
				{m.onboarding__stripe_alert__missing_requirements()}
			</p>
			{#if !canAccessStripe}
				<p>
					{m.onboarding__stripe_alert__admin_permissions_required()}
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
					if (result.type === 'failure') {
						loading = false;
						error =
							(result.data?.message as string) ??
							m.onboarding__stripe_alert__something_went_wrong();
					} else if (result.type === 'error') {
						loading = false;
						error = result.error ?? m.onboarding__stripe_alert__something_went_wrong();
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
				<span class="ml-1">{m.onboarding__stripe_alert__open_stripe()}</span>
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
				{m.onboarding__stripe_alert__setting_up()}
			</p>
		</Alert.Description>
	</Alert.Root>
{/if}
