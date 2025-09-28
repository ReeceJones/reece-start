<script lang="ts">
	import { Plus } from 'lucide-svelte';
	import { fade } from 'svelte/transition';

	const {
		loading,
		disabled
	}: {
		loading: boolean;
		disabled: boolean;
	} = $props();
</script>

<button class="btn btn-accent" type="submit" in:fade disabled={disabled || loading}>
	{#if loading}
		<span class="loading loading-spinner loading-xs"></span>
	{:else}
		<Plus class="size-4" />
	{/if}
	Create Organization
</button>
