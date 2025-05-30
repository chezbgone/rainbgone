<script lang="ts">
  import { formatUnixTime } from "./time";
	import type { Forecast } from "../Forecast/types";
	import type { Weather } from "./weather";
	import { classifyWeather, weatherInfo } from "./weather";

  interface Props {
    hours: Forecast['hourly']['data'];
  }

  interface Tick {
    time: number;
    temperature: number;
  }

  interface Stripe {
    width: number;
    weather: Weather;
  }

  let { hours }: Props = $props();

  const ticks: Tick[] = $derived(hours.map(h => ({ time: h.time, temperature: h.temperature })));
  const stripes: Stripe[] = $derived.by(() => {
    const acc: Stripe[] = [];
    let current_stripe: Stripe = {
      width: 0,
      // @ts-ignore: null will be overwritten
      weather: null,
    }
    for (let hour of hours) {
      let weather = classifyWeather(hour.precipIntensity, hour.precipType, hour.cloudCover)
      if (current_stripe.weather === null || weather !== current_stripe.weather) {
        if (current_stripe.weather !== null) {
          acc.push(current_stripe);
        }
        current_stripe = {
          width: 1,
          weather: weather,
        }
      } else {
        current_stripe.width++;
      }
    }
    acc.push(current_stripe);
    return acc;
  });
</script>

<div class="text-center my-4">
  <div class="inline-flex flex-col">
    <div class="flex w-[800px] h-10">
      {#each stripes as stripe}
        {@const weather = weatherInfo[stripe.weather]}
        <div
          style="width:{stripe.width * 33.33}px;
            background-color: {weather.color};"
          class="h-full flex justify-center items-center first:rounded-l last:rounded-r text-sm"
          class:text-neutral-800={weather.labelInfo.darkText}
          class:text-shadow-[1px_1px_0_rgba(255,255,255,0.5)]={weather.labelInfo.darkText}
          class:text-white={!weather.labelInfo.darkText}
          class:text-shadow-[1px_1px_0_rgba(0,0,0,0.5)]={!weather.labelInfo.darkText}>
          {#if stripe.width >= weather.labelInfo.requiredWidth}
            {weather.labelInfo.text}
          {/if}
        </div>
      {/each}
    </div>

    <!-- ticks -->
    <div class="w-[800px] mt-1 flex justify-between">
      {#each Array(25) as _}
        <span class="border-l border-neutral-400 even:h-1 odd:h-2"></span>
      {/each}
    </div>

    <!-- labels -->
    <div class="w-[800px] h-4 flex justify-between pb-12">
      {#each ticks as tick}
        <div class="group relative even:hidden">
          <div class="absolute top-0 not-group-first:-translate-x-1/2 text-xs">
            {formatUnixTime(tick.time)}
          </div>
          <div class="absolute top-4 not-group-first:-translate-x-[calc(50%-2px)] text-lg font-light">
            {Math.round(tick.temperature)}°
          </div>
        </div>
      {/each}
      <div></div>
    </div>
  </div>
</div>