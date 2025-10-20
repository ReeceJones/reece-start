<script lang="ts">
	import { countries, getEmojiFlag } from 'countries-list';
</script>

{#each Object.entries(countries) as [countryCode, country] (countryCode)}
	<option value={countryCode}>
		{country.name}
		{getEmojiFlag(countryCode as keyof typeof countries)}
	</option>
{/each}
