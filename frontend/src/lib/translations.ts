export const translations = {
	en: {
		// Common
		getStarted: 'Get started',
		learnMore: 'Learn more',

		// Landing Page
		landing: {
			meta: {
				title: 'reece-start - Sveltekit + Go project start template',
				description:
					'reece-start is a Sveltekit + Go project start template that comes out of the box with auth, billing, and more so you can skip the boilerplate and ship your idea from day one'
			},
			hero: {
				title:
					'Pre-packaged with everything you need to build,<br/>reece-start is how you skip the boilerplate and start shipping',
				getStarted: 'Get started',
				watchDemo: 'Watch demo',
				features: {
					authentication: 'Authentication',
					billing: 'Billing',
					stripeConnect: 'Stripe connect',
					analytics: 'Analytics',
					monitoring: 'Monitoring',
					i18n: 'i18n',
					landingAndDashboard: 'Landing & dashboard pages',
					andMore: 'and more'
				}
			},
			features: {
				title: 'Everything you need to build and scale',
				subtitle: 'Focus on your unique value proposition while we handle the infrastructure',
				authentication: {
					title: 'Advanced Authentication',
					orgMemberUser: 'Organization-Member-User model',
					rbac: 'RBAC using scopes',
					googleOauth: 'Google OAuth integration',
					sudoUsers: 'Sudo users & impersonation'
				},
				organizations: {
					title: 'Organization Management',
					invitationLinks: 'Invitation links & emails',
					roleBasedPermissions: 'Role-based permissions',
					orgSettings: 'Organization settings',
					memberManagement: 'Member management'
				},
				email: {
					title: 'Email Infrastructure',
					notificationApi: 'Notification API',
					templateSystem: 'Template system',
					invitationEmails: 'Invitation emails',
					transactionalEmails: 'Transactional emails'
				},
				billing: {
					title: 'Billing & Payments',
					orgBasedBilling: 'Organization-based billing',
					freePaidTiers: 'Free & paid tiers',
					stripeConnect: 'Stripe Connect support',
					subscriptionManagement: 'Subscription management'
				},
				userExperience: {
					title: 'User Experience',
					userSettings: 'User settings page',
					orgSettings: 'Organization settings',
					landingPages: 'Landing page templates',
					responsiveDesign: 'Responsive design'
				},
				analytics: {
					title: 'Analytics & Monitoring',
					posthog: 'PostHog integration',
					sentry: 'Sentry error tracking',
					performanceMonitoring: 'Performance monitoring',
					userAnalytics: 'User analytics'
				}
			},
			techStack: {
				title: 'Built with modern, production-ready technologies',
				subtitle: 'A carefully curated tech stack that scales with your business',
				sveltekit: {
					title: 'SvelteKit',
					description: 'Modern frontend framework'
				},
				go: {
					title: 'Go + Echo',
					description: 'High-performance backend'
				},
				postgresql: {
					title: 'PostgreSQL',
					description: 'Reliable database'
				},
				docker: {
					title: 'Docker + Railway',
					description: 'Containerized deployment'
				}
			},
			cta: {
				title: 'Ready to skip months of development?',
				subtitle:
					'Start building your SaaS with a solid foundation. No more reinventing the wheel.',
				startBuilding: 'Start building now',
				viewGithub: 'View on GitHub'
			}
		},

		// FAQ Page
		faq: {
			meta: {
				title: 'FAQ - Frequently Asked Questions',
				description:
					'Get answers to frequently asked questions about the tech stack, architecture, and features.'
			},
			header: {
				title: 'Frequently Asked Questions',
				subtitle:
					'Everything you need to know about the technology choices, architecture decisions, and how to get started.'
			},
			questions: {
				whySvelte: {
					question: 'Why Svelte over React?',
					answer:
						'React has become increasingly complicated with SSR, RSC (React Server Components), and complex data loading patterns. Svelte offers a much nicer developer experience for building reactive applications while providing a good middle ground between highly interactive client-side apps and excellent server-side rendering capabilities. The compilation approach eliminates runtime overhead and the syntax is more intuitive and less verbose than React.'
				},
				whyDaisyui: {
					question: 'Why DaisyUI for styling?',
					answer:
						'DaisyUI provides a simple, semantic approach to building components while allowing you to build everything from scratch when needed. It gives you beautiful defaults and consistent design tokens without the complexity of component libraries like Material-UI or Chakra. You get the flexibility of Tailwind CSS with the convenience of pre-built component classes.'
				},
				whyGo: {
					question: 'Why Go for the backend?',
					answer:
						'Go provides exceptional performance with low resource consumption, making it cost-effective for deployment. It has excellent concurrency support, a robust standard library, and compiles to single binaries that are easy to deploy. The strong typing system and explicit error handling lead to more reliable code, and the ecosystem around web development (Echo, GORM) is mature and well-maintained.'
				},
				whyPostgresql: {
					question: 'Why PostgreSQL?',
					answer:
						"PostgreSQL is a battle-tested, feature-rich relational database that handles complex queries, has excellent JSON support for flexibility, provides strong ACID guarantees, and scales well. It's open source, has a massive ecosystem of extensions, and is supported everywhere. Simply put - why not PostgreSQL?"
				},
				whyRailway: {
					question: 'Why Railway for deployment?',
					answer:
						"Railway provides a simple way to deploy containerized applications without the excessive costs of some other platforms. While it's not perfect, it supports most common use cases and offers a straightforward deployment experience with good Docker support, automatic deployments from Git, and reasonable pricing for small to medium-scale applications."
				},
				whyCustomAuth: {
					question: 'Why custom authentication instead of Auth0 or Clerk?',
					answer:
						'While services like Auth0 and Clerk are convenient, they often limit functionality and can become expensive at scale. I frequently found myself fighting these services to implement specific use cases, particularly the organization-member-user model with role-based access control. This project includes a complex, real-world authentication pattern out of the box that would be challenging and costly to implement with third-party auth services.'
				},
				whatIsOrgModel: {
					question: 'What is the organization-member-user model?',
					answer:
						"This is a multi-tenant architecture where users can belong to multiple organizations with different roles in each. Each organization acts as a separate workspace or tenant, and users have specific permissions within each organization they're part of. This pattern is common in B2B SaaS applications and enables complex permission structures and billing models."
				},
				whyMakefile: {
					question: 'Why use a Makefile for development?',
					answer:
						"A Makefile simplifies the monorepo workflow without adding complexity from specialized monorepo tools. It provides simple, cross-platform commands for common development tasks like starting services, running migrations, and building containers. It's a lightweight solution that most developers are familiar with and doesn't require additional tooling or configuration."
				},
				howToStart: {
					question: 'How do I get started with development?',
					answer:
						'Start by cloning the repository and running <code class="rounded bg-base-300 px-2 py-1">make dev</code> to start all services with Docker Compose. The Makefile includes commands for database migrations, seeding data, and running both frontend and backend in development mode. Check the README for detailed setup instructions and environment variable configuration.'
				},
				withoutOrgs: {
					question: 'Can I use this without the organization features?',
					answer:
						"Yes! If you don't need multi-tenant functionality, you can remove the organization-related code and simplify to a standard user authentication model. Remove the organization models, middleware, and related API endpoints, then update the authentication flow to work directly with users instead of organization memberships."
				},
				customizeEmails: {
					question: 'How do I customize the email templates?',
					answer:
						'Email templates are located in <code class="rounded bg-base-300 px-2 py-1">backend/internal/email/templates/</code>. You can modify the HTML templates and update the email service configuration to match your branding. The system supports both HTML and plain text emails with template variable substitution.'
				},
				addOauth: {
					question: 'How do I add new OAuth providers?',
					answer:
						'The project includes Google OAuth as an example. To add new providers, update the OAuth configuration in the backend, add the provider-specific endpoints, and create corresponding frontend sign-in buttons. The authentication system is designed to support multiple OAuth providers simultaneously.'
				},
				otherPlatforms: {
					question: 'Can I deploy this to other platforms besides Railway?',
					answer:
						'Absolutely! The application is containerized with Docker, so it can be deployed to any platform that supports containers: AWS ECS, Google Cloud Run, Azure Container Instances, DigitalOcean App Platform, or even your own servers with Docker Compose. Update the environment variables and database connection strings for your chosen platform.'
				},
				envVariables: {
					question: 'What environment variables do I need to configure?',
					answer:
						'Key environment variables include database connection strings, JWT secrets, OAuth client credentials, email service configuration, and any third-party API keys. Check the <code class="rounded bg-base-300 px-2 py-1">.env.example</code> files in both frontend and backend directories for a complete list of required and optional variables.'
				},
				billingIncluded: {
					question: 'Is billing/payments functionality included?',
					answer:
						'Yes! Stripe is used for payments and is integrated with the organization-based billing model.'
				},
				productionReady: {
					question: 'Is this production-ready?',
					answer:
						'Yes, this starter includes production-ready patterns: proper error handling, logging, database migrations, security middleware, rate limiting foundations, and containerized deployment. However, you should still review and adapt the code for your specific use case, add monitoring, set up proper CI/CD, and perform security audits before launching.'
				}
			}
		},

		// Pricing Page
		pricing: {
			meta: {
				title: 'Pricing - reece-start',
				description:
					'Choose the perfect plan for your business. Start free or upgrade to our Pro plan for advanced features.'
			},
			hero: {
				title: 'Simple, transparent pricing',
				subtitle: 'Start for free and scale as you grow. No hidden fees, no surprises.'
			},
			free: {
				title: 'Free',
				price: '$0',
				period: '/month',
				description: 'Everything you need to get started',
				features: {
					auth: 'Organization-Member-User authentication',
					invitations: 'Organization invitations & emails',
					rbac: 'Role-based access control (RBAC)',
					oauth: 'Google OAuth integration',
					settings: 'User & organization settings',
					emailApi: 'Email notification API',
					sudo: 'Sudo users & impersonation'
				},
				cta: 'Get started free'
			},
			pro: {
				title: 'Pro',
				badge: 'Most Popular',
				price: '$20',
				period: '/month',
				description: 'For people who like giving us money',
				features: {
					everything: 'Everything in Free (which is everything)',
					pride: 'A sense of pride and accomplishment',
					warm: 'The warm fuzzy feeling of supporting developers',
					exclusive: 'Exclusive access to... the same features',
					badge: 'A "Pro" badge that literally does nothing',
					priority: "Priority support for features that don't exist",
					advanced: 'Advanced nothing, but with extra steps',
					gratitude: 'Our eternal gratitude (worth $20/month, apparently)'
				},
				cta: 'Start Pro trial'
			},
			faq: {
				title: 'Frequently Asked Questions',
				changePlans: {
					question: 'Can I change plans at any time?',
					answer:
						"Yes! You can upgrade or downgrade your plan at any time. Changes take effect immediately and we'll prorate any billing adjustments."
				},
				freeTrial: {
					question: 'Is there a free trial for the Pro plan?',
					answer:
						'Yes! We offer a 14-day free trial for the Pro plan. No credit card required to start.'
				},
				paymentMethods: {
					question: 'What payment methods do you accept?',
					answer:
						'We accept all major credit cards (Visa, MasterCard, American Express) and PayPal. All payments are processed securely through Stripe.'
				},
				cancelSubscription: {
					question: 'Can I cancel my subscription at any time?',
					answer:
						'Absolutely. You can cancel your subscription at any time from your account settings. Your access will continue until the end of your current billing period.'
				},
				annualDiscount: {
					question: 'Do you offer discounts for annual billing?',
					answer:
						'Yes! Save 20% when you choose annual billing. Contact our sales team for custom enterprise pricing options.'
				}
			},
			cta: {
				title: 'Ready to get started?',
				subtitle: 'Join the developers already using reece-start to build their SaaS applications.',
				startFree: 'Start free today',
				learnMore: 'Learn more'
			}
		}
	}
} as const;
