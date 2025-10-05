<script lang="ts">
	import Card from '../Card/Card.svelte';
	import CardBody from '../Card/CardBody.svelte';

	const { children } = $props();
</script>

<Card>
	<CardBody>
		{@render children?.()}
	</CardBody>
</Card>
