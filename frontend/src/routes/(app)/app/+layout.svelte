<script lang="ts">
	import { setScopes } from '$lib/auth';

	const { children, data } = $props();

	const { userScopes } = data;

	setScopes(userScopes);
</script>

{@render children?.()}
