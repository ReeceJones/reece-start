<script lang="ts">
	import type { Organization } from '$lib/schemas/organization';
	import { House, Folder, DollarSign, Building2, Settings, ArrowLeftRight } from 'lucide-svelte';
	import { t } from '$lib/i18n';

	const { organization }: { organization: Organization } = $props();
</script>

<ul class="menu menu-vertical w-full space-y-1">
	<li class="w-full">
		<div class="dropdown dropdown-start dropdown-bottom w-full p-0">
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
				{organization.data.attributes.name || $t('nav.organization')}
			</div>
			<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
			<ul
				tabindex="0"
				class="dropdown-content menu z-1 ml-0 w-52 -translate-y-0.5 rounded-box bg-base-100 p-2 shadow-sm"
			>
				<li>
					<a href="/app/{organization.data.id}/settings">
						<Settings class="size-4" />
						{$t('nav.settings')}
					</a>
				</li>
				<li>
					<a href="/app">
						<ArrowLeftRight class="size-4" />
						{$t('nav.switchOrganization')}
					</a>
				</li>
			</ul>
		</div>
	</li>
	<li class="menu-title">{$t('nav.application')}</li>
	<li>
		<a href="/app/{organization.data.id}">
			<House class="size-4" />
			{$t('nav.dashboard')}
		</a>
	</li>
	<li>
		<a href="/app/{organization.data.id}/foo">
			<Folder class="size-4" />
			{$t('nav.foo')}
		</a>
	</li>
	<li>
		<a href="/app/{organization.data.id}/bar">
			<DollarSign class="size-4" />
			{$t('nav.bar')}
		</a>
	</li>
	<li>
		<a href="/app/{organization.data.id}/settings">
			<Settings class="size-4" />
			{$t('nav.settings')}
		</a>
	</li>
</ul>
