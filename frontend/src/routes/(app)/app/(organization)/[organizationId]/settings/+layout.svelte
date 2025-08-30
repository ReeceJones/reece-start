<script lang="ts">
	import { DollarSign, Receipt, Settings, Users } from 'lucide-svelte';
	import type { LayoutProps } from './$types';
	import { page } from '$app/state';
	import clsx from 'clsx/lite';
	// TODO: refactor this into common component since it is also used in user settings

	const { children, params }: LayoutProps = $props();
	const url = $derived(page.url.pathname);

	const activeClass = 'bg-base-300 rounded-md';
	const routes = [
		{
			name: 'General',
			icon: Settings,
			href: `/app/${params.organizationId}/settings`,
			exact: true
		},
		{
			name: 'Members',
			icon: Users,
			href: `/app/${params.organizationId}/settings/members`,
			exact: false
		},
		{
			name: 'Billing',
			icon: Receipt,
			href: `/app/${params.organizationId}/settings/billing`,
			exact: false
		},
		{
			name: 'Payments',
			icon: DollarSign,
			href: `/app/${params.organizationId}/settings/payments`,
			exact: false
		}
	];

	const activeRoute = $derived(routes.find((route) => url === route.href));
</script>

<div class="flex flex-col gap-8">
	<div class="space-y-4">
		<h1 class="text-3xl font-bold">Settings</h1>
		<ul
			class="menu menu-horizontal rounded-box bg-base-200 max-w-full flex-nowrap gap-1 overflow-auto shadow-sm"
		>
			{#each routes as route}
				<li
					class={clsx(
						(url === route.href || (!route.exact && url.startsWith(route.href))) && activeClass
					)}
				>
					<a href={route.href}>
						<route.icon class="size-5" />
						<span>{route.name}</span>
					</a>
				</li>
			{/each}
		</ul>
	</div>
	{@render children?.()}
</div>
