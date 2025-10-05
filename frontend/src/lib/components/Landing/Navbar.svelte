<script lang="ts">
	import { DoorOpen, Rocket, Menu, X } from 'lucide-svelte';

	const { isLoggedIn }: { isLoggedIn: boolean } = $props();

	let isMobileMenuOpen = $state(false);

	function toggleMobileMenu() {
		isMobileMenuOpen = !isMobileMenuOpen;
	}

	function closeMobileMenu() {
		isMobileMenuOpen = false;
	}
</script>

<div class="navbar">
	<div class="flex-1">
		<a class="btn btn-ghost text-lg tracking-tight" href="/">
			<Rocket class="size-6" />
			reece-start
		</a>
	</div>

	<!-- Desktop Navigation -->
	<div class="hidden lg:flex lg:flex-none lg:items-center lg:gap-2">
		<a href="/faq" class="btn btn-ghost font-medium">FAQ</a>
		<a href="/pricing" class="btn btn-ghost font-medium">Pricing</a>
		{#if isLoggedIn}
			<a href="/app" class="btn btn-neutral font-medium">
				<DoorOpen class="size-5" />
				Dashboard
			</a>
		{:else}
			<a href="/signin" class="btn btn-outline btn-neutral font-medium"> Sign in </a>
			<a href="/signup" class="btn btn-neutral font-medium"> Get started </a>
		{/if}
	</div>

	<!-- Mobile Hamburger Menu -->
	<div class="lg:hidden">
		<button class="btn btn-ghost btn-sm" aria-label="Toggle mobile menu" onclick={toggleMobileMenu}>
			{#if isMobileMenuOpen}
				<X class="size-6" />
			{:else}
				<Menu class="size-6" />
			{/if}
		</button>
	</div>
</div>

<!-- Mobile Menu Dropdown -->
{#if isMobileMenuOpen}
	<div class="bg-base-100 border-base-300 border-t shadow-lg lg:hidden">
		<div class="space-y-3 px-4 py-6">
			<a
				href="/faq"
				class="hover:bg-base-200 block rounded-lg px-3 py-2 text-base font-medium transition-colors"
				onclick={closeMobileMenu}
			>
				FAQ
			</a>
			<a
				href="/pricing"
				class="hover:bg-base-200 block rounded-lg px-3 py-2 text-base font-medium transition-colors"
				onclick={closeMobileMenu}
			>
				Pricing
			</a>
			{#if isLoggedIn}
				<a
					href="/app"
					class="hover:bg-base-200 block rounded-lg px-3 py-2 text-base font-medium transition-colors"
					onclick={closeMobileMenu}
				>
					<DoorOpen class="mr-2 inline size-5" />
					Dashboard
				</a>
			{:else}
				<a
					href="/signin"
					class="hover:bg-base-200 block rounded-lg px-3 py-2 text-base font-medium transition-colors"
					onclick={closeMobileMenu}
				>
					Sign in
				</a>
				<a
					href="/signup"
					class="bg-neutral text-neutral-content hover:bg-neutral-focus block rounded-lg px-3 py-2 text-base font-medium transition-colors"
					onclick={closeMobileMenu}
				>
					Get started
				</a>
			{/if}
		</div>
	</div>
{/if}
