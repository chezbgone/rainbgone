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
	// 25th point). The backend's `hourlyFromMidnight` is a flat series anchored at today's
	// midnight (today backfilled via the Time Machine API), so each day is a contiguous
	// 24-hour slice; the +1 includes the next midnight for the closing tick/label that
	// Stripes needs to draw the day's last hour.
	const start = dayIndex * 24;
	const hourly = forecast.hourlyFromMidnight.slice(start, start + 25);
	// The final covered day has no next hour in the series, so the slice above comes up
	// short of 25 entries and lacks that closing tick, which would make Stripes drop the
	// day's last real hour. Synthesize the boundary from the last hour instead.
	if (hourly.length > 0 && hourly.length < 25) {
		const last = hourly[hourly.length - 1];
		hourly.push({ ...last, time: last.time + 3600 });
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
