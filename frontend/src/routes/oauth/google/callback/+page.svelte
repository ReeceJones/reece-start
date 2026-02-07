<script lang="ts">
	import * as m from '$lib/paraglide/messages';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { enhance } from '$app/forms';
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Spinner } from '$lib/components/ui/spinner';

	let loading = $state(true);
	let error = $state('');
	let oauthForm: HTMLFormElement;

	onMount(() => {
		const code = page.url.searchParams.get('code');
		const state = page.url.searchParams.get('state');
		const errorParam = page.url.searchParams.get('error');

		if (errorParam) {
			error = 'OAuth authentication was cancelled or failed.';
			loading = false;
			return;
		}

		if (!code || !state) {
			error = 'Missing required OAuth parameters.';
			loading = false;
			return;
		}

		if (!oauthForm) {
			console.error('OAuth form not found');
			return;
		}

		const codeInput = oauthForm.querySelector('input[name="code"]') as HTMLInputElement;
		const stateInput = oauthForm.querySelector('input[name="state"]') as HTMLInputElement;

		if (codeInput && stateInput) {
			codeInput.value = code;
			stateInput.value = state;
			oauthForm.submit();
		}
	});
</script>

<!-- Hidden form for OAuth processing -->
<form
	bind:this={oauthForm}
	method="post"
	class="hidden"
	use:enhance={() => {
		loading = true;
		return ({ result }) => {
			if (result.type === 'redirect') {
				// Redirect will be handled automatically
				return;
			}
			if (result.type === 'failure') {
				error = (result.data?.message as string) || 'Authentication failed. Please try again.';
				loading = false;
			}
		};
	}}
>
	<input type="hidden" name="code" value="" />
	<input type="hidden" name="state" value="" />
	<input type="hidden" name="redirect" value="" />
</form>

<main class="mx-auto my-8 max-w-80">
	<Card.Root>
		<Card.Content>
			{#if loading}
				<div class="flex flex-col items-center justify-center gap-4 text-center">
					<Spinner class="size-8" />
					<Card.Title class="mt-4">{m.oauth__completing_sign_in()}</Card.Title>
					<Card.Description>{m.oauth__please_wait()}</Card.Description>
				</div>
			{:else if error}
				<div class="flex flex-col items-center gap-4 text-center">
					<Card.Title class="text-destructive">{m.oauth__authentication_error()}</Card.Title>
					<Card.Description class="mb-4">{error}</Card.Description>
					<Card.Footer class="w-full">
						<Button href="/signin" class="w-full" variant="default">
							{m.oauth__try_again()}
						</Button>
					</Card.Footer>
				</div>
			{/if}
		</Card.Content>
	</Card.Root>
</main>
