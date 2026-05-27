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
  const dayIndex = days.findIndex(day => formatDateKey(day.time, forecast.timezone) === date);

  if (dayIndex === -1) {
    error(404, 'Forecast day not found');
  }

  const day = days[dayIndex];
  const hourly = forecast.hourly.data.filter(hour => formatDateKey(hour.time, forecast.timezone) === date);
  const previousDay = days[dayIndex - 1];
  const nextDay = days[dayIndex + 1];

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
      next: nextDay
        ? `/details/${lat},${lng}/${formatDateKey(nextDay.time, forecast.timezone)}`
        : null,
    },
  };
}) satisfies PageServerLoad;
