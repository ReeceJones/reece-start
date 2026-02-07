<script lang="ts">
	import * as m from '$lib/paraglide/messages';
	import { onMount } from 'svelte';
	import { enhance } from '$app/forms';
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import { CircleX } from 'lucide-svelte';
	import { Spinner } from '$lib/components/ui/spinner';

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
	<Card.Root>
		<Card.Header>
			<Card.Title>{m.payments__redirecting_to_stripe()}</Card.Title>
		</Card.Header>
		<Card.Content>
			{#if loading}
				<Spinner class="mx-auto size-6" />
			{:else if error}
				<Alert.Root variant="destructive">
					<CircleX />
					<Alert.Description>
						{error}
					</Alert.Description>
				</Alert.Root>
			{/if}
		</Card.Content>
	</Card.Root>
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
