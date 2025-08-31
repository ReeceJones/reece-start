<script lang="ts">
	import { CircleCheck, CircleX, LogIn } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { enhance, applyAction } from '$app/forms';
	import { page } from '$app/state';

	let { form }: PageProps = $props();
	let submitting = $state(false);
	const redirect = $derived(page.url.searchParams.get('redirect'));
</script>

<main class="card card-border bg-base-200 mx-auto my-8 max-w-80 shadow-sm">
	<div class="card-body">
		<h2 class="card-title">Sign in</h2>
		<form
			method="post"
			use:enhance={() => {
				submitting = true;

				return ({ result }) => {
					applyAction(result);
					submitting = false;
				};
			}}
		>
			<fieldset class="fieldset">
				<legend class="fieldset-legend">Email</legend>
				<input
					type="email"
					name="email"
					required
					class="input validator w-full"
					placeholder="Email"
				/>
			</fieldset>

			<fieldset class="fieldset">
				<legend class="fieldset-legend">Password</legend>
				<input
					type="password"
					name="password"
					required
					class="input validator w-full"
					placeholder="Password"
				/>
			</fieldset>

			<div class="mt-3 space-y-3">
				<div class="card-actions">
					<button type="submit" class="btn btn-primary w-full" disabled={submitting}>
						{#if submitting}
							<span class="loading loading-spinner"></span>
						{:else}
							<LogIn />
						{/if}
						<span>Sign in</span>
					</button>
				</div>

				{#if form?.success}
					<div role="alert" class="alert alert-success">
						<CircleCheck />
						<span
							>You have been signed in successfully! You will be redirected to the dashboard soon.</span
						>
					</div>
				{:else if form?.success === false}
					<div role="alert" class="alert alert-error">
						<CircleX />
						<span>
							{(form as { success: boolean; message: string })?.message ??
								'There was an error signing in. Make sure you have filled out all the fields correctly.'}
						</span>
					</div>
				{/if}

				<div class="mt-3 text-center text-sm">
					<p>
						Don't have an account? <a
							href="/signup{redirect ? `?redirect=${redirect}` : ''}"
							class="link">Sign up</a
						>
					</p>
				</div>
			</div>
		</form>
	</div>
</main>
