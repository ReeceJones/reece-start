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
	import * as m from '$lib/paraglide/messages';

	let { form }: PageProps = $props();

	let submittingPasswordForm = $state(false);
	let submittingGoogleOAuthForm = $state(false);
	const submitting = $derived(submittingPasswordForm || submittingGoogleOAuthForm);

	const redirect = $derived(page.url.searchParams.get('redirect'));
	const googleClientId = $derived(env.PUBLIC_GOOGLE_OAUTH_CLIENT_ID || '');
	const isSignInDisabled = $derived(env.PUBLIC_DISABLE_SIGNIN === 'true');
</script>

<svelte:head>
	<title>{m.auth__sign_in__title()} - reece-start</title>
	<meta name="description" content={m.auth__sign_in__description()} />
</svelte:head>

<main class="mx-auto my-8 max-w-80">
	<Card.Root>
		<Card.Header>
			<Card.Title>{m.auth__sign_in__title()}</Card.Title>
			<Card.Description class="text-gray-500">{m.auth__sign_in__description()}</Card.Description>
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
							{m.auth__sign_in__sign_in_with_google()}
						{/if}
					</Button>
				</form>

				<div class="my-5 flex items-center gap-2">
					<Separator class="flex-1" />
					<span class="text-sm text-muted-foreground">{m.auth__sign_in__or_continue_with()}</span>
					<Separator class="flex-1" />
				</div>

				<form
					method="post"
					action="?/signin"
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
						<Field.Label for="email">{m.auth__sign_in__email()}</Field.Label>
						<Input
							id="email"
							type="email"
							name="email"
							required
							placeholder={m.auth__sign_in__email()}
						/>
					</Field.Field>

					<Field.Field>
						<Field.Label for="password">{m.auth__sign_in__password()}</Field.Label>
						<Input
							id="password"
							type="password"
							name="password"
							required
							placeholder={m.auth__sign_in__password()}
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
							<span>{m.auth__sign_in__sign_in_button()}</span>
						</Button>

						{#if form?.success}
							<Alert.Root>
								<CircleCheck class="h-4 w-4" />
								<Alert.Description>{m.auth__sign_in__success_message()}</Alert.Description>
							</Alert.Root>
						{:else if form?.success === false}
							<Alert.Root variant="destructive">
								<CircleX class="h-4 w-4" />
								<Alert.Description>
									{(form as { success: boolean; message: string })?.message ??
										m.auth__sign_in__error_message()}
								</Alert.Description>
							</Alert.Root>
						{/if}

						<div class="mt-3 text-center text-sm">
							<p>
								{m.auth__sign_in__no_account()}
								<Link href="/signup{redirect ? `?redirect=${redirect}` : ''}">
									{m.auth__sign_in__sign_up_link()}
								</Link>
							</p>
						</div>
					</div>
				</form>
			</div>
		</Card.Content>
	</Card.Root>
</main>
