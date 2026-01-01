<script lang="ts">
	import { getScopes } from '$lib/auth';
	import * as Card from '$lib/components/ui/card';
	import * as Select from '$lib/components/ui/select';
	import * as Alert from '$lib/components/ui/alert';
	import { t, locale, locales } from '$lib/i18n';
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';

	const scopes = $derived(getScopes());
	let submitting = $state(false);
	let error = $state('');

	let selectedLocale = $state($locale);
	let formElement: HTMLFormElement;
	let previousLocale = $state($locale);

	// Sync selectedLocale with locale store
	$effect(() => {
		selectedLocale = $locale;
		previousLocale = $locale;
	});

	// Auto-submit form when locale changes (but not on initial mount or store sync)
	$effect(() => {
		if (!formElement || selectedLocale === previousLocale) {
			return;
		}
		// Only submit if the change is user-initiated (selectedLocale differs from store)
		if (selectedLocale !== $locale) {
			previousLocale = selectedLocale;
			formElement.requestSubmit();
		}
	});
</script>

<div class="grid grid-cols-3 gap-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>{$t('noOrganization.admin.debug.userScopes')}</Card.Title>
		</Card.Header>
		<Card.Content>
			<ul class="list-inside list-disc text-sm">
				{#each scopes as scope (scope)}
					<li class="font-mono">{scope}</li>
				{/each}
			</ul>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Locale</Card.Title>
		</Card.Header>
		<Card.Content>
			<form
				bind:this={formElement}
				method="post"
				action="/locale?/setLocale"
				use:enhance={() => {
					submitting = true;
					error = '';
					return ({ result, update }) => {
						update();
						submitting = false;

						if (result.type === 'success') {
							// Invalidate all to reload with new locale
							invalidateAll();
						} else if (result.type === 'failure') {
							error = (result.data?.message as string) || 'Failed to update locale';
						}
					};
				}}
			>
				<Select.Root name="locale" type="single" bind:value={selectedLocale} disabled={submitting}>
					<Select.Trigger class="w-full">
						{selectedLocale}
					</Select.Trigger>
					<Select.Content>
						{#each locales as localeOption (localeOption)}
							<Select.Item value={localeOption}>{localeOption}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
				{#if error}
					<Alert.Root variant="destructive" class="mt-2">
						<Alert.Description>{error}</Alert.Description>
					</Alert.Root>
				{/if}
			</form>
		</Card.Content>
	</Card.Root>
</div>
