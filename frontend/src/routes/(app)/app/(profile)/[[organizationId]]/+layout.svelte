<script lang="ts">
	import UserNav from '$lib/components/Nav/UserNav.svelte';
	import IndexNav from '$lib/components/Nav/IndexNav.svelte';
	import OrganizationNav from '$lib/components/Nav/OrganizationNav.svelte';
	import type { LayoutProps } from './$types';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';

	const { children, data }: LayoutProps = $props();

	const { user, organization } = $derived(data);
</script>

<Sidebar.Provider>
	<Sidebar.Root>
		<Sidebar.Content>
			<Sidebar.Group>
				{#if organization}
					<OrganizationNav {organization} />
				{:else}
					<IndexNav />
				{/if}
			</Sidebar.Group>
		</Sidebar.Content>
		<Sidebar.Footer>
			<Sidebar.Group>
				<UserNav {user} />
			</Sidebar.Group>
		</Sidebar.Footer>
	</Sidebar.Root>
	<main class="container m-4 mx-auto inline-block flex-1 px-4">
		{@render children?.()}
	</main>
</Sidebar.Provider>
