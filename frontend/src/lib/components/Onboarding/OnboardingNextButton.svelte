<script lang="ts">
	import { ArrowRight } from 'lucide-svelte';
	import { t } from '$lib/i18n';
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

<Button href={step.path} variant="default" {disabled}>
	<span>{$t('next')}</span>
	<ArrowRight />
</Button>
