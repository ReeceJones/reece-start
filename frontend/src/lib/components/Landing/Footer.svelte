<script>
	import { t } from '$lib/i18n';
</script>

<footer class="footer bg-base-200 p-16 text-base-content sm:footer-horizontal">
	<aside>
		<p>
			<b>reece-start</b> - {$t('footer.description')}
			<br />
			{$t('footer.copyright')}
		</p>
	</aside>
	<nav>
		<div class="grid grid-flow-col gap-4">
			<a href="/pricing" class="link link-hover">{$t('footer.pricing')}</a>
			<a href="/faq" class="link link-hover">{$t('footer.faq')}</a>
			<a href="https://github.com/reecejones/reece-start" class="link link-hover"
				>{$t('footer.github')}</a
			>
		</div>
	</nav>
</footer>
