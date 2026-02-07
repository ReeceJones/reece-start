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
	import * as m from '$lib/paraglide/messages';
	import { Button } from '$lib/components/ui/button';
	import LoadingIcon from '$lib/components/Icons/LoadingIcon.svelte';

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
			error = m.payments__failed_to_open_stripe_dashboard();
		} finally {
			loading = false;
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{m.payments__title()}</Card.Title>
	</Card.Header>
	<Card.Content>
		<div class="space-y-4">
			<p class="text-muted-foreground text-sm">
				{m.payments__description()}
			</p>

			{#if error}
				<Alert.Root variant="destructive">
					<span>{error}</span>
				</Alert.Root>
			{/if}
		</div>
	</Card.Content>
	<Card.Action class="px-6">
		<Button variant="default" onclick={openStripeDashboard} disabled={loading || !canAccessStripe}>
			<LoadingIcon {loading}>
				{#snippet icon()}
					<ExternalLink class="h-4 w-4" />
				{/snippet}
			</LoadingIcon>
			{m.payments__open_stripe_dashboard()}
		</Button>
	</Card.Action>
</Card.Root>
