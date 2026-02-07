<script lang="ts">
	import { DollarSign, Receipt, Settings, Users } from 'lucide-svelte';
	import type { LayoutProps } from './$types';
	import { page } from '$app/state';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as m from '$lib/paraglide/messages';

	const { children, params }: LayoutProps = $props();
	const url = $derived(page.url.pathname);
	const baseUrl = $derived(`/app/${params.organizationId}/settings`);
	const routes = $derived([
		{
			name: m.settings__general(),
			icon: Settings,
			href: baseUrl
		},
		{
			name: m.settings__members(),
			icon: Users,
			href: `${baseUrl}/members`
		},
		{
			name: m.settings__billing(),
			icon: Receipt,
			href: `${baseUrl}/billing`
		},
		{
			name: m.settings__payments(),
			icon: DollarSign,
			href: `${baseUrl}/payments`
		}
	]);
	const activeRoute = $derived(routes.find((route) => url === route.href));
</script>

<div class="flex flex-col gap-6">
	<div class="space-y-4">
		<h1 class="text-3xl font-bold">{m.settings__title()}</h1>
		<Tabs.Root value={activeRoute?.href}>
			<Tabs.List>
				{#each routes as route (route.href)}
					<Tabs.Trigger value={route.href}>
						{#snippet child({ props })}
							<a href={route.href} {...props}>
								<route.icon class="size-4" />
								{route.name}
							</a>
						{/snippet}
					</Tabs.Trigger>
				{/each}
			</Tabs.List>
		</Tabs.Root>
	</div>
	{@render children?.()}
</div>
