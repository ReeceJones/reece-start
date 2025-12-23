<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import { post } from '$lib';
	import {
		createStripeDashboardLinkResponseSchema,
		type Organization
	} from '$lib/schemas/organization';
	import { ExternalLink } from 'lucide-svelte';
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { z } from 'zod';
	import { t } from '$lib/i18n';

	const { data }: { data: { organization: Organization } } = $props();

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

<Card.Root>
	<Card.Header>
		<Card.Title>{$t('payments.title')}</Card.Title>
	</Card.Header>
	<Card.Content>
		<div class="space-y-4">
			<p class="text-base-content/70 text-sm">
				{$t('payments.description')}
			</p>

			{#if error}
				<Alert.Root variant="destructive">
					<span>{error}</span>
				</Alert.Root>
			{/if}
		</div>
	</Card.Content>
	<Card.Action>
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
	</Card.Action>
</Card.Root>
