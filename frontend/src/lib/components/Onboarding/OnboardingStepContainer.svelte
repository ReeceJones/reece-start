<script lang="ts">
	const { children, hidden }: { children: any; hidden: boolean } = $props();
</script>

<div {hidden}>
	{@render children?.()}
</div>
