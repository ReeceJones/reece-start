<script lang="ts">
	import { ArrowLeft } from 'lucide-svelte';
	import { fade } from 'svelte/transition';
	import clsx from 'clsx/lite';

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

<a class={clsx('btn btn-outline', disabled && 'pointer-events-none')} href={step.path} in:fade>
	<ArrowLeft class="size-4" />
	Back
</a>
