<script lang="ts">
	import Day from './Day.svelte';
	import type { Forecast } from '../Forecast/types';

	interface Props {
		daily: Forecast['daily'];
		hourly: Forecast['hourly'];
		lat: number;
		lng: number;
		timezone: string;
	}

	let { daily, hourly, lat, lng, timezone }: Props = $props();

	const days = $derived(daily.data);
	const daysHourly = $derived.by(() => {
		const weekHours = [];
		for (let d = 0; d < 7; ++d) {
			weekHours.push(hourly.data.slice(d * 24, (d + 1) * 24 + 1));
		}
		return weekHours;
	});

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
				hourly={daysHourly[i]}
				{lat}
				{lng}
				{timezone}
			/>
		{/each}
	</div>
</div>
