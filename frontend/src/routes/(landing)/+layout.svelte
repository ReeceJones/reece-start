<script>
	import Footer from '$lib/components/Landing/Footer.svelte';
	import Navbar from '$lib/components/Landing/Navbar.svelte';

	const { data, children } = $props();

	const isLoggedIn = $derived(data.isLoggedIn);
</script>

<div class="min-h-screen">
	<div class="container mx-auto">
		<Navbar {isLoggedIn} />
	</div>
	<main>{@render children?.()}</main>
	<Footer />
</div>
