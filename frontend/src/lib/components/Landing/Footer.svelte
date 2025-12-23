<script>
	import { t } from '$lib/i18n';
	import { cn } from '$lib/utils';
	import { buttonVariants } from '../ui/button';
</script>

<footer class="text-base-content sm:footer-horizontal bg-footer p-16 text-footer-foreground">
	<aside>
		<p>
			<b>reece-start</b> - {$t('footer.description')}
		</p>
		<p>
			{$t('footer.copyright')}
		</p>
	</aside>
	<nav>
		<div class="flex flex-row flex-wrap gap-4">
			<a
				href="/pricing"
				class={cn(
					buttonVariants({ variant: 'link' }),
					'text-primary-foreground hover:text-primary-foreground'
				)}>{$t('footer.pricing')}</a
			>
			<a
				href="/faq"
				class={cn(
					buttonVariants({ variant: 'link' }),
					'text-primary-foreground hover:text-primary-foreground'
				)}>{$t('footer.faq')}</a
			>
			<a
				href="https://github.com/reecejones/reece-start"
				class={cn(
					buttonVariants({ variant: 'link' }),
					'text-primary-foreground hover:text-primary-foreground'
				)}>{$t('footer.github')}</a
			>
		</div>
	</nav>
</footer>
