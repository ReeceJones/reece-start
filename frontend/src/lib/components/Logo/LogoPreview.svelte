<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	const {
		logoFile,
		logoUrl,
		fallback,
		alt,
		class: className
	}: {
		logoFile?: FileList | null | undefined;
		logoUrl?: string | null | undefined;
		fallback?: Snippet<[]>;
		alt?: string;
		class?: string;
	} = $props();

	const logoPreview = $derived(
		logoFile && logoFile.length > 0 ? URL.createObjectURL(logoFile[0]) : logoUrl
	);
</script>

<div
	class={cn(
		'flex aspect-square h-48 max-h-48 w-48 max-w-48 items-center justify-center rounded-lg bg-neutral-200',
		className
	)}
>
	{#if logoPreview}
		<img src={logoPreview} alt={alt ?? 'Logo'} class="overflow-hidden rounded-lg object-cover" />
	{:else if fallback}
		{@render fallback()}
	{/if}
</div>
