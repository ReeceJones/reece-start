<script lang="ts">
	import { setIsImpersonatingUser, setScopes } from '$lib/auth';

	const { children, data } = $props();

	const { userScopes, isImpersonatingUser } = data;

	setScopes(userScopes);
	setIsImpersonatingUser(isImpersonatingUser);
</script>

<svelte:head>
	<title>Dashboard - reece-start</title>
	<meta name="description" content="Dashboard for reece-start" />
</svelte:head>

{@render children?.()}
