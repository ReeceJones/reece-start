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
	import { t } from '$lib/i18n';

	let { form }: PageProps = $props();

	let submittingPasswordForm = $state(false);
	let submittingGoogleOAuthForm = $state(false);
	const submitting = $derived(submittingPasswordForm || submittingGoogleOAuthForm);

	const redirect = $derived(page.url.searchParams.get('redirect'));
	const googleClientId = $derived(env.PUBLIC_GOOGLE_OAUTH_CLIENT_ID || '');
</script>

<svelte:head>
	<title>{$t('auth.signUp.title')} - reece-start</title>
	<meta name="description" content={$t('auth.signUp.description')} />
</svelte:head>

<main class="mx-auto my-8 max-w-80">
	<Card>
		<CardBody>
			<CardTitle>{$t('auth.signUp.title')}</CardTitle>
			<p class="text-gray-500">{$t('auth.signUp.description')}</p>
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
							{$t('auth.signUp.signUpWithGoogle')}
						{/if}
					</button>
				</form>

				<div class="divider mb-2 mt-5 text-gray-500">{$t('auth.signUp.orContinueWith')}</div>

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
						<legend class="fieldset-legend">{$t('auth.signUp.name')}</legend>
						<input
							type="text"
							name="name"
							required
							class="validator input w-full"
							placeholder={$t('auth.signUp.name')}
						/>
					</fieldset>

					<fieldset class="fieldset">
						<legend class="fieldset-legend">{$t('auth.signUp.email')}</legend>
						<input
							type="email"
							name="email"
							required
							class="validator input w-full"
							placeholder={$t('auth.signUp.email')}
						/>
					</fieldset>

					<fieldset class="fieldset">
						<legend class="fieldset-legend">{$t('auth.signUp.password')}</legend>
						<input
							type="password"
							name="password"
							required
							class="validator input w-full"
							placeholder={$t('auth.signUp.password')}
						/>
					</fieldset>

					<div class="mt-3 space-y-3">
						<CardActions>
							<button type="submit" class="btn btn-primary mt-3 w-full" disabled={submitting}>
								{#if submittingPasswordForm}
									<span class="loading loading-spinner"></span>
								{:else}
									<LogIn />
								{/if}
								<span>{$t('auth.signUp.signUpButton')}</span>
							</button>
						</CardActions>

						{#if form?.success}
							<div role="alert" class="alert alert-success">
								<CircleCheck />
								<span>{$t('auth.signUp.successMessage')}</span>
							</div>
						{:else if form?.success === false}
							<div role="alert" class="alert alert-error">
								<CircleX />
								<span>{$t('auth.signUp.errorMessage')}</span>
							</div>
						{/if}

						<div class="mt-3 text-center text-sm">
							<p>
								{$t('auth.signUp.hasAccount')}
								<a href="/signin{redirect ? `?redirect=${redirect}` : ''}" class="link"
									>{$t('auth.signUp.signInLink')}</a
								>
							</p>
						</div>
					</div>
				</form>
			</div>
		</CardBody>
	</Card>
</main>
