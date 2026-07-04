<script lang="ts">
	import Measurement from '$lib/common/Measurement.svelte';
	import type { Forecast } from './Forecast/types';

	interface Props {
		currently: Forecast['currently'];
		daily: Forecast['daily'];
	}

	let { currently, daily }: Props = $props();

	// Pirate Weather can return 'thunderstorm', which we don't have an icon for yet.
	const currentIcon = currently.icon === 'thunderstorm' ? 'rain' : currently.icon;
</script>

<div class="flex justify-center gap-4 bg-neutral-100 py-2 text-sm">
	<Measurement label="Wind">
		{currently.windSpeed.toFixed(1)} mph
	</Measurement>
	<Measurement label="Humidity">
		{currently.humidity * 100}%
	</Measurement>
	<Measurement label="Dew Point">
		{currently.dewPoint.toFixed(1)}°
	</Measurement>
	<Measurement label="UV Index">
		{currently.uvIndex}
	</Measurement>
	<Measurement label="Visibility">
		{currently.visibility.toFixed(1)} mi
	</Measurement>
	<Measurement label="Pressure">
		{Math.round(currently.pressure)} mb
	</Measurement>
</div>

<div class="mt-4 text-center">
	<div class="inline-flex items-center gap-4">
		<img width="84" height="84" src={`/weather-icons/${currentIcon}.png`} alt="" />
		<div>
			<div class="text-4xl/12 font-semibold">
				{Math.round(currently.temperature)}° {currently.summary}.
			</div>
			<div class="flex gap-2 text-sm">
				<Measurement label="Feels Like">
					{Math.round(currently.apparentTemperature)}°
				</Measurement>
				<Measurement label="Low">
					{Math.round(daily.data[0].temperatureLow)}°
				</Measurement>
				<Measurement label="High">
					{Math.round(daily.data[0].temperatureHigh)}°
				</Measurement>
			</div>
		</div>
	</div>
</div>
