<script lang="ts">
	import UserNav from '$lib/components/Nav/UserNav.svelte';
	import IndexNav from '$lib/components/Nav/IndexNav.svelte';
	import OrganizationNav from '$lib/components/Nav/OrganizationNav.svelte';
	import type { LayoutProps } from './$types';

	const { children, data }: LayoutProps = $props();

	const { user, organization } = data;
</script>

<div class="max-w-screen flex h-screen flex-row gap-4">
	<div class="border-base-300 bg-base-200 flex h-full w-56 flex-col justify-between border-r">
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
	<main class="container mx-auto mr-4 mt-4 inline-block flex-1">
		{@render children?.()}
	</main>
</div>
