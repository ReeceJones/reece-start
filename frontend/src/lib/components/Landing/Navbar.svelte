<script lang="ts">
	import { DoorOpen, Rocket } from 'lucide-svelte';

	const { isLoggedIn }: { isLoggedIn: boolean } = $props();
</script>

<div class="navbar">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl tracking-tight" href="/">
			<Rocket class="size-6" />
			reece-start
		</a>
		<a href="/faq" class="btn btn-ghost text-base font-medium">FAQ</a>
		<a href="/pricing" class="btn btn-ghost text-base font-medium">Pricing</a>
	</div>
	<div class="flex-none">
		{#if isLoggedIn}
			<a href="/app" class="btn btn-neutral text-base font-medium">
				<DoorOpen class="size-5" />
				Dashboard
			</a>
		{:else}
			<a href="/signin" class="btn btn-ghost text-base font-medium"> Sign in </a>
			<a href="/signup" class="btn btn-neutral text-base font-medium"> Get started </a>
		{/if}
	</div>
</div>
