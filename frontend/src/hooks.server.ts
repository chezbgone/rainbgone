import type { Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

const apiProxyTarget = env.API_PROXY_TARGET ?? env.VITE_API_PROXY_TARGET ?? 'http://localhost:8080';

export const handle: Handle = async ({ event, resolve }) => {
	if (!event.url.pathname.startsWith('/api/')) {
		return resolve(event);
	}

	if (!['/api/forecast', '/api/geocode'].includes(event.url.pathname)) {
		return new Response('Not found', { status: 404 });
	}
	if (!['GET', 'HEAD'].includes(event.request.method)) {
		return new Response('Method not allowed', { status: 405, headers: { Allow: 'GET, HEAD' } });
	}

	// Assign the path to a trusted URL: resolving an untrusted //path can replace its host.
	const proxyUrl = new URL(apiProxyTarget);
	proxyUrl.pathname = event.url.pathname.slice('/api'.length);
	proxyUrl.search = event.url.search;
	proxyUrl.hash = '';

	try {
		return await fetch(proxyUrl, {
			method: event.request.method,
			headers: { accept: event.request.headers.get('accept') ?? 'application/json' },
			signal: event.request.signal,
			redirect: 'manual'
		});
	} catch {
		return new Response('Backend unavailable', { status: 502 });
	}
};
