<script>
	import { CircleCheck, Code, BookOpenIcon } from 'lucide-svelte';
	import { t } from '$lib/i18n';
	import { buttonVariants } from '$lib/components/ui/button';
</script>

<svelte:head>
	<title>{$t('landing.meta.title')}</title>
	<meta name="description" content={$t('landing.meta.description')} />
</svelte:head>

<!-- Hero Section -->
<div class="container mx-auto px-4">
	<div class="min-h-screen max-w-2xl space-y-10 pt-24">
		<div class="space-y-6">
			<h1 class="text-lg font-light tracking-normal md:text-2xl">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html $t('landing.hero.title')}
			</h1>
			<div class="flex flex-wrap gap-4">
				<a
					href="https://github.com/ReeceJones/reece-start"
					class={buttonVariants({ variant: 'default' })}
				>
					<Code class="size-4" />
					{$t('landing.hero.getCode')}
				</a>
				<a href="/faq" class={buttonVariants({ variant: 'outline' })}>
					<BookOpenIcon class="size-4" />
					{$t('landing.hero.faq')}
				</a>
			</div>
		</div>

		<!-- Simple feature highlights -->
		<div class="text-base-content/70 text-sm">
			<div class="flex flex-col gap-x-8 gap-y-3">
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{$t('landing.hero.features.authentication')}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{$t('landing.hero.features.billing')}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{$t('landing.hero.features.stripeConnect')}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{$t('landing.hero.features.analytics')}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{$t('landing.hero.features.monitoring')}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{$t('landing.hero.features.i18n')}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{$t('landing.hero.features.landingAndDashboard')}</span>
				</div>
				<div class="flex items-center gap-2">
					&hellip;
					<span class="text-primary/70">{$t('landing.hero.features.andMore')}</span>
				</div>
			</div>
		</div>
	</div>
</div>
