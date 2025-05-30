<script lang="ts">
  import Current from '$lib/Current.svelte';
  import Hourly from '$lib/Hourly.svelte';
  import Daily from '$lib/Daily/Daily.svelte';
  import Map from '$lib/Map.svelte';
	import type { Geocode, Forecast } from './types';

  interface Props {
    geocode: Geocode;
    forecast: Forecast;
  }

  let { forecast, geocode }: Props = $props();
  console.log(geocode);
</script>

<Current currently={ forecast.currently } daily={ forecast.daily } />
<Hourly hourly={ forecast.hourly } />
<Map location={ geocode.geometry.location } />
<Daily daily={ forecast.daily } hourly={ forecast.hourly } />