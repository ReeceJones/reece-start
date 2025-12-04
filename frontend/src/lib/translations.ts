export const translations = {
	en: {
		// Common
		getStarted: 'Get started',
		learnMore: 'Learn more',
		signIn: 'Sign in',
		signUp: 'Sign up',
		dashboard: 'Dashboard',
		home: 'Home',
		logout: 'Logout',
		profileTitle: 'Profile',
		next: 'Next',
		back: 'Back',
		organization: 'Organization',
		switchOrganization: 'Switch Organization',
		foo: 'Foo',
		bar: 'Bar',
		close: 'Close',
		cancel: 'Cancel',
		save: 'Save',
		create: 'Create',
		edit: 'Edit',
		delete: 'Delete',
		copy: 'Copy',
		invite: 'Invite',
		admin: 'Admin',
		member: 'Member',
		users: 'Users',
		debug: 'Debug',
		application: 'Application',
		stopImpersonation: 'Stop Impersonation',
		addMember: 'Add Member',
		inviteMember: 'Invite member',
		email: 'Email',
		role: 'Role',
		copyInvitationLink: 'Copy Invitation Link',
		invitationLinkCopied: 'Invitation link copied to clipboard.',
		createOrganization: 'Create Organization',

		// Auth pages
		auth: {
			signIn: {
				title: 'Sign in',
				description: 'Enter your details below to sign in to your account.',
				signInWithGoogle: 'Sign in with Google',
				orContinueWith: 'Or continue with',
				email: 'Email',
				password: 'Password',
				signInButton: 'Sign in',
				successMessage:
					'You have been signed in successfully! You will be redirected to the dashboard soon.',
				errorMessage:
					'There was an error signing in. Make sure you have filled out all the fields correctly.',
				noAccount: "Don't have an account?",
				signUpLink: 'Sign up'
			},
			signUp: {
				title: 'Sign up',
				description: 'Enter your details below to sign up for an account.',
				signUpWithGoogle: 'Sign up with Google',
				orContinueWith: 'Or continue with',
				name: 'Name',
				email: 'Email',
				password: 'Password',
				signUpButton: 'Sign up',
				successMessage:
					'You have been signed up successfully! You will be redirected to the dashboard soon.',
				errorMessage:
					'There was an error signing up. Make sure you have filled out all the fields correctly.',
				hasAccount: 'Already have an account?',
				signInLink: 'Sign in'
			}
		},

		nav: {
			application: 'Application',
			dashboard: 'Dashboard',
			foo: 'Foo',
			bar: 'Bar',
			settings: 'Settings',
			switchOrganization: 'Switch Organization',
			organization: 'Organization',
			profile: 'Profile',
			logout: 'Logout',
			stopImpersonation: 'Stop Impersonation'
		},

		// Settings
		settings: {
			title: 'Settings',
			general: 'General',
			members: 'Members',
			billing: 'Billing',
			payments: 'Payments',
			profile: 'Profile',
			security: 'Security',
			fields: {
				email: {
					label: 'Email',
					placeholder: 'Email',
					description: 'The email you use to log into your account and receive notifications'
				},
				updatePassword: {
					label: 'Update Password',
					placeholder: 'Password',
					description: ' Update your password used to sign in to your account',
					passwordTooShort: 'Password must be at least 8 characters long'
				},
				confirmPassword: {
					label: 'Confirm Password',
					placeholder: 'Confirm Password',
					passwordDoesNotMatch: 'Passwords do not match'
				}
			},
			success: {
				profileUpdated: 'Your profile has been updated!',
				profileUpdateError:
					'There was an error updating your profile. Make sure you have filled out all the fields correctly.'
			},
			organization: {
				title: 'Organization Settings',
				general: {
					title: 'General',
					logo: {
						label: 'Organization logo',
						description: 'Upload your organization logo',
						noLogoUploaded: 'No logo uploaded',
						updateLogo: 'Update logo',
						updateLogoDescription: 'Edit the logo to your liking and click save.'
					},
					name: {
						label: 'Name',
						placeholder: 'Organization name',
						description: 'What should we call your organization?'
					},
					description: {
						label: 'Description',
						placeholder: 'Organization description',
						description: 'A brief description of your organization'
					},
					success: {
						organizationUpdated: 'Your organization has been updated!',
						organizationUpdateError:
							'There was an error updating your organization. Make sure you have filled out all the fields correctly.'
					}
				},
				members: {
					title: 'Members',
					memberInformation: 'Member Information',
					dangerZone: 'Danger Zone',
					role: {
						label: 'Role',
						admin: 'Admin',
						member: 'Member'
					},
					success: {
						memberUpdated: 'The member has been updated!',
						memberUpdateError:
							'There was an error updating the member. Make sure you have filled out all the fields correctly.'
					},
					removeMember: 'Remove member'
				}
			}
		},

		// Onboarding
		onboarding: {
			back: 'Back',
			organizationInformation: 'Organization Information',
			contactInformation: 'Contact Information',
			address: 'Address',
			businessDetails: 'Business Details',
			reviewDetails: 'Review Details',
			addressStep: {
				country: 'Country',
				selectCountry: 'Select your country',
				address: 'Address',
				enterStreetAddress: 'Enter your street address',
				addressLine2: 'Address Line 2',
				addressLine2Description:
					'If you have a second line of address (e.g. apartment or suite number), enter it here',
				city: 'City',
				enterCity: 'Enter your city',
				state: 'State',
				selectStateOrProvince: 'Select your state or province',
				zip: 'Zip',
				enterZip: 'Enter your zip'
			},
			basicInformationStep: {
				name: 'Name',
				nameDescription:
					'Enter a name for your organization. This will be shown on invoices and other communications.',
				description: 'Description',
				descriptionDescription: 'Enter a description for your organization',
				logo: 'Logo',
				logoPreview: 'Organization logo preview',
				noLogoSelected: 'No logo selected',
				uploadLogo: "Upload your organization's logo (optional)",
				updateLogo: 'Update logo',
				editLogoDescription: 'Edit the logo to your liking and click save.',
				close: 'close'
			},
			businessDetailsStep: {
				registeredBusiness: 'Is this organization associated with a registered business?',
				yes: 'Yes',
				no: 'No',
				registeredBusinessDescription:
					"Select 'Yes' if you are a registered business (e.g. LLC, corporation, etc.). Otherwise, select 'No'.",
				language: 'Language',
				languageDescription:
					"Select the language you want to use for your organization. We will use this language for your organization's invoices and any communications we send to you."
			},
			contactInformationStep: {
				contactEmail: 'Contact Email',
				email: 'Email',
				emailDescription: 'Enter an email address we can contact your organization at.',
				contactPhone: 'Contact Phone',
				selectCountry: 'Select your country',
				phoneNumber: 'Phone number',
				phoneDescription: 'Enter a phone number we can contact your organization at.'
			},
			reviewDetailsStep: {
				basicInformation: 'Basic information',
				contact: 'Contact',
				address: 'Address',
				businessDetails: 'Business details',
				name: 'Name',
				description: 'Description',
				email: 'Email',
				phone: 'Phone',
				addressLine1: 'Address line 1',
				addressLine2: 'Address line 2',
				city: 'City',
				stateProvince: 'State/Province',
				zipPostalCode: 'ZIP/Postal code',
				country: 'Country',
				organizationType: 'Organization type',
				language: 'Language',
				logoPreview: 'Logo preview'
			},
			stripeAlert: {
				missingRequirements:
					'To accept payments from your customers, Stripe needs more information about your business.',
				adminPermissionsRequired: 'You need admin permissions to complete Stripe onboarding.',
				openStripe: 'Open Stripe',
				settingUp:
					'We are setting up your Stripe account so that you can accept payments from your customers. Please check back later.',
				somethingWentWrong: 'Something went wrong. Please try again.'
			}
		},

		// Billing
		billing: {
			title: 'Billing & Subscription',
			proPlan: 'Pro Plan',
			freePlan: 'Free Plan',
			active: 'Active',
			current: 'Current',
			recommended: 'Recommended',
			proDescription: "You're subscribed to the Pro plan with all premium features.",
			freeDescription:
				"You're currently on the Free plan. Upgrade to Pro to unlock advanced features and grow your business.",
			upgradeToPro: 'Upgrade to Pro Now',
			getStartedInMinutes: 'Get started in minutes',
			billingAmount: 'Billing Amount',
			nextBillingDate: 'Next Billing Date',
			basicFeatures: 'Basic features',
			standardSupport: 'Standard support',
			communityAccess: 'Community access',
			allFreeFeatures: 'All Free features',
			advancedFeatures: 'Advanced features',
			prioritySupport: 'Priority support',
			customIntegrations: 'Custom integrations',
			getPro: 'Get Pro',
			manageSubscription: 'Manage Subscription',
			failedToStartCheckout: 'Failed to start checkout. Please try again.',
			failedToOpenBillingPortal: 'Failed to open billing portal. Please try again.',
			perMonth: '/month'
		},

		// Payments
		payments: {
			title: 'Payments',
			description:
				'Manage your payment settings and view transaction history in your Stripe dashboard.',
			openStripeDashboard: 'Open Stripe Dashboard',
			failedToOpenStripeDashboard: 'Failed to open Stripe dashboard. Please try again.',
			redirectingToStripe: 'Redirecting to Stripe...'
		},

		// Members
		members: {
			title: 'Members',
			name: 'Name',
			role: 'Role',
			noMembershipsFound: 'No memberships found',
			pendingInvitations: 'Pending Invitations',
			email: 'Email',
			noInvitationsFound: 'No invitations found',
			invitationSent: "We've sent an email to",
			withInstructionsToJoin: 'with instructions to join your organization.'
		},

		// Profile page
		profile: {
			title: 'Profile',
			profilePicture: 'Profile picture',
			uploadProfilePicture: 'Upload your profile picture',
			name: 'Name',
			namePlaceholder: 'Name',
			nameDescription: 'What should we call you?',
			profileUpdated: 'Your profile has been updated!',
			profileUpdateError:
				'There was an error updating your profile. Make sure you have filled out all the fields correctly.',
			updateImage: 'Update image',
			editImageDescription: 'Edit the image to your liking and click save.'
		},

		// Footer
		footer: {
			description:
				'Production-ready SvelteKit + Go starter template for building SaaS applications',
			copyright: 'Copyright © 2025 - All rights reserved',
			pricing: 'Pricing',
			faq: 'FAQ',
			github: 'GitHub'
		},

		// Organization roles
		roles: {
			admin: {
				title: 'Admin',
				description: 'Manage organization settings and manage members'
			},
			member: {
				title: 'Member',
				description: 'Manage XYZ'
			}
		},

		// OAuth
		oauth: {
			completingSignIn: 'Completing sign in...',
			pleaseWait: 'Please wait while we finish signing you in with Google.',
			authenticationError: 'Authentication Error',
			tryAgain: 'Try Again'
		},

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
				getCode: 'Get the code',
				faq: 'See FAQ',
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

		// No Organization Pages
		noOrganization: {
			title: 'Select Organization',
			organizations: 'Organizations',
			selectOrganization: 'Select an organization to continue to the app.',
			noOrganizations:
				'You are not a member of any organizations. Create or join an organization to get started.',
			createOrganization: 'Create Organization',
			noDescription: 'No description',
			invitation: {
				title: 'Organization Invitation',
				invitedBy: '{inviterName} invited you to join "{organizationName}"',
				invitationDescription:
					'By accepting, you will be added to the organization, and you will be able to collaborate with your team.',
				accept: 'Accept',
				decline: 'Decline',
				accepted: {
					title: 'This invitation has already been accepted.',
					description:
						'If you did not accept this invitation, please contact the organization owner for a new invitation.'
				},
				declined: {
					title: 'This invitation has already been declined.',
					description:
						'If you would like to join this organization, please contact the organization owner for a new invitation.'
				},
				expired: {
					title: 'This invitation has expired.',
					description:
						'If you would like to join this organization, please contact the organization owner for a new invitation.'
				},
				revoked: {
					title: 'This invitation has already been revoked.',
					description:
						'If you would like to join this organization, please contact the organization owner for a new invitation.'
				}
			},
			admin: {
				title: 'Admin Portal',
				welcome: 'Welcome to the admin portal',
				users: {
					title: 'Users',
					searchPlaceholder: 'Search users...',
					search: 'Search',
					name: 'Name',
					email: 'Email',
					impersonate: 'Impersonate',
					previous: 'Previous',
					next: 'Next'
				},
				debug: {
					title: 'Debug',
					userScopes: 'User Scopes'
				}
			}
		},

		// Create Organization Pages
		createOrganizationPages: {
			title: 'Create Organization',
			steps: {
				basicInformation: 'Organization Information',
				contactInformation: 'Contact Information',
				address: 'Address',
				businessDetails: 'Business Details',
				review: 'Review Details'
			},
			descriptions: {
				basicInformation: {
					intro:
						"Let's start with the basics. Give your organization a name. You can also add a logo and description if you would like.",
					note: 'You can always change this information later!'
				},
				contactInformation: {
					intro:
						"Now, we need to know how to contact your organization. We'll use this information to send you important updates and notifications. We may also use this information to confirm your identity and keep your account secure."
				},
				address: {
					intro:
						'Next, we need to know where your organization is located. This will be displayed on invoices, used to calculate taxes, and more.',
					note: 'If your address changes, you can always update it later!'
				},
				businessDetails: {
					intro:
						'We need some final details about your organization. We need this information to ensure legal compliance and provide you the best possible experience.'
				},
				review: {
					intro:
						"Lastly, let's review the details of your organization to make sure everything is correct!"
				}
			},
			progress: {
				stepOf: '{current} of {total}'
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
						'Modern react libraries are far too complicated. Nowadays many are VC backed and want you to use their cloud hosting service which introduces even more complexity. At the same time, SvelteKit provides a nicer interface for reactive UIs and supports impler data loading and mutation patterns.'
				},
				whyDaisyui: {
					question: 'Why DaisyUI for styling?',
					answer:
						"I thought it was neat and easy to use. If you don't like DaisyUI styles, you can always defined your own components and ignore DaisyUI altogether."
				},
				whyGo: {
					question: 'Why Go for the backend?',
					answer:
						'Go has what I need for most of my projects: minimal resource usage, fast build times, simple developer experience.'
				},
				whyPostgresql: {
					question: 'Why PostgreSQL?',
					answer:
						'Would you prefer I use DB2, MongoDB? Postgres is fast, easily available, and generally a good SQL database.'
				},
				whyRailway: {
					question: 'Why Railway for deployment?',
					answer:
						"I'm lazy. Maybe in the future I'll use uncloud or k8s, but railway works fine and is cheap."
				},
				whyCustomAuth: {
					question: 'Why custom authentication instead of Auth0 or Clerk?',
					answer:
						'My authentication models have always slightly difered from Auth0 or Clerk which has resulted in time wasted trying to make them work for my use case. Here you own the auth model so you can hack it to your hearts content with minimal issues. For example, this enables B2B2B models.'
				},
				whatIsOrgModel: {
					question: 'What is the organization-member-user model?',
					answer:
						'There are organizations, and users. Organizations own resources and users must be a member of the organization to access or modify the resource.'
				},
				whyMakefile: {
					question: 'Why use a Makefile for development?',
					answer: "Don't need fancy tools."
				},
				howToStart: {
					question: 'How do I get started with development?',
					answer: 'Clone the repo and open your text editor.'
				},
				withoutOrgs: {
					question: 'Can I use this without the organization features?',
					answer: 'Yes, you will just need to remove organizations from the codebase.'
				},
				addOauth: {
					question: 'How do I add new OAuth providers?',
					answer:
						'The project includes Google OAuth as an example. To add new providers, update the OAuth configuration in the backend, add the provider-specific endpoints, and create corresponding frontend sign-in buttons. The authentication system is designed to support multiple OAuth providers simultaneously.'
				},
				otherPlatforms: {
					question: 'Can I deploy this to other platforms besides Railway?',
					answer: 'Yes, there is nothing specific to railway here.'
				},
				billingIncluded: {
					question: 'Is billing/payments functionality included?',
					answer: 'Yes, Stripe is included and comes with a subscription-based model by default.'
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
				cta: 'Go Pro Now'
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
						'We accept all major credit cards (Visa, MasterCard, American Express). All payments are processed securely through Stripe.'
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

// https://stackoverflow.com/questions/65332597/typescript-is-there-a-recursive-keyofa
export type RecursiveKeyOf<TObj extends object> = {
	[TKey in keyof TObj & (string | number)]: TObj[TKey] extends unknown[]
		? `${TKey}`
		: TObj[TKey] extends object
			? `${TKey}.${RecursiveKeyOf<TObj[TKey]>}`
			: `${TKey}`;
}[keyof TObj & (string | number)];

// Use english translations as the base type for translation keys
export type TranslationKey = RecursiveKeyOf<typeof translations.en>;
