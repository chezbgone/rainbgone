<script lang="ts">
	import Stripes from '$lib/common/Stripes.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const formatDate = (time: number) =>
		new Intl.DateTimeFormat(undefined, {
			timeZone: data.timezone,
			weekday: 'long',
			month: 'long',
			day: 'numeric',
			year: 'numeric'
		}).format(time * 1000);

	const formatHourMinute = (time: number) =>
		new Intl.DateTimeFormat(undefined, {
			timeZone: data.timezone,
			hour: 'numeric',
			minute: '2-digit'
		}).format(time * 1000);

	const number = (value: number | null | undefined, digits = 0) =>
		value === null || value === undefined ? 'N/A' : value.toFixed(digits);

	const percent = (value: number | null | undefined) =>
		value === null || value === undefined ? 'N/A' : `${Math.round(value * 100)}%`;

	// UV chip color, matching Dark Sky's uv0–uv4 buckets.
	const uvLevel = (uv: number) => {
		if (uv >= 11) return 'bg-[rgba(166,89,255,0.75)]';
		if (uv >= 8) return 'bg-[rgba(255,0,0,0.6)]';
		if (uv >= 6) return 'bg-[rgba(255,127,0,0.5)]';
		if (uv >= 3) return 'bg-[rgba(255,240,0,0.6)]';
		return 'bg-[rgba(64,191,64,0.6)]';
	};

	const precipType = $derived(data.day.precipType);
	const hasPrecip = $derived(precipType !== 'none');

	const details = $derived([
		{ label: 'Precip', value: percent(data.day.precipProbability) },
		{ label: 'Pressure', value: `${number(data.day.pressure, 1)} mb` },
		{ label: 'Humidity', value: percent(data.day.humidity) },
		{ label: 'Dew Pt', value: `${number(data.day.dewPoint)}°` },
		{ label: 'Visibility', value: `${number(data.day.visibility, 1)} mi` }
	]);
</script>

<main class="mx-auto max-w-[800px] px-4 pb-12 text-neutral-800">
	<div class="my-5 text-center text-xl font-light">
		<a href={data.links.back} class="text-neutral-600 hover:text-neutral-900">← Go Back</a>
	</div>

	<div class="mt-2 mb-10 flex items-baseline justify-center gap-6 text-center">
		{#if data.links.previous}
			<a
				href={data.links.previous}
				class="flex-1 text-right text-base text-blue-600 hover:text-blue-800"
			>
				← {data.links.previousLabel}
			</a>
		{:else}
			<span class="flex-1"></span>
		{/if}
		<div class="text-lg text-black">{formatDate(data.day.time)}</div>
		{#if data.links.next}
			<a
				href={data.links.next}
				class="flex-1 text-left text-base text-blue-600 hover:text-blue-800"
			>
				{data.links.nextLabel} →
			</a>
		{:else}
			<span class="flex-1"></span>
		{/if}
	</div>

	<p class="mx-1 mt-5 mb-8 text-center text-[2rem] leading-tight font-light">{data.day.summary}</p>

	<div class="flex flex-wrap items-end justify-center gap-x-12 gap-y-4 text-base">
		<!-- High / low temperature -->
		<div class="flex items-end">
			<span class="flex items-baseline gap-1">
				<span class="text-2xl font-medium">{Math.round(data.day.temperatureMax)}°</span>
				<span class="text-sm">{formatHourMinute(data.day.temperatureMaxTime)}</span>
			</span>
			<span class="mx-1 text-xl">→</span>
			<span class="flex items-baseline gap-1">
				<span class="text-2xl font-medium">{Math.round(data.day.temperatureMin)}°</span>
				<span class="text-sm">{formatHourMinute(data.day.temperatureMinTime)}</span>
			</span>
		</div>

		<!-- Sun times -->
		<div class="flex items-center gap-4">
			<span class="flex items-center gap-1">
				<img src="/sunrise.png" width="28" height="30" alt="Sunrise" />
				{formatHourMinute(data.day.sunriseTime)}
			</span>
			<span class="flex items-center gap-1">
				<img src="/sunset.png" width="28" height="30" alt="Sunset" />
				{formatHourMinute(data.day.sunsetTime)}
			</span>
		</div>

		<!-- Precipitation accumulation -->
		{#if hasPrecip}
			<div class="flex items-baseline gap-1">
				<span class="font-medium capitalize">{precipType}</span>
				<span>{number(data.day.precipAccumulation, 2)} in</span>
			</div>
		{/if}
	</div>

	{#if data.hourly.length > 0}
		<Stripes hours={data.hourly} />
	{:else}
		<p class="mt-6 text-center text-neutral-500">
			No hourly forecast data is available for this day.
		</p>
	{/if}

	<div class="mt-8 grid grid-cols-2 gap-x-8 gap-y-3 text-center sm:grid-cols-3">
		{#each details as detail}
			<div>
				<span class="font-medium">{detail.label}:</span>
				<span class="font-light">{detail.value}</span>
			</div>
		{/each}

		<!-- Wind, with a direction arrow rotated by bearing -->
		<div>
			<span class="font-medium">Wind:</span>
			<span class="font-light">{number(data.day.windSpeed, 1)} mph</span>
			<span
				class="ml-0.5 inline-block"
				style="transform: rotate({data.day.windBearing}deg)"
				title="{Math.round(data.day.windBearing)}°"
			>
				↑
			</span>
		</div>

		<!-- UV index, with a color-coded chip -->
		<div>
			<span class="font-medium">UV Index:</span>
			<span class="ml-1 rounded px-1.5 py-0.5 font-light {uvLevel(data.day.uvIndex)}">
				{number(data.day.uvIndex)}
			</span>
		</div>
	</div>
</main>
