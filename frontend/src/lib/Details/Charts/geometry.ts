// Pure geometry helpers for the hourly line charts. No Svelte/DOM dependency so this can
// be unit-tested in isolation and is reused by the synced scrubber overlay.

export interface Point {
	x: number;
	y: number;
}

export interface Scale {
	(value: number): number;
	invert(pixel: number): number;
}

/** Maps a numeric domain onto a pixel range. `scale.invert` does the reverse (pixel →
 *  value), used by the scrubber to turn a pointer position into a fractional hour. */
export function linearScale(domain: [number, number], range: [number, number]): Scale {
	const [d0, d1] = domain;
	const [r0, r1] = range;
	const span = d1 - d0 || 1;
	const scale = ((value: number) => r0 + ((value - d0) / span) * (r1 - r0)) as Scale;
	scale.invert = (pixel: number) => d0 + ((pixel - r0) / (r1 - r0 || 1)) * span;
	return scale;
}

export interface Bounds {
	min: number;
	max: number;
	ticks: number[];
}

interface NiceBoundsOptions {
	fromZero: boolean;
	minRange: number;
	max?: number;
}

/** Computes an axis range and "nice" tick values for a set of data values, honoring a
 *  minimum visible span and an optional fixed ceiling. */
export function niceBounds(values: number[], opts: NiceBoundsOptions): Bounds {
	// Empty input would make the spread-based min/max below ±Infinity and cascade to NaN.
	// Callers with no data don't render, but bail explicitly so correctness never hinges on that.
	if (values.length === 0) {
		return { min: 0, max: 1, ticks: [0, 1] };
	}

	const dataMin = opts.fromZero ? 0 : Math.min(...values);
	const dataMax = Math.max(...values, opts.max ?? -Infinity);

	let min = dataMin;
	let max = dataMax;
	if (max - min < opts.minRange) {
		if (opts.fromZero) {
			max = min + opts.minRange;
		} else {
			const center = (min + max) / 2;
			min = center - opts.minRange / 2;
			max = center + opts.minRange / 2;
		}
	}

	const ticks = niceTicks(min, max, 5);
	return {
		min: Math.min(min, ticks[0]),
		max: Math.max(max, ticks[ticks.length - 1]),
		ticks
	};
}

// Classic "nice numbers for graph labels" algorithm (Heckbert, Graphics Gems I) — picks
// tick spacing from {1, 2, 5} x 10^n so axis labels land on round values.
function niceTicks(min: number, max: number, targetCount: number): number[] {
	if (min === max) {
		min -= 1;
		max += 1;
	}
	const span = niceNumber(max - min, false);
	const step = niceNumber(span / (targetCount - 1), true);
	const niceMin = Math.floor(min / step) * step;
	const niceMax = Math.ceil(max / step) * step;

	const ticks: number[] = [];
	for (let v = niceMin; v <= niceMax + step / 2; v += step) {
		ticks.push(Math.round(v * 1e6) / 1e6);
	}
	return ticks;
}

function niceNumber(range: number, round: boolean): number {
	const exponent = Math.floor(Math.log10(range));
	const fraction = range / 10 ** exponent;

	let niceFraction: number;
	if (round) {
		if (fraction < 1.5) niceFraction = 1;
		else if (fraction < 3) niceFraction = 2;
		else if (fraction < 7) niceFraction = 5;
		else niceFraction = 10;
	} else {
		if (fraction <= 1) niceFraction = 1;
		else if (fraction <= 2) niceFraction = 2;
		else if (fraction <= 5) niceFraction = 5;
		else niceFraction = 10;
	}
	return niceFraction * 10 ** exponent;
}

// Shared coordinate space for every hourly chart's SVG viewBox. Charts render edge-to-edge
// (no horizontal margin) so hour `i` lands at the same fraction `i / (hourCount - 1)` as the
// Stripes timeline's `justify-between` ticks. The scrubber overlay (an absolutely-positioned
// div over the whole Stripes+InstantDetails+chart stack, not part of any individual component)
// uses these same constants to convert a pointer position into an hour index or a CSS `left`
// percentage, so it never needs to reach into a chart's internals.
export const CHART_VIEW_WIDTH = 900;
export const CHART_VIEW_HEIGHT = 250;
export const CHART_MARGIN_TOP = 12;
export const CHART_MARGIN_BOTTOM = 28;

function hourXScale(hourCount: number): Scale {
	return linearScale([0, hourCount - 1], [0, CHART_VIEW_WIDTH]);
}

/** Fraction (0-1) of a chart's rendered width where a given hour index falls. Every chart
 *  shares this viewBox and is rendered at the same CSS width, so the fraction doubles as a
 *  `left` percentage for an overlay positioned over the whole chart stack. Equal to
 *  `index / (hourCount - 1)`, matching Stripes' edge-to-edge tick spacing. */
export function hourXFraction(index: number, hourCount: number): number {
	return hourXScale(hourCount)(index) / CHART_VIEW_WIDTH;
}

/** Continuous (un-rounded, clamped) fractional hour position for a pointer fraction —
 *  the position between hour columns, for interpolating values between hours. */
export function hourFromXFraction(fraction: number, hourCount: number): number {
	const raw = hourXScale(hourCount).invert(fraction * CHART_VIEW_WIDTH);
	return Math.min(hourCount - 1, Math.max(0, raw));
}

/** Clamps a raw pointer fraction to the chart's plotted range (the first/last hour's
 *  position), for a scrubber that tracks the pointer continuously instead of snapping to
 *  hour columns. */
export function clampXFraction(fraction: number, hourCount: number): number {
	const min = hourXFraction(0, hourCount);
	const max = hourXFraction(hourCount - 1, hourCount);
	return Math.min(max, Math.max(min, fraction));
}

/** Builds a smooth SVG path (Catmull-Rom converted to cubic Béziers, tension 1/6) that
 *  passes through every point, approximating a spline curve. */
export function splinePath(points: Point[]): string {
	if (points.length === 0) return '';
	if (points.length === 1) return `M ${points[0].x} ${points[0].y}`;

	let d = `M ${points[0].x} ${points[0].y}`;
	for (let i = 0; i < points.length - 1; i++) {
		const p0 = points[i - 1] ?? points[i];
		const p1 = points[i];
		const p2 = points[i + 1];
		const p3 = points[i + 2] ?? p2;

		const c1x = p1.x + (p2.x - p0.x) / 6;
		const c1y = p1.y + (p2.y - p0.y) / 6;
		const c2x = p2.x - (p3.x - p1.x) / 6;
		const c2y = p2.y - (p3.y - p1.y) / 6;

		d += ` C ${c1x} ${c1y}, ${c2x} ${c2y}, ${p2.x} ${p2.y}`;
	}
	return d;
}
