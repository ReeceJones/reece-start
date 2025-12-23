<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Spinner } from '../ui/spinner';

	const {
		loading,
		icon
	}: {
		loading: boolean;
		icon: Snippet<[]>;
	} = $props();
</script>

{#if loading}
	<Spinner />
{:else}
	{@render icon()}
{/if}
