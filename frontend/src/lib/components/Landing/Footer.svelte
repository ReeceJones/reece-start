<script>
	import * as m from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import { buttonVariants } from '../ui/button';
</script>

<footer class="text-base-content sm:footer-horizontal bg-footer p-16 text-footer-foreground">
	<aside>
		<p>
			<b>reece-start</b> - {m.footer__description()}
		</p>
		<p>
			{m.footer__copyright()}
		</p>
	</aside>
	<nav>
		<div class="flex flex-row flex-wrap gap-4">
			<a
				href="/pricing"
				class={cn(
					buttonVariants({ variant: 'link' }),
					'text-primary-foreground hover:text-primary-foreground'
				)}>{m.footer__pricing()}</a
			>
			<a
				href="/faq"
				class={cn(
					buttonVariants({ variant: 'link' }),
					'text-primary-foreground hover:text-primary-foreground'
				)}>{m.footer__faq()}</a
			>
			<a
				href="https://github.com/reecejones/reece-start"
				class={cn(
					buttonVariants({ variant: 'link' }),
					'text-primary-foreground hover:text-primary-foreground'
				)}>{m.footer__github()}</a
			>
		</div>
	</nav>
</footer>
