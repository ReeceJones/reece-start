<script lang="ts">
	import type { Organization } from '$lib/schemas/organization';
	import { House, Folder, DollarSign, Building2, Settings, RefreshCcw } from 'lucide-svelte';

	const { organization }: { organization: Organization } = $props();
</script>

<ul class="menu menu-vertical w-full">
	<li class="w-full">
		<div class="dropdown dropdown-bottom dropdown-start w-full p-0">
			<div tabindex="0" role="button" class="focus flex w-full gap-2 px-3 py-1.5">
				{#if organization.data.meta.logoDistributionUrl}
					<img
						src={organization.data.meta.logoDistributionUrl}
						alt="Organization logo"
						class="size-6 rounded-sm"
					/>
				{:else}
					<Building2 class="size-5" />
				{/if}
				{organization.data.attributes.name ?? 'Organization'}
			</div>
			<ul
				tabindex="0"
				class="dropdown-content menu bg-base-100 rounded-box z-1 -ml-3 w-52 -translate-y-1.5 p-2 shadow-sm"
			>
				<li>
					<a href="/app/{organization.data.id}/settings">
						<Settings class="size-4" />
						Settings
					</a>
				</li>
				<li>
					<a href="/app">
						<RefreshCcw class="size-4" />
						Switch Organization
					</a>
				</li>
			</ul>
		</div>
	</li>
	<li class="menu-title">Application</li>
	<li>
		<a href="/app/{organization.data.id}">
			<House class="size-4" />
			Dashboard
		</a>
	</li>
	<li>
		<a href="/app/{organization.data.id}/foo">
			<Folder class="size-4" />
			Foo
		</a>
	</li>
	<li>
		<a href="/app/{organization.data.id}/bar">
			<DollarSign class="size-4" />
			Bar
		</a>
	</li>
</ul>
