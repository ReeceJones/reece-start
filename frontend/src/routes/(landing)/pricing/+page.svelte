<script>
	import { CircleCheck, CreditCard, DoorOpen, Star, Zap } from 'lucide-svelte';
	import { t } from '$lib/i18n';
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
<section class="bg-base-200 py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto grid max-w-5xl grid-cols-1 gap-8 lg:grid-cols-2">
			<!-- Free Tier -->
			<div class="card bg-base-100 border-2 border-transparent shadow-xl">
				<div class="card-body p-8">
					<div class="mb-6">
						<div class="mb-4 flex items-center gap-3">
							<div class="bg-primary/10 rounded-lg p-2">
								<Zap class="text-primary size-6" />
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
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.free.features.auth')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.free.features.invitations')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.free.features.rbac')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.free.features.oauth')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.free.features.settings')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.free.features.emailApi')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.free.features.sudo')}</span>
						</li>
					</ul>

					<a href="/signup" class="btn btn-outline btn-lg w-full font-medium">
						<DoorOpen class="size-5" />
						{$t('pricing.free.cta')}
					</a>
				</div>
			</div>

			<!-- Pro Tier -->
			<div class="card border-accent bg-base-100 relative border-2 shadow-xl">
				<!-- Popular Badge -->
				<div class="absolute -top-4 left-1/2 -translate-x-1/2 transform">
					<div class="bg-accent text-accent-content rounded-full px-4 py-2 text-sm font-medium">
						{$t('pricing.pro.badge')}
					</div>
				</div>

				<div class="card-body p-8">
					<div class="mb-6">
						<div class="mb-4 flex items-center gap-3">
							<div class="bg-accent/10 rounded-lg p-2">
								<Star class="text-accent size-6" />
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
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.pro.features.everything')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.pro.features.pride')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.pro.features.warm')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.pro.features.exclusive')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.pro.features.badge')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.pro.features.priority')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.pro.features.advanced')}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="text-success size-5 shrink-0" />
							<span>{$t('pricing.pro.features.gratitude')}</span>
						</li>
					</ul>

					<a href="/signup" class="btn btn-lg btn-accent w-full font-medium">
						<CreditCard class="size-5" />
						{$t('pricing.pro.cta')}
					</a>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- FAQ Section -->
<section class="bg-base-100 py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-3xl">
			<h2 class="mb-12 text-center text-3xl font-bold">{$t('pricing.faq.title')}</h2>

			<div class="join-vertical join w-full">
				<div class="collapse-plus join-item border-base-300 collapse border-b">
					<input type="radio" name="pricing-accordion" />
					<div class="collapse-title text-lg font-medium">
						{$t('pricing.faq.changePlans.question')}
					</div>
					<div class="collapse-content">
						<p class="text-base-content/70">
							{$t('pricing.faq.changePlans.answer')}
						</p>
					</div>
				</div>

				<div class="collapse-plus join-item border-base-300 collapse border-b">
					<input type="radio" name="pricing-accordion" />
					<div class="collapse-title text-lg font-medium">
						{$t('pricing.faq.freeTrial.question')}
					</div>
					<div class="collapse-content">
						<p class="text-base-content/70">
							{$t('pricing.faq.freeTrial.answer')}
						</p>
					</div>
				</div>

				<div class="collapse-plus join-item border-base-300 collapse border-b">
					<input type="radio" name="pricing-accordion" />
					<div class="collapse-title text-lg font-medium">
						{$t('pricing.faq.paymentMethods.question')}
					</div>
					<div class="collapse-content">
						<p class="text-base-content/70">
							{$t('pricing.faq.paymentMethods.answer')}
						</p>
					</div>
				</div>

				<div class="collapse-plus join-item border-base-300 collapse border-b">
					<input type="radio" name="pricing-accordion" />
					<div class="collapse-title text-lg font-medium">
						{$t('pricing.faq.cancelSubscription.question')}
					</div>
					<div class="collapse-content">
						<p class="text-base-content/70">
							{$t('pricing.faq.cancelSubscription.answer')}
						</p>
					</div>
				</div>

				<div class="collapse-plus join-item border-base-300 collapse border-b">
					<input type="radio" name="pricing-accordion" />
					<div class="collapse-title text-lg font-medium">
						{$t('pricing.faq.annualDiscount.question')}
					</div>
					<div class="collapse-content">
						<p class="text-base-content/70">
							{$t('pricing.faq.annualDiscount.answer')}
						</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- CTA Section -->
<section class="bg-neutral py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-2xl text-center">
			<h2 class="text-neutral-content mb-6 text-3xl font-bold">{$t('pricing.cta.title')}</h2>
			<p class="text-neutral-content/90 mb-8 text-lg leading-relaxed">
				{$t('pricing.cta.subtitle')}
			</p>
			<div class="flex flex-col gap-4 sm:flex-row sm:justify-center">
				<a href="/signup" class="btn btn-lg btn-accent font-medium">
					<DoorOpen class="size-5" />
					{$t('pricing.cta.startFree')}
				</a>
				<a
					href="/faq"
					class="btn text-neutral-content outline-neutral-content btn-outline btn-lg btn-neutral"
				>
					{$t('pricing.cta.learnMore')}
				</a>
			</div>
		</div>
	</div>
</section>
