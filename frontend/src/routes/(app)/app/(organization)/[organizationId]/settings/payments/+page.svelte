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
	import { t } from '$lib/i18n';

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
			error = $t('payments.failedToOpenStripeDashboard');
		} finally {
			loading = false;
		}
	}
</script>

<SettingsCard>
	<SettingsCardTitle>{$t('payments.title')}</SettingsCardTitle>

	<div class="space-y-4">
		<p class="text-sm text-base-content/70">
			{$t('payments.description')}
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
					<span class="loading loading-sm loading-spinner"></span>
				{:else}
					<ExternalLink class="h-4 w-4" />
				{/if}
				{$t('payments.openStripeDashboard')}
			</button>
		</SettingsCardActions>
	</div>
</SettingsCard>
