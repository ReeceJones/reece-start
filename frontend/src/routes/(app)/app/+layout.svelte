<script lang="ts">
	import { setIsImpersonatingUser, setScopes } from '$lib/auth';
	import * as m from '$lib/paraglide/messages';

	const { children, data } = $props();

	const { userScopes, isImpersonatingUser } = $derived(data);

	setScopes(() => userScopes);
	setIsImpersonatingUser(() => isImpersonatingUser);
</script>

<svelte:head>
	<title>{m.dashboard()} - reece-start</title>
	<meta name="description" content="Dashboard for reece-start" />
</svelte:head>

{@render children?.()}
