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
	import { t } from '$lib/i18n';
</script>

<svelte:head>
	<title>{$t('landing.meta.title')}</title>
	<meta name="description" content={$t('landing.meta.description')} />
</svelte:head>

<!-- Hero Section -->
<div class="container mx-auto px-4">
	<div class="min-h-screen max-w-2xl space-y-8 pt-24">
		<div class="space-y-6">
			<h1 class="text-lg font-light tracking-normal md:text-2xl">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html $t('landing.hero.title')}
			</h1>
			<div class="flex gap-4">
				<a href="/signup" class="btn font-medium btn-neutral">
					<DoorOpen class="size-4" />
					{$t('landing.hero.getStarted')}
				</a>
				<a href="#features" class="btn font-medium btn-outline btn-neutral">
					<Play class="size-4" />
					{$t('landing.hero.watchDemo')}
				</a>
			</div>
		</div>

		<!-- Simple feature highlights -->
		<div class="pt-8 text-sm text-base-content/70">
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
					<span class="text-base-content/70">{$t('landing.hero.features.andMore')}</span>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Features Section -->
<section id="features" class="bg-base-200 py-48">
	<div class="container mx-auto px-4">
		<div class="mb-20">
			<h2 class="mb-6 text-4xl font-bold">{$t('landing.features.title')}</h2>
			<p class="max-w-2xl text-xl leading-relaxed text-base-content/70">
				{$t('landing.features.subtitle')}
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
						<h3 class="card-title">{$t('landing.features.authentication.title')}</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.authentication.orgMemberUser')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.authentication.rbac')}
						</li>
						<li class="flex items-center gap-2">
							<Star class="size-4 text-warning" />
							{$t('landing.features.authentication.googleOauth')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.authentication.sudoUsers')}
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
						<h3 class="card-title">{$t('landing.features.organizations.title')}</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.organizations.invitationLinks')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.organizations.roleBasedPermissions')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.organizations.orgSettings')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.organizations.memberManagement')}
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
						<h3 class="card-title">{$t('landing.features.email.title')}</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.email.notificationApi')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.email.templateSystem')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.email.invitationEmails')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.email.transactionalEmails')}
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
						<h3 class="card-title">{$t('landing.features.billing.title')}</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.billing.orgBasedBilling')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.billing.freePaidTiers')}
						</li>
						<li class="flex items-center gap-2">
							<Star class="size-4 text-warning" />
							{$t('landing.features.billing.stripeConnect')}
						</li>
						<li class="flex items-center gap-2">
							<Star class="size-4 text-warning" />
							{$t('landing.features.billing.subscriptionManagement')}
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
						<h3 class="card-title">{$t('landing.features.userExperience.title')}</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.userExperience.userSettings')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.userExperience.orgSettings')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.userExperience.landingPages')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.userExperience.responsiveDesign')}
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
						<h3 class="card-title">{$t('landing.features.analytics.title')}</h3>
					</div>
					<ul class="space-y-2 text-base-content/80">
						<li class="flex items-center gap-2">
							<Star class="size-4 text-warning" />
							{$t('landing.features.analytics.posthog')}
						</li>
						<li class="flex items-center gap-2">
							<Star class="size-4 text-warning" />
							{$t('landing.features.analytics.sentry')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.analytics.performanceMonitoring')}
						</li>
						<li class="flex items-center gap-2">
							<CircleCheck class="size-4 text-success" />
							{$t('landing.features.analytics.userAnalytics')}
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
			<h2 class="mb-6 text-4xl font-bold">{$t('landing.techStack.title')}</h2>
			<p class="max-w-2xl text-xl leading-relaxed text-base-content/70">
				{$t('landing.techStack.subtitle')}
			</p>
		</div>

		<div class="grid grid-cols-2 gap-12 md:grid-cols-4">
			<div class="text-center">
				<div class="mb-4 rounded-xl bg-base-200 p-6">
					<Globe class="mx-auto size-12 text-primary" />
				</div>
				<h3 class="font-semibold">{$t('landing.techStack.sveltekit.title')}</h3>
				<p class="text-sm text-base-content/70">{$t('landing.techStack.sveltekit.description')}</p>
			</div>
			<div class="text-center">
				<div class="mb-4 rounded-xl bg-base-200 p-6">
					<Zap class="mx-auto size-12 text-warning" />
				</div>
				<h3 class="font-semibold">{$t('landing.techStack.go.title')}</h3>
				<p class="text-sm text-base-content/70">{$t('landing.techStack.go.description')}</p>
			</div>
			<div class="text-center">
				<div class="mb-4 rounded-xl bg-base-200 p-6">
					<ChartColumn class="mx-auto size-12 text-accent" />
				</div>
				<h3 class="font-semibold">{$t('landing.techStack.postgresql.title')}</h3>
				<p class="text-sm text-base-content/70">
					{$t('landing.techStack.postgresql.description')}
				</p>
			</div>
			<div class="text-center">
				<div class="mb-4 rounded-xl bg-base-200 p-6">
					<Rocket class="mx-auto size-12 text-error" />
				</div>
				<h3 class="font-semibold">{$t('landing.techStack.docker.title')}</h3>
				<p class="text-sm text-base-content/70">{$t('landing.techStack.docker.description')}</p>
			</div>
		</div>
	</div>
</section>

<!-- CTA Section -->
<section class="bg-neutral py-48">
	<div class="container mx-auto px-4">
		<h2 class="mb-8 text-4xl font-bold text-neutral-content">
			{$t('landing.cta.title')}
		</h2>
		<p class="mb-12 max-w-2xl text-xl leading-relaxed text-neutral-content/90">
			{$t('landing.cta.subtitle')}
		</p>
		<div class="flex flex-col gap-4 sm:flex-row">
			<a href="/signup" class="btn font-medium btn-lg btn-accent">
				<DoorOpen class="size-5" />
				{$t('landing.cta.startBuilding')}
			</a>
			<a
				href="https://github.com/reecejones/reece-start"
				class="btn text-neutral-content outline-neutral-content btn-outline btn-lg btn-neutral"
			>
				<Code class="size-5" />
				{$t('landing.cta.viewGithub')}
			</a>
		</div>
	</div>
</section>
