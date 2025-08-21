<script lang="ts">
	import { House, Folder, DollarSign } from 'lucide-svelte';
</script>

<ul class="menu menu-vertical lg:menu-horizontal bg-base-200 rounded-box mt-2">
	<li>
		<a href="/app">
			<House class="h-4 w-4" />
			Dashboard
		</a>
	</li>
	<li>
		<a href="/app/foo">
			<Folder class="h-4 w-4" />
			Foo
		</a>
	</li>
	<li>
		<a href="/app/bar">
			<DollarSign class="h-4 w-4" />
			Bar
		</a>
	</li>
</ul>
