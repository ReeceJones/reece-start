<script lang="ts">
	const { children } = $props();
</script>

<p class="mb-4 card-title">{@render children?.()}</p>
