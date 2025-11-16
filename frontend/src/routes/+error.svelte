<script lang="ts">
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
	import CardActions from '$lib/components/Card/CardActions.svelte';
	import { AlertCircle, Home, RefreshCw } from 'lucide-svelte';
	import { page } from '$app/state';

	const status = $derived(page.status);
	const error = $derived(page.error);

	const errorMessage = $derived(() => {
		if (error instanceof Error) return error.message;
		if (typeof error === 'string') return error;
		if (error && typeof error === 'object' && 'message' in error) {
			return String(error.message);
		}
		return 'An unexpected error occurred';
	});

	const errorTitle = $derived(() => {
		if (status === 404) return 'Page Not Found';
		if (status === 401) return 'Unauthorized';
		if (status === 403) return 'Forbidden';
		if (status === 500) return 'Internal Server Error';
		if (status) return `Error ${status}`;
		return 'Error';
	});

	const errorDescription = $derived(() => {
		if (status === 404) return 'The page you are looking for could not be found.';
		if (status === 401) return 'You need to be authenticated to access this resource.';
		if (status === 403) return "You don't have permission to access this resource.";
		if (status === 500) return 'An internal server error occurred. Please try again later.';
		return 'Something went wrong. Please try again later.';
	});
</script>

<div class="flex min-h-screen items-center justify-center bg-base-100 p-4">
	<Card class="w-full max-w-md">
		<CardBody class="text-center">
			<div class="mx-auto mb-4 flex size-16 items-center justify-center rounded-full bg-error/10">
				<AlertCircle class="size-8 text-error" />
			</div>
			<CardTitle class="text-2xl">{errorTitle}</CardTitle>
			<p class="mt-2 mb-6 text-base-content/70">{errorDescription}</p>
			{#if errorMessage && errorMessage !== errorDescription}
				<div class="mb-6 alert justify-start alert-error text-left">
					<AlertCircle class="size-5 shrink-0" />
					<span class="text-sm">{errorMessage}</span>
				</div>
			{/if}
			<CardActions class="justify-center gap-3">
				<a href="/" class="btn btn-primary">
					<Home class="size-4" />
					Go Home
				</a>
				<button
					class="btn btn-ghost"
					onclick={() => {
						window.location.reload();
					}}
				>
					<RefreshCw class="size-4" />
					Reload Page
				</button>
			</CardActions>
		</CardBody>
	</Card>
</div>
