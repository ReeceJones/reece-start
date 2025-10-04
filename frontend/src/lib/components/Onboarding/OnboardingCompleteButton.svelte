<script lang="ts">
	import { Plus } from 'lucide-svelte';
	import { fade } from 'svelte/transition';

	const {
		loading
	}: {
		loading: boolean;
	} = $props();
</script>

<button class="btn btn-primary" type="submit" in:fade disabled={loading}>
	{#if loading}
		<span class="loading loading-spinner loading-xs"></span>
	{:else}
		<Plus class="size-4" />
	{/if}
	Create Organization
</button>
