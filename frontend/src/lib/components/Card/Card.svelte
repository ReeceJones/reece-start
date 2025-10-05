<script lang="ts">
	import clsx from 'clsx/lite';
	import type { Snippet } from 'svelte';

	const { children, class: className }: { children: Snippet; class?: string } = $props();
</script>

<div class={clsx('card rounded-box border border-base-300 bg-base-200 shadow-sm', className)}>
	{@render children?.()}
</div>
