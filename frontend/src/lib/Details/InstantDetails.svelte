<script lang="ts">
	import { number, percent } from '$lib/common/format';
	import { uvLevel } from '$lib/common/weather';
	import type { Forecast } from '$lib/Forecast/types';

	type HourlyDatum = Forecast['hourly']['data'][number];

	interface Props {
		hours: HourlyDatum[];
		hour: number; // fractional hour position from the scrubber
	}

	let { hours, hour }: Props = $props();

	const lerp = (a: number, b: number, t: number) => a + (b - a) * t;
	// Shortest-arc interpolation for a compass bearing (handles the 360°→0° wrap so the
	// wind arrow never spins the long way around).
	const lerpAngle = (a: number, b: number, t: number) => {
		const diff = ((b - a + 540) % 360) - 180;
		return (a + diff * t + 360) % 360;
	};

	// Interpolates every field between the two hours the scrubber sits between, except UV
	// index, which stays snapped to the nearest hour (it's a coarse, bucketed quantity).
	const values = $derived.by(() => {
		const i0 = Math.floor(hour);
		const i1 = Math.min(i0 + 1, hours.length - 1);
		const t = hour - i0;
		const lo = hours[i0];
		const hi = hours[i1];
		return {
			temperature: lerp(lo.temperature, hi.temperature, t),
			precipProbability: lerp(lo.precipProbability, hi.precipProbability, t),
			pressure: lerp(lo.pressure, hi.pressure, t),
			humidity: lerp(lo.humidity, hi.humidity, t),
			dewPoint: lerp(lo.dewPoint, hi.dewPoint, t),
			visibility: lerp(lo.visibility, hi.visibility, t),
			windSpeed: lerp(lo.windSpeed, hi.windSpeed, t),
			windBearing: lerpAngle(lo.windBearing, hi.windBearing, t),
			uvIndex: hours[Math.round(hour)].uvIndex
		};
	});

	const details = $derived([
		{ label: 'Temp', value: `${number(values.temperature)}°` },
		{ label: 'Precip', value: percent(values.precipProbability) },
		{ label: 'Pressure', value: `${number(values.pressure, 1)} mb` },
		{ label: 'Humidity', value: percent(values.humidity) },
		{ label: 'Dew Pt', value: `${number(values.dewPoint)}°` },
		{ label: 'Visibility', value: `${number(values.visibility, 1)} mi` }
	]);
</script>

<div class="mt-8 flex flex-wrap gap-y-2 text-center text-sm">
	{#each details as detail (detail.label)}
		<div class="basis-1/2 sm:basis-1/4">
			<span class="font-medium">{detail.label}:</span>
			<span class="font-light">{detail.value}</span>
		</div>
	{/each}

	<!-- UV index, with a color-coded chip. Time Machine's backfill for today's elapsed hours
	     can leave uvIndex missing; skip the color coding then so a missing reading doesn't
	     read as a real low-UV value. -->
	<div class="basis-1/2 sm:basis-1/4">
		<span class="font-medium">UV Index:</span>
		<span
			class="ml-1 rounded px-1.5 py-0.5 font-light {Number.isFinite(values.uvIndex)
				? uvLevel(values.uvIndex)
				: ''}"
		>
			{number(values.uvIndex)}
		</span>
	</div>

	<!-- Wind, with a direction arrow rotated by bearing -->
	<div class="basis-1/2 sm:basis-1/4">
		<span class="font-medium">Wind:</span>
		<span class="font-light">{number(values.windSpeed, 1)} mph</span>
		<svg
			class="ml-0.5 inline-block size-4 align-text-bottom"
			viewBox="0 0 24 24"
			style="transform: rotate({values.windBearing + 180}deg)"
			fill="none"
			stroke="currentColor"
			stroke-width="1.25"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<title>{Math.round(values.windBearing)}°</title>
			<line x1="12" y1="20" x2="12" y2="4" />
			<polyline points="6 10 12 4 18 10" />
		</svg>
	</div>
</div>
