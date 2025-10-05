<script lang="ts">
	import clsx from 'clsx/lite';
	import type { Snippet } from 'svelte';

	const { children, class: className }: { children: Snippet; class?: string } = $props();
</script>

<p class={clsx('card-title', className)}>{@render children?.()}</p>
