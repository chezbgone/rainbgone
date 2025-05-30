import type { PageServerLoad } from './$types';

export const load = (async ({ fetch, params }) => {
  const lat = params.lat
  const lng = params.lng
  let geocode = await fetch('/api/geocode?' + new URLSearchParams({
    address: `${lat},${lng}`
  }));
  let forecast = await fetch('/api/forecast?' + new URLSearchParams({
    lat: lat,
    lng: lng
  }));

  return { forecast: forecast.ok && await forecast.json(), geocode: geocode.ok && await geocode.json() };
}) satisfies PageServerLoad;