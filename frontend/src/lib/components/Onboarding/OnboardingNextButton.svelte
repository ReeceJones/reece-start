<script lang="ts">
	import { ArrowRight } from 'lucide-svelte';
	import { fade } from 'svelte/transition';
	import clsx from 'clsx/lite';
	import { t } from '$lib/i18n';

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

<a class={clsx('btn btn-neutral', disabled && 'pointer-events-none')} href={step.path} in:fade>
	<ArrowRight class="size-4" />
	{$t('next')}
</a>
