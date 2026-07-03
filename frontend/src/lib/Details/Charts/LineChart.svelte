<script lang="ts">
	import {
		linearScale,
		niceBounds,
		splinePath,
		CHART_VIEW_WIDTH as width,
		CHART_VIEW_HEIGHT as height,
		CHART_MARGIN_TOP as marginTop,
		CHART_MARGIN_BOTTOM as marginBottom
	} from './geometry';
	import { formatUnixTime } from '$lib/common/time';
	import { number } from '$lib/common/format';
	import type { Forecast } from '$lib/Forecast/types';

	type HourlyDatum = Forecast['hourly']['data'][number];

	interface SeriesConfig {
		value: (h: HourlyDatum) => number;
		color: string;
	}

	interface Props {
		hours: HourlyDatum[];
		label: string;
		series: SeriesConfig[];
		unit: string;
		fromZero: boolean;
		minRange: number;
		max?: number;
	}

	let { hours, label, series, unit, fromZero, minRange, max }: Props = $props();

	// Edge-to-edge: hour 0 at x=0, the last hour at x=width. Matches Stripes' tick spacing
	// (index / (hourCount - 1)) so the scrubber overlay lines up across both components.
	const xScale = $derived(linearScale([0, hours.length - 1], [0, width]));

	const bounds = $derived.by(() =>
		niceBounds(
			hours.flatMap((h) => series.map((s) => s.value(h))),
			{ fromZero, minRange, max }
		)
	);

	const yScale = $derived(
		linearScale([bounds.min, bounds.max], [height - marginBottom, marginTop])
	);

	const paths = $derived(
		series.map((s) => ({
			color: s.color,
			d: splinePath(hours.map((h, i) => ({ x: xScale(i), y: yScale(s.value(h)) })))
		}))
	);

	// Label every other hour, matching the hourly-stripe ticks; skip the closing boundary
	// tick so the last label doesn't crowd the axis.
	const xTicks = $derived(
		hours
			.map((h, i) => ({ i, time: h.time }))
			.filter((t) => t.i % 2 === 0 && t.i !== hours.length - 1)
	);
</script>

{#if hours.length > 1}
	<div>
		<h3 class="mb-1 text-center text-2xl font-medium text-neutral-600">{label}</h3>
		<div class="relative">
			<svg viewBox="0 0 {width} {height}" class="w-full" role="img" aria-label="{label} chart">
				{#each bounds.ticks as tick (tick)}
					<line
						x1={0}
						x2={width}
						y1={yScale(tick)}
						y2={yScale(tick)}
						stroke="#e5e5e5"
						stroke-width="1"
					/>
				{/each}

				{#each xTicks as tick (tick.i)}
					<text
						x={tick.i === 0 ? xScale(tick.i) + 4 : xScale(tick.i)}
						y={height - 8}
						text-anchor={tick.i === 0 ? 'start' : 'middle'}
						class="fill-neutral-500"
						font-size="14"
					>
						{formatUnixTime(tick.time)}
					</text>
				{/each}

				{#each paths as path (path.color)}
					<path
						d={path.d}
						fill="none"
						stroke={path.color}
						stroke-width="5"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
				{/each}
			</svg>

			<!-- Y-axis value labels as an HTML overlay (the plot itself is edge-to-edge, so
			     there's no SVG gutter for them). On wide viewports there's room outside
			     <main> to the left, so labels sit there, off the plot. Below that breakpoint
			     they fall back to overlaying above the gridlines so nothing clips. -->
			{#each bounds.ticks as tick (tick)}
				<span
					class="pointer-events-none absolute left-1 -translate-y-1/2 text-sm whitespace-nowrap text-neutral-500 min-[960px]:right-full min-[960px]:left-auto min-[960px]:mr-2"
					style="top: {(yScale(tick) / height) * 100}%"
				>
					{number(tick)}{unit}
				</span>
			{/each}
		</div>
	</div>
{/if}
