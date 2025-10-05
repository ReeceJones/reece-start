<script lang="ts">
	import { getScopes } from '$lib/auth';
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';

	const scopes = $derived(getScopes());
</script>

<div class="grid grid-cols-3 gap-6">
	<Card>
		<CardBody>
			<CardTitle>User Scopes</CardTitle>
			<ul class="list-inside list-disc">
				{#each scopes as scope}
					<li class="font-mono">{scope}</li>
				{/each}
			</ul>
		</CardBody>
	</Card>
</div>
