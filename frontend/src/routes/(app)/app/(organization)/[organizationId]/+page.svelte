<script lang="ts">
	import Card from '$lib/components/Card/Card.svelte';
	import CardBody from '$lib/components/Card/CardBody.svelte';
	import CardTitle from '$lib/components/Card/CardTitle.svelte';
	import { t } from '$lib/i18n';
</script>

<div class="grid grid-cols-6 gap-6">
	<Card class="col-span-4">
		<CardBody>
			<CardTitle>{$t('foo')}</CardTitle>
		</CardBody>
	</Card>
	<Card class="col-span-2">
		<CardBody>
			<CardTitle>{$t('bar')}</CardTitle>
		</CardBody>
	</Card>
</div>
