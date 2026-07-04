<script lang="ts">
	import Current from '$lib/Current.svelte';
	import Hourly from '$lib/Hourly.svelte';
	import Daily from '$lib/Daily/Daily.svelte';
	import LazyMap from '$lib/LazyMap.svelte';
	import type { Geocode, Forecast } from './types';
	import { classifyWeather, Weather } from '$lib/common/weather';
	import Minutely from '$lib/Minutely.svelte';

	interface Props {
		geocode: Geocode;
		forecast: Forecast;
	}

	let { forecast, geocode }: Props = $props();

	const precipitationSoon = $derived.by(() => {
		if (!forecast?.hourly?.data) return false;
		return forecast.hourly.data.slice(0, 12).some((hour) => {
			const weather = classifyWeather(hour.precipIntensity, hour.precipType, hour.cloudCover);
			return [
				Weather.Rain,
				Weather.LightRain,
				Weather.Snow,
				Weather.LightSnow,
				Weather.Sleet,
				Weather.LightSleet
			].includes(weather);
		});
	});
</script>

<Current currently={forecast.currently} daily={forecast.daily} />
<Minutely minutely={forecast.minutely} />
<Hourly hourly={forecast.hourly} timezone={forecast.timezone} />
<LazyMap location={geocode.geometry.location} {precipitationSoon} />
<Daily
	daily={forecast.daily}
	hourlyFromMidnight={forecast.hourlyFromMidnight}
	lat={geocode.geometry.location.lat}
	lng={geocode.geometry.location.lng}
	timezone={forecast.timezone}
/>
