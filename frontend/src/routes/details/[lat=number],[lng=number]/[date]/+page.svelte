<script lang="ts">
	import Stripes from '$lib/common/Stripes.svelte';
	import DayNavigation from '$lib/Details/DayNavigation.svelte';
	import DayHighlights from '$lib/Details/DayHighlights.svelte';
	import InstantDetails from '$lib/Details/InstantDetails.svelte';
	import HourlyCharts from '$lib/Details/Charts/HourlyCharts.svelte';
	import Scrubber from '$lib/Details/Scrubber.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<main class="mx-auto max-w-[800px] px-4 pb-12 text-neutral-800">
	<DayNavigation
		previous={data.links.previous}
		previousLabel={data.links.previousLabel}
		next={data.links.next}
		nextLabel={data.links.nextLabel}
		time={data.day.time}
		timezone={data.timezone}
	/>

	<p class="mx-1 mt-5 mb-8 text-center text-[2rem] leading-tight font-light">{data.day.summary}</p>

	<DayHighlights day={data.day} timezone={data.timezone} />

	{#if data.hourly.length > 0}
		<Scrubber hours={data.hourly}>
			{#snippet children(hour)}
				<Stripes hours={data.hourly} timezone={data.timezone} />
				<InstantDetails hours={data.hourly} {hour} />
				<HourlyCharts hours={data.hourly} timezone={data.timezone} />
			{/snippet}
		</Scrubber>
	{:else}
		<p class="mt-6 text-center text-neutral-500">
			No hourly forecast data is available for this day.
		</p>
	{/if}
</main>
