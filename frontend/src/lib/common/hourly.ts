import type { Forecast } from '$lib/Forecast/types';
import { formatDateKey } from './time';

// Include a closing tick for the final hour, without assuming midnight backfill
// succeeded or that a local calendar day contains exactly 24 hours.
export function hoursForDate(
	series: Forecast['hourly']['data'],
	date: string,
	timezone: string
): Forecast['hourly']['data'] {
	let lastIndex = -1;
	const hours = series.filter((hour, i) => {
		const matches = formatDateKey(hour.time, timezone) === date;
		if (matches) lastIndex = i;
		return matches;
	});
	if (hours.length === 0) return hours;
	const last = hours[hours.length - 1];
	const next = series[lastIndex + 1];
	hours.push(next?.time === last.time + 3600 ? next : { ...last, time: last.time + 3600 });
	return hours;
}
