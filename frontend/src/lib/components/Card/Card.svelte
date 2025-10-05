<script lang="ts">
	import clsx from 'clsx/lite';

	const { children, class: className }: { children: any; class?: string } = $props();
</script>

<div class={clsx('card rounded-box bg-base-200 border-base-300 border shadow-sm', className)}>
	{@render children?.()}
</div>
