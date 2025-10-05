<script lang="ts">
	import { onMount } from 'svelte';
	import { enhance } from '$app/forms';
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
	import { CircleX } from 'lucide-svelte';

	let formEl: HTMLFormElement;
	let error = $state('');
	let loading = $state(false);

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
			{#if loading}
				<span class="loading mx-auto loading-lg loading-spinner"></span>
			{:else if error}
				<div class="alert alert-error">
					<CircleX />
					<span>{error}</span>
				</div>
			{/if}
		</CardBody>
	</Card>
	<form
		method="post"
		class="hidden"
		bind:this={formEl}
		use:enhance={() => {
			loading = true;
			return ({ result, update }) => {
				update();
				if (result.type === 'error') {
					error = result.error ?? 'Something went wrong. Please try again.';
					loading = false;
				} else if (result.type === 'failure') {
					error = (result.data?.message as string) ?? 'Something went wrong. Please try again.';
					loading = false;
				}
			};
		}}
	></form>
</div>
