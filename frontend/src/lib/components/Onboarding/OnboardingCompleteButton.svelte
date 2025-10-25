<script lang="ts">
	import { Plus } from 'lucide-svelte';
	import { fade } from 'svelte/transition';
	import { t } from '$lib/i18n';

	const {
		loading
	}: {
		loading: boolean;
	} = $props();
</script>

<button class="btn btn-primary" type="submit" in:fade disabled={loading}>
	{#if loading}
		<span class="loading loading-xs loading-spinner"></span>
	{:else}
		<Plus class="size-4" />
	{/if}
	{$t('createOrganization')}
</button>
