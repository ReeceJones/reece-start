<script lang="ts">
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
</script>

<div class="grid grid-cols-6 gap-6">
	<Card class="col-span-4">
		<CardBody>
			<CardTitle>Foo</CardTitle>
		</CardBody>
	</Card>
	<Card class="col-span-2">
		<CardBody>
			<CardTitle>Bar</CardTitle>
		</CardBody>
	</Card>
</div>
