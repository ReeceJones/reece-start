<script lang="ts">
	import { getScopes } from '$lib/auth';
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
	import { t, locale, locales } from '$lib/i18n';
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';

	const scopes = $derived(getScopes());
	let submitting = $state(false);
	let error = $state('');
	let selectedLocale = $state($locale);

	// Sync selectedLocale with locale store
	$effect(() => {
		selectedLocale = $locale;
	});
</script>

<div class="grid grid-cols-3 gap-6">
	<Card>
		<CardBody>
			<CardTitle>{$t('noOrganization.admin.debug.userScopes')}</CardTitle>
			<ul class="list-inside list-disc">
				{#each scopes as scope (scope)}
					<li class="font-mono">{scope}</li>
				{/each}
			</ul>
		</CardBody>
	</Card>

	<Card>
		<CardBody>
			<CardTitle>Locale</CardTitle>
			<form
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
				<select
					name="locale"
					class="select select-bordered w-full"
					bind:value={selectedLocale}
					onchange={(e) => {
						const form = e.currentTarget.form;
						if (form) {
							form.requestSubmit();
						}
					}}
					disabled={submitting}
				>
					{#each locales as localeOption (localeOption)}
						<option value={localeOption}>{localeOption}</option>
					{/each}
				</select>
				{#if error}
					<div role="alert" class="alert alert-error mt-2">
						<span>{error}</span>
					</div>
				{/if}
			</form>
		</CardBody>
	</Card>
</div>
