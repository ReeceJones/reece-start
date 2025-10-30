<script lang="ts">
	import { t } from '$lib/i18n';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { enhance } from '$app/forms';
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
	import CardActions from '$lib/components/Card/CardActions.svelte';

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
	<Card>
		<CardBody>
			{#if loading}
				<div class="text-center">
					<span class="loading loading-lg loading-spinner"></span>
					<h2 class="mt-4">{$t('oauth.completingSignIn')}</h2>
					<p class="text-gray-500">{$t('oauth.pleaseWait')}</p>
				</div>
			{:else if error}
				<div class="text-center">
					<CardTitle class="text-error">{$t('oauth.authenticationError')}</CardTitle>
					<p class="mb-4 text-gray-500">{error}</p>
					<CardActions>
						<a href="/signin" class="btn w-full btn-primary">{$t('oauth.tryAgain')}</a>
					</CardActions>
				</div>
			{/if}
		</CardBody>
	</Card>
</main>
