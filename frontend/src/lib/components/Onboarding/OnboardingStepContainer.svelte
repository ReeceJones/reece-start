<script lang="ts">
	import type { Snippet } from 'svelte';

	const { children, hidden }: { children: Snippet; hidden: boolean } = $props();
</script>

<div {hidden} class="space-y-3">
	{@render children?.()}
</div>
