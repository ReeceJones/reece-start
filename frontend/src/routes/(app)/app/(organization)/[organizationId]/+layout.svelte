<script lang="ts">
	import UserNav from '$lib/components/UserNav.svelte';
	import OrganizationNav from '$lib/components/OrganizationNav.svelte';

	const { children } = $props();
</script>

<div class="container mx-auto">
	<div class="flex justify-between">
		<div>
			<OrganizationNav />
		</div>
		<div>
			<UserNav />
		</div>
	</div>
	<main class="mt-8">
		{@render children?.()}
	</main>
</div>
