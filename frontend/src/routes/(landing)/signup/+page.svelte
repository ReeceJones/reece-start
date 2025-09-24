<script lang="ts">
	import { CircleCheck, CircleX, LogIn } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import Google from '$lib/components/Icons/Google.svelte';
	import { env } from '$env/dynamic/public';

	let { form }: PageProps = $props();

	let submittingPasswordForm = $state(false);
	let submittingGoogleOAuthForm = $state(false);
	const submitting = $derived(submittingPasswordForm || submittingGoogleOAuthForm);

	const redirect = $derived(page.url.searchParams.get('redirect'));
	const googleClientId = $derived(env.PUBLIC_GOOGLE_OAUTH_CLIENT_ID || '');
</script>

<main class="card card-border bg-base-200 mx-auto my-8 max-w-80 shadow-sm">
	<div class="card-body">
		<h2 class="card-title">Sign up</h2>
		<p class="text-gray-500">Enter your details below to sign up for an account.</p>
		<div class="mt-4">
			<form
				method="post"
				action="?/oauthGoogle"
				use:enhance={() => {
					submittingGoogleOAuthForm = true;

					return ({ update }) => {
						update();
						submittingGoogleOAuthForm = false;
					};
				}}
			>
				<button class="btn btn-neutral w-full" disabled={!googleClientId || submitting}>
					{#if submittingGoogleOAuthForm}
						<span class="loading loading-spinner"></span>
					{:else}
						<Google />
						Sign up with Google
					{/if}
				</button>
			</form>

			<div class="divider mb-2 mt-5 text-gray-500">Or continue with</div>

			<form
				method="post"
				action="?/signup"
				use:enhance={() => {
					submittingPasswordForm = true;

					return ({ update }) => {
						update();
						submittingPasswordForm = false;
					};
				}}
			>
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Name</legend>
					<input
						type="text"
						name="name"
						required
						class="input validator w-full"
						placeholder="Name"
					/>
				</fieldset>

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
						<button type="submit" class="btn btn-primary mt-3 w-full" disabled={submitting}>
							{#if submittingPasswordForm}
								<span class="loading loading-spinner"></span>
							{:else}
								<LogIn />
							{/if}
							<span>Sign up</span>
						</button>
					</div>

					{#if form?.success}
						<div role="alert" class="alert alert-success">
							<CircleCheck />
							<span
								>You have been signed up successfully! You will be redirected to the dashboard soon.</span
							>
						</div>
					{:else if form?.success === false}
						<div role="alert" class="alert alert-error">
							<CircleX />
							<span
								>There was an error signing up. Make sure you have filled out all the fields
								correctly.</span
							>
						</div>
					{/if}

					<div class="mt-3 text-center text-sm">
						<p>
							Already have an account? <a
								href="/signin{redirect ? `?redirect=${redirect}` : ''}"
								class="link">Sign in</a
							>
						</p>
					</div>
				</div>
			</form>
		</div>
	</div>
</main>
