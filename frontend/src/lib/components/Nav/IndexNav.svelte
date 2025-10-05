<script lang="ts">
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { Bug, House, Shield, Users } from 'lucide-svelte';
	import { page } from '$app/state';
	import clsx from 'clsx/lite';

	const isAdmin = $derived(hasScope(UserScope.Admin));
	const activeClass = 'bg-base-300';
</script>

<ul class="menu menu-vertical mt-2 w-full space-y-1 rounded-box bg-base-200">
	<li class="menu-title">Application</li>
	<li>
		<a href="/app" class={clsx(page.url.pathname === '/app' && activeClass)}>
			<House class="h-4 w-4" />
			Home
		</a>
	</li>
	{#if isAdmin}
		<li>
			<details open>
				<summary>
					<Shield class="h-4 w-4" />
					Admin
				</summary>
				<ul class="space-y-1">
					<li class="mt-1">
						<a
							href="/app/admin/users"
							class={clsx(page.url.pathname === '/app/admin/users' && activeClass)}
						>
							<Users class="h-4 w-4" />
							Users</a
						>
					</li>
					<li>
						<a
							href="/app/admin/debug"
							class={clsx(page.url.pathname === '/app/admin/debug' && activeClass)}
						>
							<Bug class="h-4 w-4" />
							Debug</a
						>
					</li>
				</ul>
			</details>
		</li>
	{/if}
</ul>
