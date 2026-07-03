<script lang="ts">
	import { clampXFraction, hourFromXFraction } from './Charts/geometry';
	import { Spring } from 'svelte/motion';
	import { untrack, type Snippet } from 'svelte';
	import type { Forecast } from '$lib/Forecast/types';

	type HourlyDatum = Forecast['hourly']['data'][number];

	interface Props {
		hours: HourlyDatum[];
		children: Snippet<[number]>;
	}

	let { hours, children }: Props = $props();

	// Scrubber: a single overlay div positioned over the whole Stripes + InstantDetails + chart
	// stack, driven purely by pointer position + the shared hour<->fraction geometry
	// (src/lib/Details/Charts/geometry.ts). None of the wrapped components know about it. It's
	// always visible (starting centered on the day), tracks the pointer continuously (not
	// snapped to hour columns), and only moves while click/touch-dragging, not on plain hover.
	// The rendered position trails the pointer with a critically damped spring (drag
	// "friction") rather than snapping straight to it. Because the charts and Stripes' ticks
	// both place hour `i` at fraction `i / (hours.length - 1)`, one straight line stays
	// correct over both.
	let targetFraction = $state(0.5);
	// precision is in the same units as the spring's value — since position lives in [0,1]
	// fraction space (not pixels), the default precision (0.01, tuned for pixel-scale values)
	// is ~1% of the container width and causes a visible final-frame jump. Lower it so the
	// spring's terminal step is sub-pixel instead.
	const position = new Spring(0.5, { stiffness: 0.15, damping: 1, precision: 0.0001 });
	let dragging = $state(false);

	// Follow the pointer's target fraction with spring friction.
	$effect(() => {
		position.target = targetFraction;
	});

	// Recenter only when the day actually changes (new hourly series), snapping instantly.
	// `dayKey` is the effect's only dependency and is stable during a drag, so this can't fire
	// mid-drag. The untrack is belt-and-suspenders: the block only writes (writes don't create
	// dependencies anyway), but it guarantees no future read in here accidentally subscribes.
	const dayKey = $derived(hours[0]?.time);
	$effect(() => {
		dayKey;
		untrack(() => {
			targetFraction = 0.5;
			position.set(0.5, { instant: true });
		});
	});

	function updateFromPointer(e: PointerEvent) {
		const rect = (e.currentTarget as HTMLDivElement).getBoundingClientRect();
		const fraction = (e.clientX - rect.left) / rect.width;
		targetFraction = clampXFraction(fraction, hours.length);
	}

	function handlePointerDown(e: PointerEvent) {
		if (hours.length < 2) return;
		dragging = true;
		(e.currentTarget as HTMLDivElement).setPointerCapture(e.pointerId);
		updateFromPointer(e);
	}

	function handlePointerMove(e: PointerEvent) {
		if (!dragging) return;
		updateFromPointer(e);
	}

	function handlePointerUp(e: PointerEvent) {
		dragging = false;
		(e.currentTarget as HTMLDivElement).releasePointerCapture(e.pointerId);
	}

	const scrubberLeftPercent = $derived(position.current * 100);

	// The fractional hour the line is currently drawn over, derived from the spring's rendered
	// position (not the raw pointer target) so it stays visually tied to where the line
	// actually is. Un-rounded so children can interpolate between the bracketing hours.
	const scrubbedHour = $derived(hourFromXFraction(position.current, hours.length));
</script>

<!-- Known a11y gap: this drag region is pointer-only. There's no keyboard/slider semantics
     yet (role="slider", aria-valuemin/max/now, ArrowLeft/Right to step by an hour), so keyboard
     and screen-reader users can't move the scrubber — InstantDetails stays at the centered hour
     for them. TODO: add slider semantics + arrow-key handling. -->
<div
	class="relative touch-none select-none"
	onpointerdown={handlePointerDown}
	onpointermove={handlePointerMove}
	onpointerup={handlePointerUp}
	onpointercancel={handlePointerUp}
>
	<!-- Translucent band centered on the scrubber position, with the 1px line nested inside.
	     A same-height band (inset-y-0) glows only horizontally — unlike a box-shadow spread,
	     which would bleed above and below the line. Position is set once, on the wrapper. -->
	<div
		class="pointer-events-none absolute inset-y-0 flex w-8 -translate-x-1/2 justify-center bg-red-500/5"
		style="left: {scrubberLeftPercent}%"
	>
		<div class="w-px bg-red-400"></div>
	</div>

	{@render children(scrubbedHour)}
</div>
