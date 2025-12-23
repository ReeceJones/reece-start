<script lang="ts">
	import UserNav from '$lib/components/Nav/UserNav.svelte';
	import IndexNav from '$lib/components/Nav/IndexNav.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';

	const { children, data } = $props();

	const { user } = $derived(data);
</script>

<Sidebar.Provider>
	<Sidebar.Root>
		<Sidebar.Content>
			<Sidebar.Group>
				<IndexNav />
			</Sidebar.Group>
		</Sidebar.Content>
		<Sidebar.Footer>
			<Sidebar.Group>
				<UserNav {user} />
			</Sidebar.Group>
		</Sidebar.Footer>
	</Sidebar.Root>
	<main class="container mx-auto mt-4 inline-block flex-1 px-4">
		{@render children?.()}
	</main>
</Sidebar.Provider>
