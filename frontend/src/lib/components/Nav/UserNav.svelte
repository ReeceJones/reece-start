<script lang="ts">
	import { ChevronsUpDown, EyeOff, LogOut, Settings, User } from 'lucide-svelte';
	import { getSelfUserResponseSchema } from '$lib/schemas/user';
	import type { z } from 'zod';
	import { page } from '$app/state';
	import { getIsImpersonatingUser } from '$lib/auth';
	import * as m from '$lib/paraglide/messages';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

	const { user }: { user: z.infer<typeof getSelfUserResponseSchema> } = $props();
	const organizationId = $derived(page.params.organizationId);
	const profileHref = $derived(organizationId ? `/app/${organizationId}/profile` : '/app/profile');
	const isImpersonatingUser = $derived(getIsImpersonatingUser());

	let signoutForm: HTMLFormElement;
	let stopImpersonationForm: HTMLFormElement;

	function handleSignOut() {
		signoutForm?.requestSubmit();
	}

	function handleStopImpersonation() {
		stopImpersonationForm?.requestSubmit();
	}
</script>

<Sidebar.Menu class="w-full">
	<Sidebar.MenuItem class="w-full">
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						{...props}
						class="w-full py-6 data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
					>
						{#if user.data.meta.logoDistributionUrl}
							<img
								src={user.data.meta.logoDistributionUrl}
								alt="User logo"
								class="size-8 rounded-md"
							/>
						{:else}
							<div
								class="flex size-8 shrink-0 items-center justify-center rounded-lg bg-sidebar-accent"
							>
								<User class="size-6" />
							</div>
						{/if}
						<div class="flex flex-col items-start">
							<span>{user.data.attributes.name || m.nav__profile()}</span>
							<span class="text-xs text-sidebar-foreground/70">{user.data.attributes.email}</span>
						</div>
						<ChevronsUpDown class="ms-auto" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content side="top" class="w-(--bits-dropdown-menu-anchor-width)">
				<DropdownMenu.Item>
					{#snippet child({ props })}
						<a href={profileHref} {...props}>
							<Settings class="size-4" />
							{m.nav__settings()}
						</a>
					{/snippet}
				</DropdownMenu.Item>
				{#if isImpersonatingUser}
					<DropdownMenu.Item onSelect={handleStopImpersonation} variant="destructive">
						<EyeOff class="size-4" />
						{m.nav__stop_impersonation()}
					</DropdownMenu.Item>
				{/if}
				<DropdownMenu.Item onSelect={handleSignOut} variant="destructive">
					<LogOut class="size-4" />
					{m.nav__logout()}
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
		<form
			bind:this={signoutForm}
			action="/app?/signout"
			method="POST"
			enctype="multipart/form-data"
		></form>
		<form
			bind:this={stopImpersonationForm}
			action="/app?/stopImpersonation"
			method="POST"
			enctype="multipart/form-data"
		></form>
	</Sidebar.MenuItem>
</Sidebar.Menu>
