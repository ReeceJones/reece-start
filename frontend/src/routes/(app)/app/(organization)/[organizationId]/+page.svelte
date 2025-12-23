<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { t } from '$lib/i18n';
</script>

<div class="grid grid-cols-6 gap-6">
	<Card.Root class="col-span-4">
		<Card.Header>
			<Card.Title>{$t('foo')}</Card.Title>
		</Card.Header>
	</Card.Root>
	<Card.Root class="col-span-2">
		<Card.Header>
			<Card.Title>{$t('bar')}</Card.Title>
		</Card.Header>
	</Card.Root>
</div>
