<script lang="ts">
	import { CircleCheck, CircleX, LogIn } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import Google from '$lib/components/Icons/Google.svelte';
	import { env } from '$env/dynamic/public';
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
	import CardActions from '$lib/components/Card/CardActions.svelte';

	let { form }: PageProps = $props();

	let submittingPasswordForm = $state(false);
	let submittingGoogleOAuthForm = $state(false);
	const submitting = $derived(submittingPasswordForm || submittingGoogleOAuthForm);

	const redirect = $derived(page.url.searchParams.get('redirect'));
	const googleClientId = $derived(env.PUBLIC_GOOGLE_OAUTH_CLIENT_ID || '');
</script>

<svelte:head>
	<title>Sign up - reece-start</title>
	<meta name="description" content="Sign up for an account to continue to the dashboard." />
</svelte:head>

<main class="mx-auto my-8 max-w-80">
	<Card>
		<CardBody>
			<CardTitle>Sign up</CardTitle>
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
					<button class="btn w-full btn-neutral" disabled={!googleClientId || submitting}>
						{#if submittingGoogleOAuthForm}
							<span class="loading loading-spinner"></span>
						{:else}
							<Google />
							Sign up with Google
						{/if}
					</button>
				</form>

				<div class="divider mt-5 mb-2 text-gray-500">Or continue with</div>

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
							class="validator input w-full"
							placeholder="Name"
						/>
					</fieldset>

					<fieldset class="fieldset">
						<legend class="fieldset-legend">Email</legend>
						<input
							type="email"
							name="email"
							required
							class="validator input w-full"
							placeholder="Email"
						/>
					</fieldset>

					<fieldset class="fieldset">
						<legend class="fieldset-legend">Password</legend>
						<input
							type="password"
							name="password"
							required
							class="validator input w-full"
							placeholder="Password"
						/>
					</fieldset>

					<div class="mt-3 space-y-3">
						<CardActions>
							<button type="submit" class="btn mt-3 w-full btn-primary" disabled={submitting}>
								{#if submittingPasswordForm}
									<span class="loading loading-spinner"></span>
								{:else}
									<LogIn />
								{/if}
								<span>Sign up</span>
							</button>
						</CardActions>

						{#if form?.success}
							<div role="alert" class="alert alert-success">
								<CircleCheck />
								<span
									>You have been signed up successfully! You will be redirected to the dashboard
									soon.</span
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
		</CardBody>
	</Card>
</main>
