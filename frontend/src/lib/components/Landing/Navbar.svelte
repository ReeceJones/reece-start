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
		<a class="btn text-lg tracking-tight btn-ghost" href="/">
			<Rocket class="size-6" />
			reece-start
		</a>
	</div>

	<!-- Desktop Navigation -->
	<div class="hidden lg:flex lg:flex-none lg:items-center lg:gap-2">
		<a href="/faq" class="btn font-medium btn-ghost">FAQ</a>
		<a href="/pricing" class="btn font-medium btn-ghost">Pricing</a>
		{#if isLoggedIn}
			<a href="/app" class="btn font-medium btn-neutral">
				<DoorOpen class="size-5" />
				Dashboard
			</a>
		{:else}
			<a href="/signin" class="btn font-medium btn-outline btn-neutral"> Sign in </a>
			<a href="/signup" class="btn font-medium btn-neutral"> Get started </a>
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
	<div class="border-t border-base-300 bg-base-100 shadow-lg lg:hidden">
		<div class="space-y-3 px-4 py-6">
			<a
				href="/faq"
				class="block rounded-lg px-3 py-2 text-base font-medium transition-colors hover:bg-base-200"
				onclick={closeMobileMenu}
			>
				FAQ
			</a>
			<a
				href="/pricing"
				class="block rounded-lg px-3 py-2 text-base font-medium transition-colors hover:bg-base-200"
				onclick={closeMobileMenu}
			>
				Pricing
			</a>
			{#if isLoggedIn}
				<a
					href="/app"
					class="block rounded-lg px-3 py-2 text-base font-medium transition-colors hover:bg-base-200"
					onclick={closeMobileMenu}
				>
					<DoorOpen class="mr-2 inline size-5" />
					Dashboard
				</a>
			{:else}
				<a
					href="/signin"
					class="block rounded-lg px-3 py-2 text-base font-medium transition-colors hover:bg-base-200"
					onclick={closeMobileMenu}
				>
					Sign in
				</a>
				<a
					href="/signup"
					class="hover:bg-neutral-focus block rounded-lg bg-neutral px-3 py-2 text-base font-medium text-neutral-content transition-colors"
					onclick={closeMobileMenu}
				>
					Get started
				</a>
			{/if}
		</div>
	</div>
{/if}
