<script lang="ts">
	import { number, percent } from '$lib/common/format';
	import { uvLevel } from '$lib/common/weather';
	import type { DailyDatum } from '$lib/Forecast/types';

	interface Props {
		day: DailyDatum;
	}

	let { day }: Props = $props();

	const details = $derived([
		{ label: 'Precip', value: percent(day.precipProbability) },
		{ label: 'Pressure', value: `${number(day.pressure, 1)} mb` },
		{ label: 'Humidity', value: percent(day.humidity) },
		{ label: 'Dew Pt', value: `${number(day.dewPoint)}°` },
		{ label: 'Visibility', value: `${number(day.visibility, 1)} mi` }
	]);
</script>

<div class="mt-8 grid grid-cols-2 gap-x-8 gap-y-3 text-center sm:grid-cols-4">
	{#each details as detail}
		<div>
			<span class="font-medium">{detail.label}:</span>
			<span class="font-light">{detail.value}</span>
		</div>
	{/each}

	<!-- Wind, with a direction arrow rotated by bearing -->
	<div>
		<span class="font-medium">Wind:</span>
		<span class="font-light">{number(day.windSpeed, 1)} mph</span>
		<span
			class="ml-0.5 inline-block"
			style="transform: rotate({day.windBearing + 180}deg)"
			title="{Math.round(day.windBearing)}°"
		>
			↑
		</span>
	</div>

	<!-- UV index, with a color-coded chip -->
	<div>
		<span class="font-medium">UV Index:</span>
		<span class="ml-1 rounded px-1.5 py-0.5 font-light {uvLevel(day.uvIndex)}">
			{number(day.uvIndex)}
		</span>
	</div>
</div>
