import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load = (async ({ fetch, params }) => {
	const address = params.address;
	const geocode = await fetch(
		'/api/geocode?' +
			new URLSearchParams({
				address: address
			})
	);
	if (!geocode.ok) {
		return { forecast: null, geocode: null };
	}
	const geocodeJson = await geocode.json();
	const lat = geocodeJson.geometry.location.lat;
	const lng = geocodeJson.geometry.location.lng;

	redirect(303, `/${lat},${lng}`);
}) satisfies PageServerLoad;
