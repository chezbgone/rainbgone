<script lang="ts">
	import Day from './Day.svelte';
	import type { Forecast } from '../Forecast/types';

	interface Props {
		daily: Forecast['daily'];
		hourlyFromMidnight: Forecast['hourlyFromMidnight'];
		lat: number;
		lng: number;
		timezone: string;
	}

	let { daily, hourlyFromMidnight, lat, lng, timezone }: Props = $props();

	const days = $derived(daily.data);

	const minTemp = $derived(Math.min(...days.map((d) => d.temperatureMin)));
	const maxTemp = $derived(Math.max(...days.map((d) => d.temperatureMax)));
</script>

<div class="m-8 text-center">
	<div class="m-4 text-xl font-light">{daily.summary}</div>
	<div class="flex flex-col items-center gap-2">
		{#each days.slice(0, 7) as day, i}
			<Day
				today={i === 0}
				daily={day}
				{minTemp}
				{maxTemp}
				hourly={hourlyFromMidnight.slice(i * 24, i * 24 + 25)}
				{lat}
				{lng}
				{timezone}
			/>
		{/each}
	</div>
</div>
