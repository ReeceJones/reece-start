<script lang="ts">
	import { getScopes } from '$lib/auth';
	import * as Card from '$lib/components/ui/card';
	import * as Select from '$lib/components/ui/select';
	import * as Alert from '$lib/components/ui/alert';
	import * as m from '$lib/paraglide/messages';
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import { getLocale, locales } from '$lib/paraglide/runtime';

	const scopes = $derived(getScopes());
	let submitting = $state(false);
	let error = $state('');

	let selectedLocale = $state(getLocale());
	let formElement: HTMLFormElement;
	let previousLocale = $state(getLocale());

	// Auto-submit form when locale changes
	$effect(() => {
		if (!formElement || selectedLocale === previousLocale) {
			return;
		}
		previousLocale = selectedLocale;
		formElement.requestSubmit();
	});
</script>

<div class="grid grid-cols-3 gap-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>{m.no_organization__admin__debug__user_scopes()}</Card.Title>
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
