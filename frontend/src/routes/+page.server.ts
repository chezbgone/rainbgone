import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load = (async () => {
	redirect(303, `/47.6062,-122.3321`); // Seattle
}) satisfies PageServerLoad;
