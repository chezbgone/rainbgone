<script lang="ts">
	import Stripes from '$lib/common/Stripes.svelte';
	import { formatDateKey, formatHourMinute, formatUnixTime } from '$lib/common/time';
	import type { Forecast } from '../Forecast/types';

	interface Props {
		daily: Forecast['daily']['data'][number];
		hourly: Forecast['hourly']['data'];
		minTemp: number;
		maxTemp: number;
		today: boolean;
		lat: number;
		lng: number;
		timezone: string;
	}

	let { daily, hourly, minTemp, maxTemp, today, lat, lng, timezone }: Props = $props();

	// Pirate Weather can return 'thunderstorm', which we don't have an icon for yet.
	const dailyIcon = daily.icon === 'thunderstorm' ? 'rain' : daily.icon;

	const getDayName = (dt: number) => {
		if (today) return 'Today';
		return new Intl.DateTimeFormat(undefined, { timeZone: timezone, weekday: 'short' }).format(
			dt * 1000
		);
	};

	const left = ((daily.temperatureMin - minTemp) / (maxTemp - minTemp)) * 100;
	const width = ((daily.temperatureMax - daily.temperatureMin) / (maxTemp - minTemp)) * 100;
	const right = ((daily.temperatureMax - minTemp) / (maxTemp - minTemp)) * 100;

	let opened = $state(false);
	const detailsHref = $derived(`/details/${lat},${lng}/${formatDateKey(daily.time, timezone)}`);
</script>

<details bind:open={opened}>
	<summary class="group flex items-center px-4 py-2 hover:cursor-pointer hover:bg-gray-50">
		<img src={`/weather-icons/${dailyIcon}.png`} alt={dailyIcon} class="h-8" />
		<span class="mx-2 w-16">{getDayName(daily.time)}</span>
		<span class="relative mx-16 flex h-4 w-lg items-center">
			<span class="absolute ml-[-30px] text-right" style="left: {left}%"
				>{Math.round(daily.temperatureMin)}°</span
			>
			<span
				class="absolute left-0 h-4 rounded-full bg-gray-700"
				style="
          margin-left: {left}%;
          width: {width}%;
        "
			></span>
			<span class="absolute ml-[5px] text-left" style="left: {right}%"
				>{Math.round(daily.temperatureMax)}°</span
			>
		</span>
		<svg
			class="fill-neutral-500 group-hover:opacity-80"
			width="24"
			height="24"
			viewBox="0 0 24 24"
			xmlns="http://www.w3.org/2000/svg"
		>
			<defs>
				<mask id="plus-mask">
					<circle cx="12" cy="12" r="12" fill="white" />
					<rect x="11" y="6" width="2" height="12" fill="black" />
					<rect x="6" y="11" width="12" height="2" fill="black" />
				</mask>
				<mask id="minus-mask">
					<circle cx="12" cy="12" r="12" fill="white" />
					<rect x="6" y="11" width="12" height="2" fill="black" />
				</mask>
			</defs>
			<circle cx="12" cy="12" r="12" mask={opened ? 'url(#minus-mask)' : 'url(#plus-mask)'} />
		</svg>
	</summary>
	<div class="border-b border-neutral-400 py-2">
		<div class="text-xl font-light">{daily.summary}</div>
		<div class="my-2 inline-flex items-center gap-16 text-neutral-800">
			<div>
				<span class="text-2xl">{Math.round(daily.temperatureMin)}°</span>
				<span class="font-light">{formatUnixTime(daily.temperatureMinTime, timezone)}</span>
				<span class="mx-2 text-2xl font-light">&#8594;</span>
				<span class="text-2xl">{Math.round(daily.temperatureMax)}°</span>
				<span class="font-light">{formatUnixTime(daily.temperatureMaxTime, timezone)}</span>
			</div>
			<div class="font-light">
				<span>
					<img src="sunrise.png" alt="Sunrise icon" class="inline h-6 w-6" />
					<span>{formatHourMinute(daily.sunriseTime, timezone)}</span>
				</span>
				<span>
					<img src="sunset.png" alt="Sunset icon" class="ml-4 inline h-6 w-6" />
					<span>{formatHourMinute(daily.sunsetTime, timezone)}</span>
				</span>
			</div>
			<div>
				<span class="mr-1">Rain</span>
				<span class="font-light">{daily.precipIntensity.toFixed(2)}</span>
				<span class="font-light">in</span>
			</div>
		</div>
		<Stripes hours={hourly} {timezone} />
		<a
			href={detailsHref}
			class="mb-4 inline-block rounded-sm bg-blue-500 px-4 py-2 text-white uppercase hover:bg-blue-600"
		>
			more details
		</a>
	</div>
</details>
