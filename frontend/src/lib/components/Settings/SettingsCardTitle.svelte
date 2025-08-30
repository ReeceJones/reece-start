<script lang="ts">
	const { children } = $props();
</script>

<p class="card-title mb-4">{@render children?.()}</p>
