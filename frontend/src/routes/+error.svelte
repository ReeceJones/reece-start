<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import { House, RefreshCw, TriangleAlert } from 'lucide-svelte';
	import { page } from '$app/state';
	import { Button, buttonVariants } from '$lib/components/ui/button';

	const status = $derived(page.status);
	const error = $derived(page.error);

	const errorMessage = $derived.by(() => {
		if (error instanceof Error) return error.message;
		if (typeof error === 'string') return error;
		if (error && typeof error === 'object' && 'message' in error) {
			return String(error.message);
		}
		return 'An unexpected error occurred';
	});

	const errorTitle = $derived.by(() => {
		if (status === 404) return 'Page Not Found';
		if (status === 401) return 'Unauthorized';
		if (status === 403) return 'Forbidden';
		if (status === 500) return 'Internal Server Error';
		if (status) return `Error ${status}`;
		return 'Error';
	});

	const errorDescription = $derived.by(() => {
		if (status === 404) return 'The page you are looking for could not be found.';
		if (status === 401) return 'You need to be authenticated to access this resource.';
		if (status === 403) return "You don't have permission to access this resource.";
		if (status === 500) return 'An internal server error occurred. Please try again later.';
		return 'Something went wrong. Please try again later.';
	});
</script>

<div class="bg-base-100 flex min-h-screen items-center justify-center p-4">
	<Card.Root class="w-full max-w-md">
		<Card.Header>
			<div class="bg-error/10 mx-auto mb-4 flex size-16 items-center justify-center rounded-full">
				<TriangleAlert class="text-error size-8" />
			</div>
			<Card.Title>{errorTitle}</Card.Title>
			<Card.Description>{errorDescription}</Card.Description>
		</Card.Header>
		{#if errorMessage && errorMessage !== errorDescription && status !== 404}
			<Card.Content>
				<Alert.Root variant="destructive">
					<Alert.Description class="overflow-scroll font-semibold">
						{errorMessage}</Alert.Description
					>
				</Alert.Root>
			</Card.Content>
		{/if}
		<Card.Footer>
			<Card.Action>
				<a href="/" class={buttonVariants()}>
					<House class="size-4" />
					Go Home
				</a>
				{#if status !== 404}
					<Button
						variant="ghost"
						onclick={() => {
							window.location.reload();
						}}
					>
						<RefreshCw class="size-4" />
						Reload Page
					</Button>
				{/if}
			</Card.Action>
		</Card.Footer>
	</Card.Root>
</div>
