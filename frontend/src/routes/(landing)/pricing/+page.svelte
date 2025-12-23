<script>
	import { CircleCheck, CreditCard, DoorOpen, Star, Zap } from 'lucide-svelte';
	import { t } from '$lib/i18n';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Accordion from '$lib/components/ui/accordion';
</script>

<svelte:head>
	<title>{$t('pricing.meta.title')}</title>
	<meta name="description" content={$t('pricing.meta.description')} />
</svelte:head>

<!-- Hero Section -->
<div class="container mx-auto px-4">
	<div class="mx-auto max-w-4xl py-20 text-center">
		<h1 class="mb-6 text-3xl font-bold md:text-4xl">{$t('pricing.hero.title')}</h1>
		<p class="text-base-content/70 mx-auto max-w-2xl text-xl leading-relaxed">
			{$t('pricing.hero.subtitle')}
		</p>
	</div>
</div>

<!-- Pricing Cards Section -->
<section class="py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto grid max-w-5xl grid-cols-1 gap-8 lg:grid-cols-2">
			<!-- Free Tier -->
			<Card.Root class="shadow-none">
				<Card.Content class="p-8">
					<div class="mb-6">
						<div class="mb-4 flex items-center gap-3">
							<div class="rounded-lg bg-primary/10 p-2">
								<Zap class="size-6 text-primary" />
							</div>
							<h3 class="text-2xl font-bold">{$t('pricing.free.title')}</h3>
						</div>
						<div class="mb-4 flex items-baseline gap-2">
							<span class="text-4xl font-bold">{$t('pricing.free.price')}</span>
							<span class="text-base-content/70">{$t('pricing.free.period')}</span>
						</div>
						<p class="text-base-content/70">{$t('pricing.free.description')}</p>
					</div>

					<!-- Features List -->
					<ul class="mb-8 space-y-3">
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.free.features.auth')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.free.features.invitations')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.free.features.rbac')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.free.features.oauth')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.free.features.settings')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.free.features.emailApi')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.free.features.sudo')}</span>
						</li>
					</ul>

					<a
						href="/signup"
						class={buttonVariants({
							variant: 'secondary',
							size: 'lg',
							class: 'w-full text-base'
						})}
					>
						<DoorOpen class="size-5" />
						{$t('pricing.free.cta')}
					</a>
				</Card.Content>
			</Card.Root>

			<!-- Pro Tier -->
			<Card.Root class="relative border-2 border-accent">
				<!-- Popular Badge -->
				<div class="absolute -top-4 left-1/2 -translate-x-1/2 transform">
					<div class="rounded-full bg-accent px-4 py-2 text-sm font-medium text-accent-foreground">
						{$t('pricing.pro.badge')}
					</div>
				</div>

				<Card.Content class="p-8">
					<div class="mb-6">
						<div class="mb-4 flex items-center gap-3">
							<div class="rounded-lg bg-accent/10 p-2">
								<Star class="size-6 text-accent" />
							</div>
							<h3 class="text-2xl font-bold">{$t('pricing.pro.title')}</h3>
						</div>
						<div class="mb-4 flex items-baseline gap-2">
							<span class="text-4xl font-bold">{$t('pricing.pro.price')}</span>
							<span class="text-base-content/70">{$t('pricing.pro.period')}</span>
						</div>
						<p class="text-base-content/70">{$t('pricing.pro.description')}</p>
					</div>

					<!-- Features List -->
					<ul class="mb-8 space-y-3">
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.pro.features.everything')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.pro.features.pride')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.pro.features.warm')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.pro.features.exclusive')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.pro.features.badge')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.pro.features.priority')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.pro.features.advanced')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{$t('pricing.pro.features.gratitude')}</span>
						</li>
					</ul>

					<a
						href="/signup"
						class={buttonVariants({
							variant: 'accent',
							size: 'lg',
							class: 'w-full text-base'
						})}
					>
						<CreditCard class="size-5" />
						{$t('pricing.pro.cta')}
					</a>
				</Card.Content>
			</Card.Root>
		</div>
	</div>
</section>

<!-- FAQ Section -->
<section class="bg-neutral-100 py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-3xl">
			<h2 class="mb-12 text-center text-3xl font-bold">{$t('pricing.faq.title')}</h2>

			<Accordion.Root type="single">
				<Accordion.Item value="changePlans">
					<Accordion.Trigger class="text-lg font-medium">
						{$t('pricing.faq.changePlans.question')}
					</Accordion.Trigger>
					<Accordion.Content>
						<p class="text-base-content/70">
							{$t('pricing.faq.changePlans.answer')}
						</p>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="freeTrial">
					<Accordion.Trigger class="text-lg font-medium">
						{$t('pricing.faq.freeTrial.question')}
					</Accordion.Trigger>
					<Accordion.Content>
						<p class="text-base-content/70">
							{$t('pricing.faq.freeTrial.answer')}
						</p>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="paymentMethods">
					<Accordion.Trigger class="text-lg font-medium">
						{$t('pricing.faq.paymentMethods.question')}
					</Accordion.Trigger>
					<Accordion.Content>
						<p class="text-base-content/70">
							{$t('pricing.faq.paymentMethods.answer')}
						</p>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="cancelSubscription">
					<Accordion.Trigger class="text-lg font-medium">
						{$t('pricing.faq.cancelSubscription.question')}
					</Accordion.Trigger>
					<Accordion.Content>
						<p class="text-base-content/70">
							{$t('pricing.faq.cancelSubscription.answer')}
						</p>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="annualDiscount">
					<Accordion.Trigger class="text-lg font-medium">
						{$t('pricing.faq.annualDiscount.question')}
					</Accordion.Trigger>
					<Accordion.Content>
						<p class="text-base-content/70">
							{$t('pricing.faq.annualDiscount.answer')}
						</p>
					</Accordion.Content>
				</Accordion.Item>
			</Accordion.Root>
		</div>
	</div>
</section>

<!-- CTA Section -->
<section class="py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-2xl text-center">
			<h2 class="text-neutral-content mb-6 text-3xl font-bold">{$t('pricing.cta.title')}</h2>
			<p class="text-neutral-content/90 mb-8 text-lg leading-relaxed">
				{$t('pricing.cta.subtitle')}
			</p>
			<div class="flex flex-col gap-4 sm:flex-row sm:justify-center">
				<a
					href="/signup"
					class={buttonVariants({
						variant: 'accent',
						size: 'lg',
						class: 'font-medium'
					})}
				>
					<DoorOpen class="size-5" />
					{$t('pricing.cta.startFree')}
				</a>
				<a
					href="/faq"
					class={buttonVariants({
						variant: 'outline',
						size: 'lg',
						class: 'text-neutral-content outline-neutral-content'
					})}
				>
					{$t('pricing.cta.learnMore')}
				</a>
			</div>
		</div>
	</div>
</section>
