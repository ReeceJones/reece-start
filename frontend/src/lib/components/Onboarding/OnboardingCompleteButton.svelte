<script lang="ts">
	import { Plus } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
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
	<span>{m.create_organization()}</span>
</Button>
