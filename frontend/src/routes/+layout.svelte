<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { assertIsLocale, setLocale } from '$lib/paraglide/runtime';
	import type { LayoutData } from './$types';
	import type { Snippet } from 'svelte';
	import { Toaster } from '$lib/components/ui/sonner';

	let { children, data }: { children: Snippet; data: LayoutData } = $props();

	// Keep client locale in sync with the server-provided locale.
	$effect(() => {
		setLocale(assertIsLocale(data.locale), { reload: false });
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Toaster />

{@render children?.()}
