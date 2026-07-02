<script lang="ts">
	import { formatHourMinute } from '$lib/common/time';
	import { number } from '$lib/common/format';
	import type { DailyDatum } from '$lib/Forecast/types';

	interface Props {
		day: DailyDatum;
		timezone: string;
	}

	let { day, timezone }: Props = $props();

	const precipType = $derived(day.precipType);
	const hasPrecip = $derived(precipType !== 'none');
</script>

<div class="flex flex-wrap items-end justify-center gap-x-12 gap-y-4 text-base">
	<!-- High / low temperature -->
	<div class="flex items-end">
		<span class="flex items-baseline gap-1">
			<span class="text-2xl font-medium">{Math.round(day.temperatureMax)}°</span>
			<span class="text-sm">{formatHourMinute(day.temperatureMaxTime, timezone)}</span>
		</span>
		<span class="mx-1 text-xl">→</span>
		<span class="flex items-baseline gap-1">
			<span class="text-2xl font-medium">{Math.round(day.temperatureMin)}°</span>
			<span class="text-sm">{formatHourMinute(day.temperatureMinTime, timezone)}</span>
		</span>
	</div>

	<!-- Sun times -->
	<div class="flex items-center gap-4">
		<span class="flex items-center gap-1">
			<img src="/sunrise.png" width="28" height="30" alt="Sunrise" />
			{formatHourMinute(day.sunriseTime, timezone)}
		</span>
		<span class="flex items-center gap-1">
			<img src="/sunset.png" width="28" height="30" alt="Sunset" />
			{formatHourMinute(day.sunsetTime, timezone)}
		</span>
	</div>

	<!-- Precipitation accumulation -->
	{#if hasPrecip}
		<div class="flex items-baseline gap-1">
			<span class="font-medium capitalize">{precipType}</span>
			<span>{number(day.precipAccumulation, 2)} in</span>
		</div>
	{/if}
</div>
