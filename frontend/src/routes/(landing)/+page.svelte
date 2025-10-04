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
				<a href="/signup" class="btn btn-neutral font-medium">
					<DoorOpen class="size-4" />
					Get started
				</a>
				<a href="#features" class="btn btn-neutral btn-outline font-medium">
					<Play class="size-4" />
					Watch demo
				</a>
			</div>
		</div>

		<!-- Simple feature highlights -->
		<div class="text-base-content/70 pt-8 text-sm">
			<div class="flex flex-col gap-x-8 gap-y-3">
				<div class="flex items-center gap-2">
					<CircleCheck class="text-success size-4" />
					<span>Authentication</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="text-success size-4" />
					<span>Stripe connect</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="text-success size-4" />
					<span>Analytics & monitoring</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="text-success size-4" />
					<span>i18n</span>
				</div>
				<div class="flex items-center gap-2">
					<CircleCheck class="text-success size-4" />
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
			<p class="text-base-content/70 max-w-2xl text-xl leading-relaxed">
				Focus on your unique value proposition while we handle the infrastructure
			</p>
		</div>

		<div class="grid grid-cols-1 gap-12 md:grid-cols-2 lg:grid-cols-3">
			<!-- Authentication -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="bg-primary/10 rounded-lg p-2">
							<Shield class="text-primary size-6" />
						</div>
						<h3 class="card-title">Advanced Authentication</h3>
					</div>
					<ul class="text-base-content/80 space-y-2">
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Organization-Member-User model
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							RBAC using scopes
						</li>
						<li class="flex items-center gap-2">
							<Star class="text-warning size-4" />
							Google OAuth integration
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Sudo users & impersonation
						</li>
					</ul>
				</div>
			</div>

			<!-- Organizations -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="bg-secondary/10 rounded-lg p-2">
							<Users class="text-secondary size-6" />
						</div>
						<h3 class="card-title">Organization Management</h3>
					</div>
					<ul class="text-base-content/80 space-y-2">
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Invitation links & emails
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Role-based permissions
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Organization settings
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Member management
						</li>
					</ul>
				</div>
			</div>

			<!-- Email System -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="bg-accent/10 rounded-lg p-2">
							<Mail class="text-accent size-6" />
						</div>
						<h3 class="card-title">Email Infrastructure</h3>
					</div>
					<ul class="text-base-content/80 space-y-2">
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Notification API
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Template system
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Invitation emails
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Transactional emails
						</li>
					</ul>
				</div>
			</div>

			<!-- Billing -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="bg-warning/10 rounded-lg p-2">
							<CreditCard class="text-warning size-6" />
						</div>
						<h3 class="card-title">Billing & Payments</h3>
					</div>
					<ul class="text-base-content/80 space-y-2">
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Organization-based billing
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Free & paid tiers
						</li>
						<li class="flex items-center gap-2">
							<Star class="text-warning size-4" />
							Stripe Connect support
						</li>
						<li class="flex items-center gap-2">
							<Star class="text-warning size-4" />
							Subscription management
						</li>
					</ul>
				</div>
			</div>

			<!-- User Experience -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="bg-info/10 rounded-lg p-2">
							<Settings class="text-info size-6" />
						</div>
						<h3 class="card-title">User Experience</h3>
					</div>
					<ul class="text-base-content/80 space-y-2">
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							User settings page
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Organization settings
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Landing page templates
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Responsive design
						</li>
					</ul>
				</div>
			</div>

			<!-- Analytics & Monitoring -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<div class="mb-4 flex items-center gap-3">
						<div class="bg-success/10 rounded-lg p-2">
							<ChartColumn class="text-success size-6" />
						</div>
						<h3 class="card-title">Analytics & Monitoring</h3>
					</div>
					<ul class="text-base-content/80 space-y-2">
						<li class="flex items-center gap-2">
							<Star class="text-warning size-4" />
							PostHog integration
						</li>
						<li class="flex items-center gap-2">
							<Star class="text-warning size-4" />
							Sentry error tracking
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
							Performance monitoring
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="text-success size-4" />
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
			<p class="text-base-content/70 max-w-2xl text-xl leading-relaxed">
				A carefully curated tech stack that scales with your business
			</p>
		</div>

		<div class="grid grid-cols-2 gap-12 md:grid-cols-4">
			<div class="text-center">
				<div class="bg-base-200 mb-4 rounded-xl p-6">
					<Globe class="text-primary mx-auto size-12" />
				</div>
				<h3 class="font-semibold">SvelteKit</h3>
				<p class="text-base-content/70 text-sm">Modern frontend framework</p>
			</div>
			<div class="text-center">
				<div class="bg-base-200 mb-4 rounded-xl p-6">
					<Zap class="text-warning mx-auto size-12" />
				</div>
				<h3 class="font-semibold">Go + Echo</h3>
				<p class="text-base-content/70 text-sm">High-performance backend</p>
			</div>
			<div class="text-center">
				<div class="bg-base-200 mb-4 rounded-xl p-6">
					<ChartColumn class="text-accent mx-auto size-12" />
				</div>
				<h3 class="font-semibold">PostgreSQL</h3>
				<p class="text-base-content/70 text-sm">Reliable database</p>
			</div>
			<div class="text-center">
				<div class="bg-base-200 mb-4 rounded-xl p-6">
					<Rocket class="text-error mx-auto size-12" />
				</div>
				<h3 class="font-semibold">Docker + Railway</h3>
				<p class="text-base-content/70 text-sm">Containerized deployment</p>
			</div>
		</div>
	</div>
</section>

<!-- CTA Section -->
<section class="bg-neutral py-48">
	<div class="container mx-auto px-4">
		<h2 class="text-neutral-content mb-8 text-4xl font-bold">
			Ready to skip months of development?
		</h2>
		<p class="text-neutral-content/90 mb-12 max-w-2xl text-xl leading-relaxed">
			Start building your SaaS with a solid foundation. No more reinventing the wheel.
		</p>
		<div class="flex flex-col gap-4 sm:flex-row">
			<a href="/signup" class="btn btn-accent btn-lg font-medium">
				<DoorOpen class="size-5" />
				Start building now
			</a>
			<a
				href="https://github.com/reecejones/reece-start"
				class="btn btn-outline btn-lg btn-neutral text-neutral-content outline-neutral-content"
			>
				<Code class="size-5" />
				View on GitHub
			</a>
		</div>
	</div>
</section>
