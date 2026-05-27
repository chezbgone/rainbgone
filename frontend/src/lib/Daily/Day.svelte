<script lang="ts">
  import Stripes from "$lib/common/Stripes.svelte";
  import { formatDateKey, formatUnixTime } from "$lib/common/time";
  import type { Forecast } from "../Forecast/types";

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

  const getDayName = (dt: number) => {
    if (today) return 'Today';
    return new Date(dt * 1000).toLocaleDateString(undefined, { weekday: 'short' });
  };
  const formatHourMinute = (unixTime: number) =>
    new Date(unixTime * 1000).toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' });

  const left = (daily.temperatureMin - minTemp) / (maxTemp - minTemp) * 100;
  const width = (daily.temperatureMax - daily.temperatureMin) / (maxTemp - minTemp) * 100;
  const right = (daily.temperatureMax - minTemp) / (maxTemp - minTemp) * 100;

  let opened = $state(false);
  const detailsHref = $derived(`/details/${lat},${lng}/${formatDateKey(daily.time, timezone)}`);
</script>

<details bind:open={opened}>
  <summary class="group flex items-center px-4 py-2 hover:bg-gray-50 hover:cursor-pointer">
    <img src={`/weather-icons/${daily.icon}.png`} alt={daily.icon} class="h-8" />
    <span class="w-16 mx-2">{getDayName(daily.time)}</span>
    <span class="relative flex items-center w-lg h-4 mx-16">
      <span
        class="absolute text-right ml-[-30px]"
        style="left: {left}%"
      >{Math.round(daily.temperatureMin)}°</span>
      <span
        class="absolute left-0 h-4 rounded-full bg-gray-700"
        style="
          margin-left: {left}%;
          width: {width}%;
        "
      ></span>
      <span
        class="absolute text-left ml-[5px]"
        style="left: {right}%"
      >{Math.round(daily.temperatureMax)}°</span>
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
      <circle cx="12" cy="12" r="12" mask={opened ? "url(#minus-mask)" : "url(#plus-mask)"} />
    </svg>
  </summary>
  <div class="py-2 border-b border-neutral-400">
    <div class="text-xl font-light">{daily.summary}</div>
    <div class="inline-flex gap-16 my-2 items-center text-neutral-800">
      <div>
        <span class="text-2xl">{Math.round(daily.temperatureMin)}°</span>
        <span class="font-light">{formatUnixTime(daily.temperatureMinTime)}</span>
        <span class="mx-2 text-2xl font-light">&#8594;</span>
        <span class="text-2xl">{Math.round(daily.temperatureMax)}°</span>
        <span class="font-light">{formatUnixTime(daily.temperatureMaxTime)}</span>
      </div>
      <div class="font-light">
        <span>
          <img src="sunrise.png" alt="Sunrise icon" class="inline w-6 h-6" />
          <span>{formatHourMinute(daily.sunriseTime)}</span>
        </span>
        <span>
          <img src="sunset.png" alt="Sunset icon" class="inline w-6 h-6 ml-4" />
          <span>{formatHourMinute(daily.sunsetTime)}</span>
        </span>
      </div>
      <div>
        <span class="mr-1">Rain</span>
        <span class="font-light">{daily.precipIntensity.toFixed(2)}</span>
        <span class="font-light">in</span>
      </div>
    </div>
    <Stripes hours={hourly} />
    <a href={detailsHref} class="inline-block px-4 py-2 mb-4 rounded-sm bg-blue-500 hover:bg-blue-600 text-white uppercase">
      more details
    </a>
  </div>
</details>
