<script lang="ts">
	import { ArrowLeft } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
	import { Button } from '$lib/components/ui/button';

	const {
		step,
		disabled
	}: {
		step: {
			path: string;
		};
		disabled: boolean;
	} = $props();
</script>

<Button href={step.path} variant="outline" {disabled}>
	<ArrowLeft />
	<span>{m.back()}</span>
</Button>
