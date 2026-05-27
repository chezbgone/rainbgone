<script lang="ts">
  import Stripes from '$lib/common/Stripes.svelte';
  import { formatUnixTime } from '$lib/common/time';
  import type { PageProps } from './$types';

  let { data }: PageProps = $props();

  const formatDate = (time: number) =>
    new Intl.DateTimeFormat(undefined, {
      timeZone: data.timezone,
      weekday: 'long',
      month: 'long',
      day: 'numeric',
      year: 'numeric',
    }).format(time * 1000);

  const formatHourMinute = (time: number) =>
    new Intl.DateTimeFormat(undefined, {
      timeZone: data.timezone,
      hour: 'numeric',
      minute: '2-digit',
    }).format(time * 1000);

  const number = (value: number | null | undefined, digits = 0) =>
    value === null || value === undefined ? 'N/A' : value.toFixed(digits);

  const percent = (value: number | null | undefined) =>
    value === null || value === undefined ? 'N/A' : `${Math.round(value * 100)}%`;

  const bearing = (value: number | null | undefined) =>
    value === null || value === undefined ? 'N/A' : `${Math.round(value)}°`;

  const precipType = data.day.precipType === 'none' ? 'None' : data.day.precipType;

  const detailRows = [
    ['Wind speed', `${number(data.day.windSpeed, 1)} mph`],
    ['Wind gust', `${number(data.day.windGust, 1)} mph at ${formatHourMinute(data.day.windGustTime)}`],
    ['Wind bearing', bearing(data.day.windBearing)],
    ['Humidity', percent(data.day.humidity)],
    ['Pressure', `${number(data.day.pressure, 1)} mb`],
    ['Dew point', `${number(data.day.dewPoint)}°`],
    ['UV index', `${number(data.day.uvIndex)} at ${formatHourMinute(data.day.uvIndexTime)}`],
    ['Visibility', `${number(data.day.visibility, 1)} mi`],
    ['Cloud cover', percent(data.day.cloudCover)],
    ['Moon phase', percent(data.day.moonPhase)],
  ];
</script>

<main class="mx-auto max-w-5xl px-4 py-8 text-neutral-800">
  <a href={data.links.back} class="text-blue-600 hover:text-blue-800">Back to forecast</a>

  <header class="mt-6 flex flex-col gap-4 border-b border-neutral-300 pb-6 sm:flex-row sm:items-center sm:justify-between">
    <div>
      <div class="text-sm uppercase tracking-wide text-neutral-500">{data.formattedAddress ?? `${data.lat}, ${data.lng}`}</div>
      <h1 class="mt-1 text-4xl font-light">{formatDate(data.day.time)}</h1>
      <p class="mt-2 text-xl font-light">{data.day.summary}</p>
    </div>
    <img src={`/weather-icons/${data.day.icon}.png`} alt={data.day.icon} class="h-20 w-20 self-start sm:self-center" />
  </header>

  <section class="grid gap-6 py-6 md:grid-cols-2">
    <div>
      <h2 class="text-xl font-light">Temperature</h2>
      <dl class="mt-3 grid gap-3">
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Minimum</dt>
          <dd>{Math.round(data.day.temperatureMin)}° at {formatHourMinute(data.day.temperatureMinTime)}</dd>
        </div>
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Maximum</dt>
          <dd>{Math.round(data.day.temperatureMax)}° at {formatHourMinute(data.day.temperatureMaxTime)}</dd>
        </div>
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Feels-like minimum</dt>
          <dd>{Math.round(data.day.apparentTemperatureMin)}° at {formatHourMinute(data.day.apparentTemperatureMinTime)}</dd>
        </div>
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Feels-like maximum</dt>
          <dd>{Math.round(data.day.apparentTemperatureMax)}° at {formatHourMinute(data.day.apparentTemperatureMaxTime)}</dd>
        </div>
      </dl>
    </div>

    <div>
      <h2 class="text-xl font-light">Sun and Precipitation</h2>
      <dl class="mt-3 grid gap-3">
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Sunrise</dt>
          <dd>{formatHourMinute(data.day.sunriseTime)}</dd>
        </div>
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Sunset</dt>
          <dd>{formatHourMinute(data.day.sunsetTime)}</dd>
        </div>
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Type</dt>
          <dd class="capitalize">{precipType}</dd>
        </div>
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Probability</dt>
          <dd>{percent(data.day.precipProbability)}</dd>
        </div>
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Intensity</dt>
          <dd>{number(data.day.precipIntensity, 2)} in/hr</dd>
        </div>
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Max intensity</dt>
          <dd>{number(data.day.precipIntensityMax, 2)} in/hr at {formatHourMinute(data.day.precipIntensityMaxTime)}</dd>
        </div>
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>Accumulation</dt>
          <dd>{number(data.day.precipAccumulation, 2)} in</dd>
        </div>
      </dl>
    </div>
  </section>

  <section class="border-t border-neutral-300 py-6">
    <h2 class="text-xl font-light">24-hour timeline</h2>
    {#if data.hourly.length > 0}
      <Stripes hours={data.hourly} />
    {:else}
      <p class="mt-3 text-neutral-500">No hourly forecast data is available for this day.</p>
    {/if}
  </section>

  <section class="border-t border-neutral-300 py-6">
    <h2 class="text-xl font-light">Daily details</h2>
    <dl class="mt-3 grid gap-x-8 gap-y-3 sm:grid-cols-2">
      {#each detailRows as [label, value]}
        <div class="flex justify-between border-b border-neutral-200 pb-2">
          <dt>{label}</dt>
          <dd>{value}</dd>
        </div>
      {/each}
    </dl>
  </section>

  <nav class="flex items-center justify-between border-t border-neutral-300 pt-6">
    {#if data.links.previous}
      <a href={data.links.previous} class="text-blue-600 hover:text-blue-800">Previous day</a>
    {:else}
      <span></span>
    {/if}
    {#if data.links.next}
      <a href={data.links.next} class="text-blue-600 hover:text-blue-800">Next day</a>
    {/if}
  </nav>
</main>
