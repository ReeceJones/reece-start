<script lang="ts">
	const { children } = $props();
	import { page } from '$app/state';
	import clsx from 'clsx/lite';
	import { User, Lock } from 'lucide-svelte';

	const url = $derived(page.url.pathname);
	const activeClass = 'bg-base-300 rounded-md';
	const routes = [
		{
			name: 'Profile',
			icon: User,
			href: '/app/profile'
		},
		{
			name: 'Security',
			icon: Lock,
			href: '/app/profile/security'
		}
	];

	const activeRoute = $derived(routes.find((route) => url === route.href));
</script>

<div class="flex gap-6">
	<div>
		<ul class="menu rounded-box bg-base-200 w-56 gap-1 shadow-sm">
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
