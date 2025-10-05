<script lang="ts">
	import CardActions from '$lib/components/Card/CardActions.svelte';

	const { children } = $props();
</script>

<CardActions class="mt-3 justify-start">
	{@render children?.()}
</CardActions>
