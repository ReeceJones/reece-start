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
			href: `/app/${params.organizationId}/settings`
		},
		{
			name: 'Members',
			icon: Users,
			href: `/app/${params.organizationId}/settings/members`
		},
		{
			name: 'Billing',
			icon: Receipt,
			href: `/app/${params.organizationId}/settings/billing`
		},
		{
			name: 'Payments',
			icon: DollarSign,
			href: `/app/${params.organizationId}/settings/payments`
		}
	];

	const activeRoute = $derived(routes.find((route) => url === route.href));
</script>

<div class="flex flex-col gap-6 lg:flex-row">
	<div>
		<ul
			class="menu menu-horizontal lg:menu-vertical rounded-box bg-base-200 gap-1 shadow-sm lg:w-56"
		>
			{#each routes as route}
				<li class={clsx(url === route.href && activeClass)}>
					<a href={route.href}>
						<route.icon class="size-5" />
						<span>{route.name}</span>
					</a>
				</li>
			{/each}
		</ul>
	</div>
	<div class="card rounded-box bg-base-200 flex-1 shadow-sm">
		<div class="card-body">
			{#if activeRoute}
				<p class="card-title">{activeRoute?.name}</p>
			{/if}
			{@render children?.()}
		</div>
	</div>
</div>
