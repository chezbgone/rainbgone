import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { formatDateKey } from '$lib/common/time';
import type { Forecast } from '$lib/Forecast/types';

const datePattern = /^\d{4}-\d{2}-\d{2}$/;

export const load = (async ({ fetch, params }) => {
	const { lat, lng, date } = params;

	if (!datePattern.test(date)) {
		error(404, 'Forecast day not found');
	}

	const forecastResponse = await fetch('/api/forecast?' + new URLSearchParams({ lat, lng }));

	if (!forecastResponse.ok) {
		error(forecastResponse.status, 'Error fetching weather data');
	}

	const forecast = (await forecastResponse.json()) as Forecast & { formatted_address?: string };
	const days = forecast.daily.data;
	const dayIndex = days.findIndex((day) => formatDateKey(day.time, forecast.timezone) === date);

	if (dayIndex === -1) {
		error(404, 'Forecast day not found');
	}

	const day = days[dayIndex];
	// The day's hourly entries, midnight to midnight (plus the closing midnight as the
	// closing tick Stripes needs to draw the day's last hour). Selected by matching each
	// entry's local calendar date against `date`, rather than assuming `hourlyFromMidnight`
	// is a clean N*24-long series starting at midnight: the backend's Time Machine backfill
	// is best-effort (see server/forecast.go), so a degraded/failing backfill can leave the
	// series short or starting later than midnight. Filtering by local date means each day
	// still renders whatever hours it actually has — worst case today starts later than
	// midnight — instead of every day sliding off its boundary in lockstep.
	const series = forecast.hourlyFromMidnight;
	let lastIndex = -1;
	const hourly = series.filter((h, i) => {
		const matches = formatDateKey(h.time, forecast.timezone) === date;
		if (matches) lastIndex = i;
		return matches;
	});
	// Close the boundary with the real next hour when the series has one (e.g. today, or
	// any day that isn't the last one covered); otherwise synthesize it from the last hour,
	// same as before.
	if (hourly.length > 0) {
		const next = series[lastIndex + 1];
		if (next) {
			hourly.push(next);
		} else {
			const last = hourly[hourly.length - 1];
			hourly.push({ ...last, time: last.time + 3600 });
		}
	}
	const previousDay = days[dayIndex - 1];
	const nextDay = days[dayIndex + 1];

	const weekday = (time: number) =>
		new Intl.DateTimeFormat(undefined, {
			timeZone: forecast.timezone,
			weekday: 'long'
		}).format(time * 1000);

	return {
		lat: Number(lat),
		lng: Number(lng),
		date,
		timezone: forecast.timezone,
		formattedAddress: forecast.formatted_address,
		day,
		hourly,
		links: {
			back: `/${lat},${lng}`,
			previous: previousDay
				? `/details/${lat},${lng}/${formatDateKey(previousDay.time, forecast.timezone)}`
				: null,
			previousLabel: previousDay ? weekday(previousDay.time) : null,
			next: nextDay
				? `/details/${lat},${lng}/${formatDateKey(nextDay.time, forecast.timezone)}`
				: null,
			nextLabel: nextDay ? weekday(nextDay.time) : null
		}
	};
}) satisfies PageServerLoad;
