import type { PageServerLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load = (async ({ fetch, params }) => {
  const lat = params.lat
  const lng = params.lng
  let forecast = await fetch('/api/forecast?' + new URLSearchParams({
    lat: lat,
    lng: lng
  }));

  if (!forecast.ok) {
    error(forecast.status, 'Error fetching weather data');
  }

  const forecastJson = await forecast.json();
  const latNumber = Number(lat);
  const lngNumber = Number(lng);

  return {
    forecast: forecastJson,
    geocode: {
      formatted_address: forecastJson.formatted_address,
      geometry: {
        location: {
          lat: latNumber,
          lng: lngNumber
        },
        viewport: {
          northeast: {
            lat: latNumber,
            lng: lngNumber
          },
          southwest: {
            lat: latNumber,
            lng: lngNumber
          }
        }
      }
    }
  };
}) satisfies PageServerLoad;
