<script lang="ts">
	const { children, params } = $props();
	import { page } from '$app/state';
	import { User, Lock } from 'lucide-svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import { t } from '$lib/i18n';

	const url = $derived(page.url.pathname);
	const baseUrl = $derived(
		params.organizationId ? `/app/${params.organizationId}/profile` : '/app/profile'
	);
	const routes = $derived([
		{
			name: $t('settings.profile'),
			icon: User,
			href: baseUrl
		},
		{
			name: $t('settings.security'),
			icon: Lock,
			href: `${baseUrl}/security`
		}
	]);
	const activeRoute = $derived(routes.find((route) => route.href === url));
</script>

<div class="flex flex-col gap-6">
	<div class="space-y-4">
		<h1 class="text-3xl font-bold">{$t('settings.title')}</h1>
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
