<script lang="ts">
	import UserNav from '$lib/components/Nav/UserNav.svelte';
	import IndexNav from '$lib/components/Nav/IndexNav.svelte';
	import OrganizationNav from '$lib/components/Nav/OrganizationNav.svelte';
	import type { LayoutProps } from './$types';

	const { children, data }: LayoutProps = $props();

	const { user, organization } = data;
</script>

<div class="flex h-screen max-w-screen flex-row gap-4">
	<div class="flex h-full w-56 flex-col justify-between border-r border-base-300 bg-base-200">
		<div>
			{#if organization}
				<OrganizationNav {organization} />
			{:else}
				<IndexNav />
			{/if}
		</div>
		<div>
			<UserNav {user} />
		</div>
	</div>
	<main class="container mx-auto mt-4 mr-4 inline-block flex-1">
		{@render children?.()}
	</main>
</div>
