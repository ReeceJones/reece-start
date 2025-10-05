<script>
	import { CircleCheck, CreditCard, DoorOpen, Star, Zap } from 'lucide-svelte';
</script>

<svelte:head>
	<title>Pricing - reece-start</title>
	<meta
		name="description"
		content="Choose the perfect plan for your business. Start free or upgrade to our Pro plan for advanced features."
	/>
</svelte:head>

<!-- Hero Section -->
<div class="container mx-auto px-4">
	<div class="mx-auto max-w-4xl py-20 text-center">
		<h1 class="mb-6 text-3xl font-bold md:text-4xl">Simple, transparent pricing</h1>
		<p class="mx-auto max-w-2xl text-xl leading-relaxed text-base-content/70">
			Start for free and scale as you grow. No hidden fees, no surprises.
		</p>
	</div>
</div>

<!-- Pricing Cards Section -->
<section class="bg-base-200 py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto grid max-w-5xl grid-cols-1 gap-8 lg:grid-cols-2">
			<!-- Free Tier -->
			<div class="card border-2 border-transparent bg-base-100 shadow-xl">
				<div class="card-body p-8">
					<div class="mb-6">
						<div class="mb-4 flex items-center gap-3">
							<div class="rounded-lg bg-primary/10 p-2">
								<Zap class="size-6 text-primary" />
							</div>
							<h3 class="text-2xl font-bold">Free</h3>
						</div>
						<div class="mb-4 flex items-baseline gap-2">
							<span class="text-4xl font-bold">$0</span>
							<span class="text-base-content/70">/month</span>
						</div>
						<p class="text-base-content/70">Everything you need to get started</p>
					</div>

					<!-- Features List -->
					<ul class="mb-8 space-y-3">
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Organization-Member-User authentication</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Organization invitations & emails</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Role-based access control (RBAC)</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Google OAuth integration</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>User & organization settings</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Email notification API</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Sudo users & impersonation</span>
						</li>
					</ul>

					<a href="/signup" class="btn w-full font-medium btn-outline btn-lg">
						<DoorOpen class="size-5" />
						Get started free
					</a>
				</div>
			</div>

			<!-- Pro Tier -->
			<div class="card relative border-2 border-accent bg-base-100 shadow-xl">
				<!-- Popular Badge -->
				<div class="absolute -top-4 left-1/2 -translate-x-1/2 transform">
					<div class="rounded-full bg-accent px-4 py-2 text-sm font-medium text-accent-content">
						Most Popular
					</div>
				</div>

				<div class="card-body p-8">
					<div class="mb-6">
						<div class="mb-4 flex items-center gap-3">
							<div class="rounded-lg bg-accent/10 p-2">
								<Star class="size-6 text-accent" />
							</div>
							<h3 class="text-2xl font-bold">Pro</h3>
						</div>
						<div class="mb-4 flex items-baseline gap-2">
							<span class="text-4xl font-bold">$20</span>
							<span class="text-base-content/70">/month</span>
						</div>
						<p class="text-base-content/70">For people who like giving us money</p>
					</div>

					<!-- Features List -->
					<ul class="mb-8 space-y-3">
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Everything in Free (which is everything)</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>A sense of pride and accomplishment</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>The warm fuzzy feeling of supporting developers</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Exclusive access to... the same features</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>A "Pro" badge that literally does nothing</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Priority support for features that don't exist</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Advanced nothing, but with extra steps</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>Our eternal gratitude (worth $20/month, apparently)</span>
						</li>
					</ul>

					<a href="/signup" class="btn w-full font-medium btn-lg btn-accent">
						<CreditCard class="size-5" />
						Start Pro trial
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
			<h2 class="mb-12 text-center text-3xl font-bold">Frequently Asked Questions</h2>

			<div class="join-vertical join w-full">
				<div class="collapse-plus collapse join-item border-b border-base-300">
					<input type="radio" name="pricing-accordion" />
					<div class="collapse-title text-lg font-medium">Can I change plans at any time?</div>
					<div class="collapse-content">
						<p class="text-base-content/70">
							Yes! You can upgrade or downgrade your plan at any time. Changes take effect
							immediately and we'll prorate any billing adjustments.
						</p>
					</div>
				</div>

				<div class="collapse-plus collapse join-item border-b border-base-300">
					<input type="radio" name="pricing-accordion" />
					<div class="collapse-title text-lg font-medium">
						Is there a free trial for the Pro plan?
					</div>
					<div class="collapse-content">
						<p class="text-base-content/70">
							Yes! We offer a 14-day free trial for the Pro plan. No credit card required to start.
						</p>
					</div>
				</div>

				<div class="collapse-plus collapse join-item border-b border-base-300">
					<input type="radio" name="pricing-accordion" />
					<div class="collapse-title text-lg font-medium">What payment methods do you accept?</div>
					<div class="collapse-content">
						<p class="text-base-content/70">
							We accept all major credit cards (Visa, MasterCard, American Express) and PayPal. All
							payments are processed securely through Stripe.
						</p>
					</div>
				</div>

				<div class="collapse-plus collapse join-item border-b border-base-300">
					<input type="radio" name="pricing-accordion" />
					<div class="collapse-title text-lg font-medium">
						Can I cancel my subscription at any time?
					</div>
					<div class="collapse-content">
						<p class="text-base-content/70">
							Absolutely. You can cancel your subscription at any time from your account settings.
							Your access will continue until the end of your current billing period.
						</p>
					</div>
				</div>

				<div class="collapse-plus collapse join-item border-b border-base-300">
					<input type="radio" name="pricing-accordion" />
					<div class="collapse-title text-lg font-medium">
						Do you offer discounts for annual billing?
					</div>
					<div class="collapse-content">
						<p class="text-base-content/70">
							Yes! Save 20% when you choose annual billing. Contact our sales team for custom
							enterprise pricing options.
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
			<h2 class="mb-6 text-3xl font-bold text-neutral-content">Ready to get started?</h2>
			<p class="mb-8 text-lg leading-relaxed text-neutral-content/90">
				Join the developers already using reece-start to build their SaaS applications.
			</p>
			<div class="flex flex-col gap-4 sm:flex-row sm:justify-center">
				<a href="/signup" class="btn font-medium btn-lg btn-accent">
					<DoorOpen class="size-5" />
					Start free today
				</a>
				<a
					href="/faq"
					class="btn text-neutral-content outline-neutral-content btn-outline btn-lg btn-neutral"
				>
					Learn more
				</a>
			</div>
		</div>
	</div>
</section>
