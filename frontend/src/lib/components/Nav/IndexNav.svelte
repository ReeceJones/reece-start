<script lang="ts">
	import { getScopes, hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { Bug, Building, House, Shield, Users } from 'lucide-svelte';
	import { page } from '$app/state';
	import clsx from 'clsx/lite';

	const isAdmin = $derived(hasScope(UserScope.Admin));

	const routes = [
		{
			name: 'Home',
			icon: House,
			href: '/app'
		}
	];

	const activeClass = 'bg-base-300';
</script>

<ul class="menu menu-vertical bg-base-200 rounded-box mt-2 w-full space-y-1">
	<li class="menu-title">Application</li>
	{#each routes as route}
		<li>
			<a href={route.href} class={clsx(page.url.pathname === route.href && activeClass)}>
				<route.icon class="h-4 w-4" />
				{route.name}
			</a>
		</li>
	{/each}
	<li>
		<details open>
			<summary>
				<Shield class="h-4 w-4" />
				Admin
			</summary>
			<ul class="space-y-1">
				<li class="mt-1">
					<a href="/app/admin/users">
						<Users class="h-4 w-4" />
						Users</a
					>
				</li>
				<li>
					<a href="/app/admin/debug">
						<Bug class="h-4 w-4" />
						Debug</a
					>
				</li>
			</ul>
		</details>
	</li>
</ul>
