<script lang="ts">
	import clsx from 'clsx/lite';
	const { children, class: className }: { children: any; class?: string } = $props();
</script>

<div class={clsx('card-actions justify-end', className)}>
	{@render children?.()}
</div>
