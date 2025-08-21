<script lang="ts">
	import UserNav from '$lib/components/UserNav.svelte';
	import IndexNav from '$lib/components/IndexNav.svelte';

	const { children } = $props();
</script>

<div class="container mx-auto">
	<div class="flex justify-between">
		<div>
			<IndexNav />
		</div>
		<div>
			<UserNav />
		</div>
	</div>
	<main class="mt-8">
		{@render children?.()}
	</main>
</div>
