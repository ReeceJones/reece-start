<script lang="ts">
	import { CircleQuestionMark } from 'lucide-svelte';

	function handleAccordionClick(event: MouseEvent) {
		// Find the radio input within the clicked accordion
		const radioInput = (event.currentTarget as HTMLElement).querySelector(
			'input[type="radio"]'
		) as HTMLInputElement;
		if (radioInput) {
			radioInput.checked = true;
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		// Handle Enter and Space key presses for accessibility
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			handleAccordionClick(event as any);
		}
	}
</script>

<svelte:head>
	<title>FAQ - Frequently Asked Questions</title>
	<meta
		name="description"
		content="Get answers to frequently asked questions about the tech stack, architecture, and features."
	/>
</svelte:head>

<!-- FAQ Header -->
<div class="container mx-auto px-4">
	<div class="mx-auto max-w-4xl py-24">
		<div class="mb-12 text-center">
			<div class="mb-6 flex justify-center">
				<div class="bg-info/10 rounded-xl p-2">
					<CircleQuestionMark class="text-info size-12" />
				</div>
			</div>
			<h1 class="mb-6 text-3xl font-bold md:text-4xl">Frequently Asked Questions</h1>
			<p class="text-base-content/70 mx-auto max-w-2xl text-xl leading-relaxed">
				Everything you need to know about the technology choices, architecture decisions, and how to
				get started.
			</p>
		</div>

		<!-- FAQ Accordion -->
		<div class="space-y-3">
			<!-- Tech Stack Questions -->
			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why Svelte over React section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">Why Svelte over React?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						React has become increasingly complicated with SSR, RSC (React Server Components), and
						complex data loading patterns. Svelte offers a much nicer developer experience for
						building reactive applications while providing a good middle ground between highly
						interactive client-side apps and excellent server-side rendering capabilities. The
						compilation approach eliminates runtime overhead and the syntax is more intuitive and
						less verbose than React.
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why DaisyUI for styling section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">Why DaisyUI for styling?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						DaisyUI provides a simple, semantic approach to building components while allowing you
						to build everything from scratch when needed. It gives you beautiful defaults and
						consistent design tokens without the complexity of component libraries like Material-UI
						or Chakra. You get the flexibility of Tailwind CSS with the convenience of pre-built
						component classes.
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why Go for the backend section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">Why Go for the backend?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						Go provides exceptional performance with low resource consumption, making it
						cost-effective for deployment. It has excellent concurrency support, a robust standard
						library, and compiles to single binaries that are easy to deploy. The strong typing
						system and explicit error handling lead to more reliable code, and the ecosystem around
						web development (Echo, GORM) is mature and well-maintained.
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why PostgreSQL section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">Why PostgreSQL?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						PostgreSQL is a battle-tested, feature-rich relational database that handles complex
						queries, has excellent JSON support for flexibility, provides strong ACID guarantees,
						and scales well. It's open source, has a massive ecosystem of extensions, and is
						supported everywhere. Simply put - why not PostgreSQL?
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why Railway for deployment section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">Why Railway for deployment?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						Railway provides a simple way to deploy containerized applications without the excessive
						costs of some other platforms. While it's not perfect, it supports most common use cases
						and offers a straightforward deployment experience with good Docker support, automatic
						deployments from Git, and reasonable pricing for small to medium-scale applications.
					</p>
				</div>
			</div>

			<!-- Architecture Questions -->
			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why custom authentication instead of Auth0 or Clerk section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">
					Why custom authentication instead of Auth0 or Clerk?
				</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						While services like Auth0 and Clerk are convenient, they often limit functionality and
						can become expensive at scale. I frequently found myself fighting these services to
						implement specific use cases, particularly the organization-member-user model with
						role-based access control. This project includes a complex, real-world authentication
						pattern out of the box that would be challenging and costly to implement with
						third-party auth services.
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle What is the organization-member-user model section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">What is the organization-member-user model?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						This is a multi-tenant architecture where users can belong to multiple organizations
						with different roles in each. Each organization acts as a separate workspace or tenant,
						and users have specific permissions within each organization they're part of. This
						pattern is common in B2B SaaS applications and enables complex permission structures and
						billing models.
					</p>
				</div>
			</div>

			<!-- Development Questions -->
			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Why use a Makefile for development section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">Why use a Makefile for development?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						A Makefile simplifies the monorepo workflow without adding complexity from specialized
						monorepo tools. It provides simple, cross-platform commands for common development tasks
						like starting services, running migrations, and building containers. It's a lightweight
						solution that most developers are familiar with and doesn't require additional tooling
						or configuration.
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle How do I get started with development section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">How do I get started with development?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						Start by cloning the repository and running <code class="bg-base-300 rounded px-2 py-1"
							>make dev</code
						> to start all services with Docker Compose. The Makefile includes commands for database
						migrations, seeding data, and running both frontend and backend in development mode. Check
						the README for detailed setup instructions and environment variable configuration.
					</p>
				</div>
			</div>

			<!-- Customization Questions -->
			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Can I use this without the organization features section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">
					Can I use this without the organization features?
				</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						Yes! If you don't need multi-tenant functionality, you can remove the
						organization-related code and simplify to a standard user authentication model. Remove
						the organization models, middleware, and related API endpoints, then update the
						authentication flow to work directly with users instead of organization memberships.
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle How do I customize the email templates section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">How do I customize the email templates?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						Email templates are located in <code class="bg-base-300 rounded px-2 py-1"
							>backend/internal/email/templates/</code
						>. You can modify the HTML templates and update the email service configuration to match
						your branding. The system supports both HTML and plain text emails with template
						variable substitution.
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle How do I add new OAuth providers section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">How do I add new OAuth providers?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						The project includes Google OAuth as an example. To add new providers, update the OAuth
						configuration in the backend, add the provider-specific endpoints, and create
						corresponding frontend sign-in buttons. The authentication system is designed to support
						multiple OAuth providers simultaneously.
					</p>
				</div>
			</div>

			<!-- Deployment Questions -->
			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Can I deploy this to other platforms besides Railway section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">
					Can I deploy this to other platforms besides Railway?
				</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						Absolutely! The application is containerized with Docker, so it can be deployed to any
						platform that supports containers: AWS ECS, Google Cloud Run, Azure Container Instances,
						DigitalOcean App Platform, or even your own servers with Docker Compose. Update the
						environment variables and database connection strings for your chosen platform.
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle What environment variables do I need to configure section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">
					What environment variables do I need to configure?
				</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						Key environment variables include database connection strings, JWT secrets, OAuth client
						credentials, email service configuration, and any third-party API keys. Check the <code
							class="bg-base-300 rounded px-2 py-1">.env.example</code
						> files in both frontend and backend directories for a complete list of required and optional
						variables.
					</p>
				</div>
			</div>

			<!-- Billing & Business Questions -->
			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Is billing/payments functionality included section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">Is billing/payments functionality included?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						Yes! Stripe is used for payments and is integrated with the organization-based billing
						model.
					</p>
				</div>
			</div>

			<div
				class="collapse-arrow bg-base-200 border-base-300 collapse border"
				onclick={handleAccordionClick}
				onkeydown={handleKeydown}
				role="button"
				tabindex="0"
				aria-label="Toggle Is this production-ready section"
			>
				<input type="radio" name="faq-accordion" />
				<div class="collapse-title font-semibold">Is this production-ready?</div>
				<div class="collapse-content text-sm">
					<p class="leading-relaxed">
						Yes, this starter includes production-ready patterns: proper error handling, logging,
						database migrations, security middleware, rate limiting foundations, and containerized
						deployment. However, you should still review and adapt the code for your specific use
						case, add monitoring, set up proper CI/CD, and perform security audits before launching.
					</p>
				</div>
			</div>
		</div>
	</div>
</div>
