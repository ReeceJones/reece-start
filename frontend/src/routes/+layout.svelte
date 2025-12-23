<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { initializeLocale } from '$lib/i18n';
	import type { LayoutData } from './$types';
	import type { Snippet } from 'svelte';
	import { Toaster } from '$lib/components/ui/sonner';

	let { children, data }: { children: Snippet; data: LayoutData } = $props();

	// Initialize locale from server-provided value and keep it in sync
	$effect(() => {
		initializeLocale(data.locale);
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Toaster />

{@render children?.()}
