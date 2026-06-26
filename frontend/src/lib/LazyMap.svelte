<script lang="ts">
	import { onMount } from 'svelte';
	import type Map from '$lib/Map.svelte';

	interface Props {
		location: {
			lat: number;
			lng: number;
		};
		precipitationSoon: boolean;
	}

	let { location, precipitationSoon }: Props = $props();
	let MapComponent: typeof Map | null = $state(null);

	onMount(async () => {
		MapComponent = (await import('$lib/Map.svelte')).default;
	});
</script>

{#if MapComponent}
	<MapComponent {location} {precipitationSoon} />
{:else}
	<div
		class="relative max-h-[800px] min-h-[350px] w-full bg-neutral-200 after:block after:pt-[35%]"
	></div>
{/if}
