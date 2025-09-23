<script lang="ts">
	import { setIsImpersonatingUser, setScopes } from '$lib/auth';

	const { children, data } = $props();

	const { userScopes, isImpersonatingUser } = data;

	setScopes(userScopes);
	setIsImpersonatingUser(isImpersonatingUser);
</script>

{@render children?.()}
