<script lang="ts">
	import { setIsImpersonatingUser, setScopes } from '$lib/auth';
	import { t } from '$lib/i18n';

	const { children, data } = $props();

	const { userScopes, isImpersonatingUser } = $derived(data);

	setScopes(() => userScopes);
	setIsImpersonatingUser(() => isImpersonatingUser);
</script>

<svelte:head>
	<title>{$t('dashboard')} - reece-start</title>
	<meta name="description" content="Dashboard for reece-start" />
</svelte:head>

{@render children?.()}
