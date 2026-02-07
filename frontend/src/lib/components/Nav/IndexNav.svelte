<script lang="ts">
	import { hasScope } from '$lib/auth';
	import { UserScope } from '$lib/schemas/jwt';
	import { Bug, House, Users } from 'lucide-svelte';
	import { page } from '$app/state';
	import * as m from '$lib/paraglide/messages';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';

	const isAdmin = $derived(hasScope(UserScope.Admin));
	const currentPath = $derived(page.url.pathname);
</script>

<Sidebar.Menu class="w-full space-y-1">
	<Sidebar.GroupLabel>{m.application()}</Sidebar.GroupLabel>
	<Sidebar.MenuItem>
		<Sidebar.MenuButton isActive={currentPath === '/app'}>
			{#snippet child({ props })}
				<a href="/app" {...props}>
					<House class="size-4" />
					<span>{m.home()}</span>
				</a>
			{/snippet}
		</Sidebar.MenuButton>
	</Sidebar.MenuItem>
	{#if isAdmin}
		<Sidebar.GroupLabel>{m.admin()}</Sidebar.GroupLabel>
		<Sidebar.MenuItem>
			<Sidebar.MenuButton isActive={currentPath === '/app/admin/users'}>
				{#snippet child({ props })}
					<a href="/app/admin/users" {...props}>
						<Users class="size-4" />
						<span>{m.users()}</span>
					</a>
				{/snippet}
			</Sidebar.MenuButton>
		</Sidebar.MenuItem>
		<Sidebar.MenuItem>
			<Sidebar.MenuButton isActive={currentPath === '/app/admin/debug'}>
				{#snippet child({ props })}
					<a href="/app/admin/debug" {...props}>
						<Bug class="size-4" />
						<span>{m.debug()}</span>
					</a>
				{/snippet}
			</Sidebar.MenuButton>
		</Sidebar.MenuItem>
	{/if}
</Sidebar.Menu>
