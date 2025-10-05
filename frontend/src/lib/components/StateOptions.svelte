<script lang="ts">
	import { State } from 'country-state-city';

	const { countryCode } = $props();

	const states = $derived(State.getStatesOfCountry(countryCode));
</script>

{#each states as state (state.isoCode)}
	<option value={state.isoCode}>
		{state.name}
	</option>
{/each}
