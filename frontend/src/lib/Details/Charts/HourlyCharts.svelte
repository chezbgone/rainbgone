<script lang="ts">
	import LineChart from './LineChart.svelte';
	import type { Forecast } from '$lib/Forecast/types';

	type HourlyDatum = Forecast['hourly']['data'][number];

	interface Props {
		hours: HourlyDatum[];
	}

	let { hours }: Props = $props();

	interface SeriesConfig {
		value: (h: HourlyDatum) => number;
		color: string;
	}

	interface ChartConfig {
		label: string;
		unit: string;
		fromZero: boolean;
		minRange: number;
		max?: number;
		series: SeriesConfig[];
	}

	const solid = '#111111';
	const faint = 'rgba(0, 0, 0, 0.25)';

	const charts: ChartConfig[] = [
		{
			label: 'Temperature',
			unit: '°',
			fromZero: false,
			minRange: 10,
			series: [
				{ value: (h) => h.temperature, color: solid },
				{ value: (h) => h.apparentTemperature, color: faint }
			]
		},
		{
			label: 'Precipitation Chance',
			unit: '%',
			fromZero: true,
			minRange: 0,
			max: 100,
			series: [{ value: (h) => h.precipProbability * 100, color: solid }]
		},
		{
			label: 'Humidity',
			unit: '%',
			fromZero: true,
			minRange: 0,
			max: 100,
			series: [{ value: (h) => h.humidity * 100, color: solid }]
		},
		{
			label: 'Dew Point',
			unit: '°',
			fromZero: false,
			minRange: 5,
			series: [{ value: (h) => h.dewPoint, color: solid }]
		},
		{
			label: 'Wind Speed',
			unit: ' mph',
			fromZero: true,
			minRange: 5,
			series: [{ value: (h) => h.windSpeed, color: solid }]
		},
		{
			label: 'Pressure',
			unit: ' mb',
			fromZero: false,
			minRange: 0,
			series: [{ value: (h) => h.pressure, color: solid }]
		},
		{
			label: 'UV Index',
			unit: '',
			fromZero: true,
			minRange: 8,
			series: [{ value: (h) => h.uvIndex, color: solid }]
		},
		{
			label: 'Visibility',
			unit: ' mi',
			fromZero: true,
			minRange: 8,
			max: 10,
			series: [{ value: (h) => h.visibility, color: solid }]
		}
	];
</script>

<div class="mt-8 flex flex-col gap-8">
	{#each charts as chart (chart.label)}
		<LineChart
			{hours}
			label={chart.label}
			unit={chart.unit}
			fromZero={chart.fromZero}
			minRange={chart.minRange}
			max={chart.max}
			series={chart.series}
		/>
	{/each}
</div>
