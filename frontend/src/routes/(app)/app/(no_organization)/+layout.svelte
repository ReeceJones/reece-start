<script lang="ts">
	import UserNav from '$lib/components/UserNav.svelte';
	import IndexNav from '$lib/components/IndexNav.svelte';

	const { children, data } = $props();

	const { user } = data;
</script>

<div class="max-w-screen flex h-screen flex-row gap-4">
	<div class="border-base-300 bg-base-200 flex h-full w-56 flex-col justify-between border-r">
		<div>
			<IndexNav />
		</div>
		<div>
			<UserNav {user} />
		</div>
	</div>
	<main class="container mx-auto mr-4 mt-4 flex-1">
		{@render children?.()}
	</main>
</div>
