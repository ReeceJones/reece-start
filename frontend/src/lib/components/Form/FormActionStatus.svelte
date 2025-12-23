<script lang="ts">
	import { CircleCheck, CircleX } from 'lucide-svelte';
	import * as Alert from '$lib/components/ui/alert';

	const {
		form,
		success,
		failure
	}: {
		form:
			| {
					success: boolean;
					message: string;
			  }
			| null
			| undefined;
		success: string;
		failure: string;
	} = $props();
</script>

{#if form?.success}
	<Alert.Root variant="success">
		<CircleCheck />
		<span>{success}</span>
	</Alert.Root>
{:else if form?.success === false}
	<Alert.Root variant="destructive">
		<CircleX />
		<span>{failure}</span>
	</Alert.Root>
{/if}
