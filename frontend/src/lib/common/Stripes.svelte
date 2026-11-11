<script lang="ts">
	import { formatUnixTime } from './time';
	import type { Forecast } from '../Forecast/types';
	import type { Weather } from './weather';
	import { classifyWeather, weatherInfo } from './weather';

	interface Props {
		hours: Forecast['hourly']['data'];
		timezone: string;
	}

	interface Tick {
		time: number;
		temperature: number;
	}

	interface Stripe {
		width: number;
		weather: Weather;
	}

	let { hours, timezone }: Props = $props();

	const ticks: Tick[] = $derived(hours.map((h) => ({ time: h.time, temperature: h.temperature })));
	const stripes: Stripe[] = $derived.by(() => {
		const acc: Stripe[] = [];
		let currentStripe: Stripe | undefined;
		// The last entry is the closing boundary tick, not a full hour period.
		for (const hour of hours.slice(0, -1)) {
			const weather = classifyWeather(hour.precipIntensity, hour.precipType, hour.cloudCover);
			if (currentStripe?.weather === weather) {
				currentStripe.width++;
			} else {
				currentStripe = { width: 1, weather };
				acc.push(currentStripe);
			}
		}
		return acc;
	});
</script>

{#if hours.length < 2}
	<p class="my-4 text-center text-neutral-500">No hourly timeline is available.</p>
{:else}
	<!-- This max-w-[800px] must stay in sync with <main>'s own max-w-[800px] on the details page:
     the scrubber line and hourly charts span main's content box, and the ticks below line up
     with them only because this cap is non-binding (main is already ≤800px). If main ever
     widens past 800px, this cap would clip the stripes and desync the ticks from the charts. -->
	<div class="mx-auto my-4 flex w-full max-w-[800px] flex-col text-center">
		<div class="flex h-10 w-full">
			{#each stripes as stripe}
				{@const weather = weatherInfo[stripe.weather]}
				<div
					style="flex-grow: {stripe.width};
            flex-basis: 0;
            background-color: {weather.color};"
					class="flex h-full items-center justify-center text-sm first:rounded-l last:rounded-r"
					class:text-neutral-800={weather.labelInfo.darkText}
					class:text-shadow-[1px_1px_0_rgba(255,255,255,0.2)]={weather.labelInfo.darkText}
					class:text-white={!weather.labelInfo.darkText}
					class:text-shadow-[1px_1px_0_rgba(0,0,0,0.2)]={!weather.labelInfo.darkText}
				>
					{#if stripe.width >= weather.labelInfo.requiredWidth}
						{weather.labelInfo.text}
					{/if}
				</div>
			{/each}
		</div>

		<!-- ticks -->
		<div class="mt-1 flex w-full justify-between">
			{#each Array(ticks.length) as _}
				<span class="border-l border-neutral-400 odd:h-2 even:h-1"></span>
			{/each}
		</div>

		<!-- labels -->
		<div class="flex h-4 w-full justify-between pb-12">
			{#each ticks as tick}
				<div class="group relative even:hidden">
					<div class="absolute top-0 text-xs not-group-first:-translate-x-1/2">
						{formatUnixTime(tick.time, timezone)}
					</div>
					<div
						class="absolute top-4 text-lg font-light not-group-first:-translate-x-[calc(50%-2px)]"
					>
						{Math.round(tick.temperature)}°
					</div>
				</div>
			{/each}
		</div>
	</div>
{/if}
