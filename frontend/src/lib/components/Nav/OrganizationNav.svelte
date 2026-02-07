<script lang="ts">
	import type { Organization } from '$lib/schemas/organization';
	import {
		House,
		Folder,
		DollarSign,
		Building2,
		Settings,
		ArrowLeftRight,
		ChevronsUpDown
	} from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
	import { page } from '$app/state';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

	const { organization }: { organization: Organization } = $props();
	const currentPath = $derived(page.url.pathname);
</script>

<Sidebar.Menu class="w-full space-y-1">
	<Sidebar.MenuItem class="w-full">
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						{...props}
						class="w-full py-6 data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
					>
						{#if organization.data.meta.logoDistributionUrl}
							<img
								src={organization.data.meta.logoDistributionUrl}
								alt="Organization logo"
								class="size-8 rounded-md"
							/>
						{:else}
							<div
								class="flex size-8 shrink-0 items-center justify-center rounded-lg bg-sidebar-accent"
							>
								<Building2 class="size-6" />
							</div>
						{/if}
						<span>{organization.data.attributes.name || m.nav__organization()}</span>
						<ChevronsUpDown class="ms-auto" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				side="bottom"
				align="start"
				class="w-(--bits-dropdown-menu-anchor-width)"
			>
				<DropdownMenu.Item>
					{#snippet child({ props })}
						<a href="/app/{organization.data.id}/settings" {...props}>
							<Settings class="size-4" />
							{m.nav__settings()}
						</a>
					{/snippet}
				</DropdownMenu.Item>
				<DropdownMenu.Item>
					{#snippet child({ props })}
						<a href="/app" {...props}>
							<ArrowLeftRight class="size-4" />
							{m.nav__switch_organization()}
						</a>
					{/snippet}
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
	<Sidebar.GroupLabel>{m.nav__application()}</Sidebar.GroupLabel>
	<Sidebar.MenuItem>
		<Sidebar.MenuButton isActive={currentPath === `/app/${organization.data.id}`}>
			{#snippet child({ props })}
				<a href="/app/{organization.data.id}" {...props}>
					<House class="size-4" />
					<span>{m.nav__dashboard()}</span>
				</a>
			{/snippet}
		</Sidebar.MenuButton>
	</Sidebar.MenuItem>
	<Sidebar.MenuItem>
		<Sidebar.MenuButton isActive={currentPath === `/app/${organization.data.id}/foo`}>
			{#snippet child({ props })}
				<a href="/app/{organization.data.id}/foo" {...props}>
					<Folder class="size-4" />
					<span>{m.nav__foo()}</span>
				</a>
			{/snippet}
		</Sidebar.MenuButton>
	</Sidebar.MenuItem>
	<Sidebar.MenuItem>
		<Sidebar.MenuButton isActive={currentPath === `/app/${organization.data.id}/bar`}>
			{#snippet child({ props })}
				<a href="/app/{organization.data.id}/bar" {...props}>
					<DollarSign class="size-4" />
					<span>{m.nav__bar()}</span>
				</a>
			{/snippet}
		</Sidebar.MenuButton>
	</Sidebar.MenuItem>
	<Sidebar.MenuItem>
		<Sidebar.MenuButton isActive={currentPath === `/app/${organization.data.id}/settings`}>
			{#snippet child({ props })}
				<a href="/app/{organization.data.id}/settings" {...props}>
					<Settings class="size-4" />
					<span>{m.nav__settings()}</span>
				</a>
			{/snippet}
		</Sidebar.MenuButton>
	</Sidebar.MenuItem>
</Sidebar.Menu>
