<script lang="ts">
	import clsx from 'clsx/lite';

	const { children, class: className }: { children: any; class?: string; component?: any } =
		$props();
</script>

<p class={clsx('card-title', className)}>{@render children?.()}</p>
