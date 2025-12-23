<script lang="ts">
	import { Plus } from 'lucide-svelte';
	import { t } from '$lib/i18n';
	import { Button } from '$lib/components/ui/button';
	import LoadingIcon from '$lib/components/Icons/LoadingIcon.svelte';

	const {
		loading
	}: {
		loading: boolean;
	} = $props();
</script>

<Button type="submit" disabled={loading}>
	<LoadingIcon {loading}>
		{#snippet icon()}
			<Plus />
		{/snippet}
	</LoadingIcon>
	<span>{$t('createOrganization')}</span>
</Button>
