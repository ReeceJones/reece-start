<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { enhance } from '$app/forms';

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
				error = (result.data as any)?.message || 'Authentication failed. Please try again.';
				loading = false;
			}
		};
	}}
>
	<input type="hidden" name="code" value="" />
	<input type="hidden" name="state" value="" />
	<input type="hidden" name="redirect" value="" />
</form>

<main class="card card-border bg-base-200 mx-auto my-8 max-w-80 shadow-sm">
	<div class="card-body">
		{#if loading}
			<div class="text-center">
				<span class="loading loading-spinner loading-lg"></span>
				<h2 class="mt-4">Completing sign in...</h2>
				<p class="text-gray-500">Please wait while we finish signing you in with Google.</p>
			</div>
		{:else if error}
			<div class="text-center">
				<h2 class="card-title text-error">Authentication Error</h2>
				<p class="mb-4 text-gray-500">{error}</p>
				<div class="card-actions">
					<a href="/signin" class="btn btn-primary">Try Again</a>
				</div>
			</div>
		{/if}
	</div>
</main>
