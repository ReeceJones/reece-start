<script lang="ts">
	import { getScopes, hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { House, Shield } from 'lucide-svelte';
	import { page } from '$app/state';
	import clsx from 'clsx/lite';

	const isAdmin = $derived(hasScope(UserScope.Admin));

	const routes = $derived(
		[
			{
				name: 'Home',
				icon: House,
				href: '/app',
				shouldShow: true
			},
			{
				name: 'Admin',
				icon: Shield,
				href: '/app/admin',
				shouldShow: isAdmin
			}
		].filter((route) => route.shouldShow)
	);

	const activeClass = 'bg-base-300';
</script>

<ul class="menu menu-vertical lg:menu-horizontal bg-base-200 rounded-box mt-2 w-full space-y-1">
	<li class="menu-title">Application</li>
	{#each routes as route}
		<li class="w-full">
			<a href={route.href} class={clsx(page.url.pathname === route.href && activeClass)}>
				<route.icon class="h-4 w-4" />
				{route.name}
			</a>
		</li>
	{/each}
</ul>
