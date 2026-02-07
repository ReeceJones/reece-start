<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as m from '$lib/paraglide/messages';
</script>

<div class="grid grid-cols-6 gap-6">
	<Card.Root class="col-span-4">
		<Card.Header>
			<Card.Title>{m.foo()}</Card.Title>
		</Card.Header>
	</Card.Root>
	<Card.Root class="col-span-2">
		<Card.Header>
			<Card.Title>{m.bar()}</Card.Title>
		</Card.Header>
	</Card.Root>
</div>
