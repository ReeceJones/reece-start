<script lang="ts">
	import UserNav from '$lib/components/Nav/UserNav.svelte';
	import OrganizationNav from '$lib/components/Nav/OrganizationNav.svelte';
	import { setScopes } from '$lib/auth.js';

	const { children, data } = $props();

	const { user, organization, membershipScopes } = data;

	// Give the javascript client access to the scopes stored in the token cookie
	setScopes(membershipScopes);
</script>

<div class="max-w-screen flex h-screen flex-row gap-4">
	<div class="border-base-300 bg-base-200 flex h-full w-56 flex-col justify-between border-r">
		<div>
			<OrganizationNav {organization} />
		</div>
		<div>
			<UserNav {user} />
		</div>
	</div>
	<main class="container mx-auto mr-4 mt-4 flex-1">
		{@render children?.()}
	</main>
</div>
