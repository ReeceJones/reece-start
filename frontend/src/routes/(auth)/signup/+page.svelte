<script lang="ts">
	import { CircleCheck, CircleX, LogIn } from 'lucide-svelte';
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import Google from '$lib/components/Icons/Google.svelte';
	import { env } from '$env/dynamic/public';
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Field from '$lib/components/ui/field';
	import { Separator } from '$lib/components/ui/separator';
	import * as Alert from '$lib/components/ui/alert';
	import { Link } from '$lib/components/ui/link';
	import { Spinner } from '$lib/components/ui/spinner';
	import { t } from '$lib/i18n';

	let { form }: PageProps = $props();

	let submittingPasswordForm = $state(false);
	let submittingGoogleOAuthForm = $state(false);
	const submitting = $derived(submittingPasswordForm || submittingGoogleOAuthForm);

	const redirect = $derived(page.url.searchParams.get('redirect'));
	const googleClientId = $derived(env.PUBLIC_GOOGLE_OAUTH_CLIENT_ID || '');
	const isSignInDisabled = $derived(env.PUBLIC_DISABLE_SIGNIN === 'true');
</script>

<svelte:head>
	<title>{$t('auth.signUp.title')} - reece-start</title>
	<meta name="description" content={$t('auth.signUp.description')} />
</svelte:head>

<main class="mx-auto my-8 max-w-80">
	<Card.Root>
		<Card.Header>
			<Card.Title>{$t('auth.signUp.title')}</Card.Title>
			<Card.Description class="text-muted-foreground"
				>{$t('auth.signUp.description')}</Card.Description
			>
		</Card.Header>
		<Card.Content>
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
					<Button
						class="w-full"
						variant="outline"
						type="submit"
						disabled={!googleClientId || submitting || isSignInDisabled}
					>
						{#if submittingGoogleOAuthForm}
							<Spinner class="h-4 w-4" />
						{:else}
							<Google />
							{$t('auth.signUp.signUpWithGoogle')}
						{/if}
					</Button>
				</form>

				<div class="my-5 flex items-center gap-2">
					<Separator class="flex-1" />
					<span class="text-sm text-muted-foreground">{$t('auth.signUp.orContinueWith')}</span>
					<Separator class="flex-1" />
				</div>

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
					class="space-y-4"
				>
					<Field.Field>
						<Field.Label for="name">{$t('auth.signUp.name')}</Field.Label>
						<Input
							id="name"
							type="text"
							name="name"
							required
							placeholder={$t('auth.signUp.name')}
						/>
					</Field.Field>

					<Field.Field>
						<Field.Label for="email">{$t('auth.signUp.email')}</Field.Label>
						<Input
							id="email"
							type="email"
							name="email"
							required
							placeholder={$t('auth.signUp.email')}
						/>
					</Field.Field>

					<Field.Field>
						<Field.Label for="password">{$t('auth.signUp.password')}</Field.Label>
						<Input
							id="password"
							type="password"
							name="password"
							required
							placeholder={$t('auth.signUp.password')}
						/>
					</Field.Field>

					<div class="mt-3 space-y-3">
						<Button
							type="submit"
							variant="default"
							class="w-full"
							disabled={submitting || isSignInDisabled}
						>
							{#if submittingPasswordForm}
								<Spinner class="h-4 w-4" />
							{:else}
								<LogIn class="h-4 w-4" />
							{/if}
							<span>{$t('auth.signUp.signUpButton')}</span>
						</Button>

						{#if form?.success}
							<Alert.Root>
								<CircleCheck class="h-4 w-4" />
								<Alert.Description>{$t('auth.signUp.successMessage')}</Alert.Description>
							</Alert.Root>
						{:else if form?.success === false}
							<Alert.Root variant="destructive">
								<CircleX class="h-4 w-4" />
								<Alert.Description>{$t('auth.signUp.errorMessage')}</Alert.Description>
							</Alert.Root>
						{/if}

						<div class="mt-3 text-center text-sm">
							<p>
								{$t('auth.signUp.hasAccount')}
								<Link href="/signin{redirect ? `?redirect=${redirect}` : ''}">
									{$t('auth.signUp.signInLink')}
								</Link>
							</p>
						</div>
					</div>
				</form>
			</div>
		</Card.Content>
	</Card.Root>
</main>
