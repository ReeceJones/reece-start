<script>
	import { Rocket } from 'lucide-svelte';
</script>

<footer class="footer bg-base-200 p-16 text-base-content sm:footer-horizontal">
	<aside>
		<Rocket class="size-4 text-primary" />
		<p>
			<b>reece-start</b> - Production-ready SvelteKit + Go starter template for building SaaS
			applications
			<br />
			Copyright © 2024 - All rights reserved
		</p>
	</aside>
	<nav>
		<div class="grid grid-flow-col gap-4">
			<a href="/pricing" class="link link-hover">Pricing</a>
			<a href="/faq" class="link link-hover">FAQ</a>
			<a href="https://github.com/reecejones/reece-start" class="link link-hover">GitHub</a>
		</div>
	</nav>
</footer>
