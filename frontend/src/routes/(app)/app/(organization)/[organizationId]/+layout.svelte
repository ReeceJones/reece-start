<script lang="ts">
	import UserNav from '$lib/components/Nav/UserNav.svelte';
	import OrganizationNav from '$lib/components/Nav/OrganizationNav.svelte';
	import { setScopes } from '$lib/auth.js';
	import OnboardingStripeAlert from '$lib/components/Onboarding/OnboardingStripeAlert.svelte';

	const { children, data } = $props();

	const { user, organization, userScopes } = data;

	// Give the javascript client access to the scopes stored in the token cookie
	setScopes(userScopes);
</script>

<div class="flex h-screen max-w-screen flex-row gap-4">
	<div class="flex h-full w-56 flex-col justify-between border-r border-base-300 bg-base-200">
		<div>
			<OrganizationNav {organization} />
		</div>
		<div>
			<UserNav {user} />
		</div>
	</div>
	<main class="container mx-auto mt-4 mr-4 flex-1">
		<OnboardingStripeAlert {organization} />
		{@render children?.()}
	</main>
</div>
