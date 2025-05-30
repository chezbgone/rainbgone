import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load = (async ({ fetch, params }) => {
  redirect(303, `/seattle`);
}) satisfies PageServerLoad;