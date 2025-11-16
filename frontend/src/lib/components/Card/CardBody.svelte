<script lang="ts">
	import clsx from 'clsx/lite';
	import type { Snippet } from 'svelte';

	const { children, class: className }: { children: Snippet; class?: string } = $props();
</script>

<div class={clsx('card-body', className)}>
	{@render children?.()}
</div>
