<script lang="ts">
	import UserNav from '$lib/components/Nav/UserNav.svelte';
	import OrganizationNav from '$lib/components/Nav/OrganizationNav.svelte';
	import { setScopes } from '$lib/auth.js';
	import OnboardingStripeAlert from '$lib/components/Onboarding/OnboardingStripeAlert.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';

	const { children, data } = $props();

	const { user, organization, userScopes } = $derived(data);

	// Give the javascript client access to the scopes stored in the token cookie
	setScopes(() => userScopes);
</script>

<Sidebar.Provider>
	<Sidebar.Root>
		<Sidebar.Content>
			<Sidebar.Group>
				<OrganizationNav {organization} />
			</Sidebar.Group>
		</Sidebar.Content>
		<Sidebar.Footer>
			<Sidebar.Group>
				<UserNav {user} />
			</Sidebar.Group>
		</Sidebar.Footer>
	</Sidebar.Root>
	<main class="container m-4 mx-auto flex-1 px-4">
		<OnboardingStripeAlert {organization} />
		{@render children?.()}
	</main>
</Sidebar.Provider>
