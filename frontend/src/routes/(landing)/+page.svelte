<script>
	import {
		DoorOpen,
		Play,
		Shield,
		Users,
		Mail,
		Settings,
		CreditCard,
		Zap,
		CircleCheck,
		Star,
		Globe,
		ChartColumn,
		Code,
		Rocket
	} from 'lucide-svelte';
</script>

<svelte:head>
	<title>reece-start - Sveltekit + Go project start template</title>
	<meta
		name="description"
		content="reece-start is a Sveltekit + Go project start template that comes out of the box with auth, billing, and more so you can skip the boilerplate and ship your idea from day one"
	/>
</svelte:head>

<!-- Hero Section -->
<div class="container mx-auto px-4">
	<div class="min-h-screen max-w-2xl space-y-8 pt-24">
		<div class="space-y-6">
			<h1 class="text-lg font-light tracking-normal md:text-2xl">
				Pre-packaged with everything you need to build,
				<br />
				reece-start is how you skip the boilerplate and start shipping
			</h1>
			<div class="flex gap-4">
				<a href="/signup" class="btn font-medium btn-neutral">
					<DoorOpen class="size-4" />
					Get started
				</a>
				<a href="#features" class="btn font-medium btn-outline btn-neutral">
					<Play class="size-4" />
					Watch demo
				</a>
			</div>
		</div>

		<!-- Simple feature highlights -->
		<div class="pt-8 text-sm text-base-content/70">
			<div class="flex flex-col gap-x-8 gap-y-3">
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>Authentication</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>Billing</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>Stripe connect</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>Analytics</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>Monitoring</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>i18n</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="size-4 text-success" />
					<span>Landing & dashboard pages</span>
				</div>
				<div class="flex items-center gap-2">
					&hellip;
					<span class="text-base-content/70">and more</span>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Features Section -->
<section id="features" class="bg-base-200 py-48">
	<div class="container mx-auto px-4">
		<div class="mb-20">
			<h2 class="mb-6 text-4xl font-bold">Everything you need to build and scale</h2>
			<p class="max-w-2xl text-xl leading-relaxed text-base-content/70">
				Focus on your unique value proposition while we handle the infrastructure
			</p>
		</div>

		<div class="grid grid-cols-1 gap-12 md:grid-cols-2 lg:grid-cols-3">
			<!-- Authentication -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="rounded-lg bg-primary/10 p-2">
							<Shield class="size-6 text-primary" />
						</div>
						<h3 class="card-title">Advanced Authentication</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Organization-Member-User model
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							RBAC using scopes
						</li>
						<li class="flex items-center gap-2">
							<Star class="size-4 text-warning" />
							Google OAuth integration
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Sudo users & impersonation
						</li>
					</ul>
				</div>
			</div>

			<!-- Organizations -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="rounded-lg bg-secondary/10 p-2">
							<Users class="size-6 text-secondary" />
						</div>
						<h3 class="card-title">Organization Management</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Invitation links & emails
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Role-based permissions
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Organization settings
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Member management
						</li>
					</ul>
				</div>
			</div>

			<!-- Email System -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="rounded-lg bg-accent/10 p-2">
							<Mail class="size-6 text-accent" />
						</div>
						<h3 class="card-title">Email Infrastructure</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Notification API
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Template system
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Invitation emails
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Transactional emails
						</li>
					</ul>
				</div>
			</div>

			<!-- Billing -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="rounded-lg bg-warning/10 p-2">
							<CreditCard class="size-6 text-warning" />
						</div>
						<h3 class="card-title">Billing & Payments</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Organization-based billing
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Free & paid tiers
						</li>
						<li class="flex items-center gap-2">
							<Star class="size-4 text-warning" />
							Stripe Connect support
						</li>
						<li class="flex items-center gap-2">
							<Star class="size-4 text-warning" />
							Subscription management
						</li>
					</ul>
				</div>
			</div>

			<!-- User Experience -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="rounded-lg bg-info/10 p-2">
							<Settings class="size-6 text-info" />
						</div>
						<h3 class="card-title">User Experience</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							User settings page
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Organization settings
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Landing page templates
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Responsive design
						</li>
					</ul>
				</div>
			</div>

			<!-- Analytics & Monitoring -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="rounded-lg bg-success/10 p-2">
							<ChartColumn class="size-6 text-success" />
						</div>
						<h3 class="card-title">Analytics & Monitoring</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<Star class="size-4 text-warning" />
							PostHog integration
						</li>
						<li class="flex items-center gap-2">
							<Star class="size-4 text-warning" />
							Sentry error tracking
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							Performance monitoring
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							User analytics
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- Tech Stack Section -->
<section class="bg-base-100 py-48">
	<div class="container mx-auto px-4">
		<div class="mb-20">
			<h2 class="mb-6 text-4xl font-bold">Built with modern, production-ready technologies</h2>
			<p class="max-w-2xl text-xl leading-relaxed text-base-content/70">
				A carefully curated tech stack that scales with your business
			</p>
		</div>

		<div class="grid grid-cols-2 gap-12 md:grid-cols-4">
			<div class="text-center">
				<div class="mb-4 rounded-xl bg-base-200 p-6">
					<Globe class="mx-auto size-12 text-primary" />
				</div>
				<h3 class="font-semibold">SvelteKit</h3>
				<p class="text-sm text-base-content/70">Modern frontend framework</p>
			</div>
			<div class="text-center">
				<div class="mb-4 rounded-xl bg-base-200 p-6">
					<Zap class="mx-auto size-12 text-warning" />
				</div>
				<h3 class="font-semibold">Go + Echo</h3>
				<p class="text-sm text-base-content/70">High-performance backend</p>
			</div>
			<div class="text-center">
				<div class="mb-4 rounded-xl bg-base-200 p-6">
					<ChartColumn class="mx-auto size-12 text-accent" />
				</div>
				<h3 class="font-semibold">PostgreSQL</h3>
				<p class="text-sm text-base-content/70">Reliable database</p>
			</div>
			<div class="text-center">
				<div class="mb-4 rounded-xl bg-base-200 p-6">
					<Rocket class="mx-auto size-12 text-error" />
				</div>
				<h3 class="font-semibold">Docker + Railway</h3>
				<p class="text-sm text-base-content/70">Containerized deployment</p>
			</div>
		</div>
	</div>
</section>

<!-- CTA Section -->
<section class="bg-neutral py-48">
	<div class="container mx-auto px-4">
		<h2 class="mb-8 text-4xl font-bold text-neutral-content">
			Ready to skip months of development?
		</h2>
		<p class="mb-12 max-w-2xl text-xl leading-relaxed text-neutral-content/90">
			Start building your SaaS with a solid foundation. No more reinventing the wheel.
		</p>
		<div class="flex flex-col gap-4 sm:flex-row">
			<a href="/signup" class="btn font-medium btn-lg btn-accent">
				<DoorOpen class="size-5" />
				Start building now
			</a>
			<a
				href="https://github.com/reecejones/reece-start"
				class="btn text-neutral-content outline-neutral-content btn-outline btn-lg btn-neutral"
			>
				<Code class="size-5" />
				View on GitHub
			</a>
		</div>
	</div>
</section>
