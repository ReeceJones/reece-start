<script lang="ts">
	const { children } = $props();
</script>

<div class="card rounded-box bg-base-200 border-base-300 border shadow-sm">
	<div class="card-body">
		{@render children?.()}
	</div>
</div>
