<script lang="ts">
	import { countries, getEmojiFlag } from 'countries-list';
</script>

{#each Object.entries(countries) as [countryCode, country] (countryCode)}
	<option value={countryCode}>
		{getEmojiFlag(countryCode as keyof typeof countries)}
		{country.name}
	</option>
{/each}
