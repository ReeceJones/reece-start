<script lang="ts">
	import { getScopes } from '$lib/auth';

	const scopes = $derived(getScopes());
</script>

<div class="grid grid-cols-3 gap-6">
	<div class="card bg-base-200 border-base-300 border shadow-sm">
		<div class="card-body">
			<h2 class="card-title">User Scopes</h2>
			<pre>{JSON.stringify(scopes, null, 2)}</pre>
		</div>
	</div>
</div>
