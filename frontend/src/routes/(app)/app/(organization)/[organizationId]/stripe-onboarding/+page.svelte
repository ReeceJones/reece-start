<script lang="ts">
	import { onMount } from 'svelte';
	import { enhance } from '$app/forms';
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';

	let formEl: HTMLFormElement;

	onMount(() => {
		if (formEl) {
			formEl.submit();
		}
	});
</script>

<div class="mx-auto max-w-80">
	<Card>
		<CardBody>
			<CardTitle>Redirecting to Stripe...</CardTitle>
			<span class="loading loading-spinner loading-lg mx-auto"></span>
		</CardBody>
	</Card>
	<form
		method="post"
		class="hidden"
		bind:this={formEl}
		use:enhance={() =>
			({ result }) => {
				/* handled by server redirect */
			}}
	></form>
</div>
