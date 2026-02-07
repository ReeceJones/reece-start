<script>
	import { CircleCheck, CreditCard, DoorOpen, Star, Zap } from 'lucide-svelte';
	import * as m from '$lib/paraglide/messages';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Accordion from '$lib/components/ui/accordion';
</script>

<svelte:head>
	<title>{m.pricing__meta__title()}</title>
	<meta name="description" content={m.pricing__meta__description()} />
</svelte:head>

<!-- Hero Section -->
<div class="container mx-auto px-4">
	<div class="mx-auto max-w-4xl py-20 text-center">
		<h1 class="mb-6 text-3xl font-bold md:text-4xl">{m.pricing__hero__title()}</h1>
		<p class="text-base-content/70 mx-auto max-w-2xl text-xl leading-relaxed">
			{m.pricing__hero__subtitle()}
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
							<h3 class="text-2xl font-bold">{m.pricing__free__title()}</h3>
						</div>
						<div class="mb-4 flex items-baseline gap-2">
							<span class="text-4xl font-bold">{m.pricing__free__price()}</span>
							<span class="text-base-content/70">{m.pricing__free__period()}</span>
						</div>
						<p class="text-base-content/70">{m.pricing__free__description()}</p>
					</div>

					<!-- Features List -->
					<ul class="mb-8 space-y-3">
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__free__features__auth()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__free__features__invitations()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__free__features__rbac()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__free__features__oauth()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__free__features__settings()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__free__features__email_api()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__free__features__sudo()}</span>
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
						{m.pricing__free__cta()}
					</a>
				</Card.Content>
			</Card.Root>

			<!-- Pro Tier -->
			<Card.Root class="relative border-2 border-accent">
				<!-- Popular Badge -->
				<div class="absolute -top-4 left-1/2 -translate-x-1/2 transform">
					<div class="rounded-full bg-accent px-4 py-2 text-sm font-medium text-accent-foreground">
						{m.pricing__pro__badge()}
					</div>
				</div>

				<Card.Content class="p-8">
					<div class="mb-6">
						<div class="mb-4 flex items-center gap-3">
							<div class="rounded-lg bg-accent/10 p-2">
								<Star class="size-6 text-accent" />
							</div>
							<h3 class="text-2xl font-bold">{m.pricing__pro__title()}</h3>
						</div>
						<div class="mb-4 flex items-baseline gap-2">
							<span class="text-4xl font-bold">{m.pricing__pro__price()}</span>
							<span class="text-base-content/70">{m.pricing__pro__period()}</span>
						</div>
						<p class="text-base-content/70">{m.pricing__pro__description()}</p>
					</div>

					<!-- Features List -->
					<ul class="mb-8 space-y-3">
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__pro__features__everything()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__pro__features__pride()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__pro__features__warm()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__pro__features__exclusive()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__pro__features__badge()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__pro__features__priority()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__pro__features__advanced()}</span>
						</li>
						<li class="flex items-center gap-3">
							<CircleCheck class="size-5 shrink-0 text-success" />
							<span>{m.pricing__pro__features__gratitude()}</span>
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
						{m.pricing__pro__cta()}
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
			<h2 class="mb-12 text-center text-3xl font-bold">{m.pricing__faq__title()}</h2>

			<Accordion.Root type="single">
				<Accordion.Item value="changePlans">
					<Accordion.Trigger class="text-lg font-medium">
						{m.pricing__faq__change_plans__question()}
					</Accordion.Trigger>
					<Accordion.Content>
						<p class="text-base-content/70">
							{m.pricing__faq__change_plans__answer()}
						</p>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="freeTrial">
					<Accordion.Trigger class="text-lg font-medium">
						{m.pricing__faq__free_trial__question()}
					</Accordion.Trigger>
					<Accordion.Content>
						<p class="text-base-content/70">
							{m.pricing__faq__free_trial__answer()}
						</p>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="paymentMethods">
					<Accordion.Trigger class="text-lg font-medium">
						{m.pricing__faq__payment_methods__question()}
					</Accordion.Trigger>
					<Accordion.Content>
						<p class="text-base-content/70">
							{m.pricing__faq__payment_methods__answer()}
						</p>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="cancelSubscription">
					<Accordion.Trigger class="text-lg font-medium">
						{m.pricing__faq__cancel_subscription__question()}
					</Accordion.Trigger>
					<Accordion.Content>
						<p class="text-base-content/70">
							{m.pricing__faq__cancel_subscription__answer()}
						</p>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="annualDiscount">
					<Accordion.Trigger class="text-lg font-medium">
						{m.pricing__faq__annual_discount__question()}
					</Accordion.Trigger>
					<Accordion.Content>
						<p class="text-base-content/70">
							{m.pricing__faq__annual_discount__answer()}
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
			<h2 class="text-neutral-content mb-6 text-3xl font-bold">{m.pricing__cta__title()}</h2>
			<p class="text-neutral-content/90 mb-8 text-lg leading-relaxed">
				{m.pricing__cta__subtitle()}
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
					{m.pricing__cta__start_free()}
				</a>
				<a
					href="/faq"
					class={buttonVariants({
						variant: 'outline',
						size: 'lg',
						class: 'text-neutral-content outline-neutral-content'
					})}
				>
					{m.pricing__cta__learn_more()}
				</a>
			</div>
		</div>
	</div>
</section>
