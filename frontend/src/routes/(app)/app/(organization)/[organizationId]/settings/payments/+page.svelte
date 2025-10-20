<script lang="ts">
	import SettingsCard from '$lib/components/Settings/SettingsCard.svelte';
	import SettingsCardTitle from '$lib/components/Settings/SettingsCardTitle.svelte';
	import SettingsCardActions from '$lib/components/Settings/SettingsCardActions.svelte';
	import { post } from '$lib';
	import { createStripeDashboardLinkResponseSchema } from '$lib/schemas/organization';
	import { ExternalLink } from 'lucide-svelte';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { z } from 'zod';

	const { data }: { data: { organization: any } } = $props();

	let loading = $state(false);
	let error = $state<string | null>(null);

	const canAccessStripe = $derived(hasScope(UserScope.OrganizationStripeUpdate));

	async function openStripeDashboard() {
		if (!data.organization?.data?.id) return;

		loading = true;
		error = null;

		try {
			const response = await post(
				`/api/organizations/${data.organization.data.id}/stripe-dashboard-link`,
				{},
				{
					fetch: window.fetch,
					requestSchema: z.object({}),
					responseSchema: createStripeDashboardLinkResponseSchema
				}
			);

			// Open the dashboard URL in a new tab
			window.open(response.data.attributes.url, '_blank');
		} catch (err) {
			console.error('Failed to create Stripe dashboard link:', err);
			error = 'Failed to open Stripe dashboard. Please try again.';
		} finally {
			loading = false;
		}
	}
</script>

<SettingsCard>
	<SettingsCardTitle>Payments</SettingsCardTitle>

	<div class="space-y-4">
		<p class="text-base-content/70 text-sm">
			Manage your payment settings and view transaction history in your Stripe dashboard.
		</p>

		{#if error}
			<div class="alert alert-error">
				<span>{error}</span>
			</div>
		{/if}

		<SettingsCardActions>
			<button
				class="btn btn-primary"
				onclick={openStripeDashboard}
				disabled={loading || !canAccessStripe}
			>
				{#if loading}
					<span class="loading loading-spinner loading-sm"></span>
				{:else}
					<ExternalLink class="h-4 w-4" />
				{/if}
				Open Stripe Dashboard
			</button>
		</SettingsCardActions>
	</div>
</SettingsCard>
