<script lang="ts">
	const { children } = $props();
</script>

<div class="card-body">
	{@render children?.()}
</div>
