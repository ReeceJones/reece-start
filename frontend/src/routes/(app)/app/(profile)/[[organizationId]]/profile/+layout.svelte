<script lang="ts">
	const { children, params } = $props();
	import { page } from '$app/state';
	import clsx from 'clsx/lite';
	import { User, Lock } from 'lucide-svelte';

	const url = $derived(page.url.pathname);
	const activeClass = 'bg-base-300 rounded-md';
	const baseUrl = $derived(
		params.organizationId ? `/app/${params.organizationId}/profile` : '/app/profile'
	);
	const routes = $derived([
		{
			name: 'Profile',
			icon: User,
			href: baseUrl
		},
		{
			name: 'Security',
			icon: Lock,
			href: `${baseUrl}/security`
		}
	]);
</script>

<div class="flex flex-col gap-6">
	<div class="space-y-4">
		<h1 class="text-3xl font-bold">Settings</h1>
		<ul class="menu menu-horizontal rounded-box bg-base-200 gap-1 shadow-sm">
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
	{@render children?.()}
</div>
