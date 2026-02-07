<script>
	import { CircleCheck, Code, BookOpenIcon } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
	import { buttonVariants } from '$lib/components/ui/button';
</script>

<svelte:head>
	<title>{m.landing__meta__title()}</title>
	<meta name="description" content={m.landing__meta__description()} />
</svelte:head>

<!-- Hero Section -->
<div class="container mx-auto px-4">
	<div class="min-h-screen max-w-2xl space-y-10 pt-24">
		<div class="space-y-6">
			<h1 class="text-lg font-light tracking-normal md:text-2xl">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html m.landing__hero__title()}
			</h1>
			<div class="flex flex-wrap gap-4">
				<a
					href="https://github.com/ReeceJones/reece-start"
					class={buttonVariants({ variant: 'default' })}
				>
					<Code class="size-4" />
					{m.landing__hero__get_code()}
				</a>
				<a href="/faq" class={buttonVariants({ variant: 'outline' })}>
					<BookOpenIcon class="size-4" />
					{m.landing__hero__faq()}
				</a>
			</div>
		</div>

		<!-- Simple feature highlights -->
		<div class="text-base-content/70 text-sm">
			<div class="flex flex-col gap-x-8 gap-y-3">
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{m.landing__hero__features__authentication()}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{m.landing__hero__features__billing()}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{m.landing__hero__features__stripe_connect()}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{m.landing__hero__features__analytics()}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{m.landing__hero__features__monitoring()}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{m.landing__hero__features__i18n()}</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>{m.landing__hero__features__landing_and_dashboard()}</span>
				</div>
				<div class="flex items-center gap-2">
					&hellip;
					<span class="text-primary/70">{m.landing__hero__features__and_more()}</span>
				</div>
			</div>
		</div>
	</div>
</div>
