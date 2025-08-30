<script lang="ts">
	const { children } = $props();
</script>

<div class="card rounded-box bg-base-200 shadow-sm">
	<div class="card-body">
		{@render children?.()}
	</div>
</div>
