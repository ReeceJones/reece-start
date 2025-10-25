<script lang="ts">
	import { DollarSign, Receipt, Settings, Users } from 'lucide-svelte';
	import type { LayoutProps } from './$types';
	import { page } from '$app/state';
	import clsx from 'clsx/lite';
	import { t } from '$lib/i18n';
	// TODO: refactor this into common component since it is also used in user settings

	const { children, params }: LayoutProps = $props();
	const url = $derived(page.url.pathname);

	const activeClass = 'bg-base-300 rounded-md';
	const routes = [
		{
			name: $t('settings.general'),
			icon: Settings,
			href: `/app/${params.organizationId}/settings`,
			exact: true
		},
		{
			name: $t('settings.members'),
			icon: Users,
			href: `/app/${params.organizationId}/settings/members`,
			exact: false
		},
		{
			name: $t('settings.billing'),
			icon: Receipt,
			href: `/app/${params.organizationId}/settings/billing`,
			exact: false
		},
		{
			name: $t('settings.payments'),
			icon: DollarSign,
			href: `/app/${params.organizationId}/settings/payments`,
			exact: false
		}
	];
</script>

<div class="flex flex-col gap-8">
	<div class="space-y-4">
		<h1 class="text-3xl font-bold">{$t('settings.title')}</h1>
		<ul
			class="menu menu-horizontal rounded-box bg-base-200 max-w-full flex-nowrap gap-1 overflow-auto shadow-sm"
		>
			{#each routes as route (route.href)}
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
