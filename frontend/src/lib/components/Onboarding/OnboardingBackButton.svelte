<script lang="ts">
	import { ArrowLeft } from 'lucide-svelte';
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

<Button href={step.path} variant="outline" {disabled}>
	<ArrowLeft />
	<span>{$t('back')}</span>
</Button>
