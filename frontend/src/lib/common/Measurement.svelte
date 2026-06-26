<script lang="ts">
	let { label, children } = $props();
</script>

<div>
	<span class="font-medium">{label}:</span>
	{#if label === 'UV Index'}
		<span class="rounded-sm bg-emerald-200 px-1 py-0.5 font-light">
			{@render children()}
		</span>
		<!-- TODO: wind direction icon -->
	{:else}
		<span class="font-light">
			{@render children()}
		</span>
	{/if}
</div>
