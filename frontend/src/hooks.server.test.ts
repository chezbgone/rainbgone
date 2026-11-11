import { describe, expect, it, vi } from 'vitest';
import type { RequestEvent } from '@sveltejs/kit';
import { handle } from './hooks.server';

vi.mock('$env/dynamic/private', () => ({ env: { API_PROXY_TARGET: 'http://backend:8080' } }));

function proxy(path: string, method = 'GET') {
	const url = new URL(path, 'http://app.example');
	const event = {
		url,
		request: new Request(url, {
			method,
			headers: { cookie: 'dummy=secret', authorization: 'Bearer dummy', accept: 'application/json' }
		})
	} as RequestEvent;
	const resolve = vi.fn(async () => new Response('page'));
	return { event, resolve };
}

describe('API proxy', () => {
	it.each([
		'/api//127.0.0.1:9000/private',
		'/api///attacker.example/private',
		'/api/%2f%2fattacker.example/private',
		'/api/unrecognized'
	])('rejects unrecognized or host-changing paths: %s', async (path) => {
		const fetch = vi.fn(async () => new Response('must not be reached'));
		vi.stubGlobal('fetch', fetch);
		const response = await handle(proxy(path));
		expect(response.status).toBe(404);
		expect(fetch).not.toHaveBeenCalled();
	});

	it.each(['forecast?lat=47&lng=-122', 'geocode?address=Seattle'])(
		'forwards %s to the configured backend without private headers',
		async (path) => {
			const fetch = vi.fn(async () => new Response('fixture'));
			vi.stubGlobal('fetch', fetch);
			const response = await handle(proxy(`/api/${path}`));
			expect(await response.text()).toBe('fixture');
			const [url, options] = fetch.mock.calls[0] as unknown as [URL, RequestInit];
			expect(url.href).toBe(`http://backend:8080/${path}`);
			const headers = new Headers(options.headers);
			expect(headers.get('cookie')).toBeNull();
			expect(headers.get('authorization')).toBeNull();
			expect(headers.get('accept')).toBe('application/json');
			expect(options.redirect).toBe('manual');
		}
	);

	it('passes ordinary pages to SvelteKit', async () => {
		const args = proxy('/47,-122');
		const response = await handle(args);
		expect(await response.text()).toBe('page');
		expect(args.resolve).toHaveBeenCalledWith(args.event);
	});

	it('rejects unsupported methods without contacting the backend', async () => {
		const fetch = vi.fn();
		vi.stubGlobal('fetch', fetch);
		const response = await handle(proxy('/api/forecast', 'POST'));
		expect(response.status).toBe(405);
		expect(response.headers.get('allow')).toBe('GET, HEAD');
		expect(fetch).not.toHaveBeenCalled();
	});

	it('returns a generic error when the backend cannot be reached', async () => {
		vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('private connection details')));
		const response = await handle(proxy('/api/forecast'));
		expect(response.status).toBe(502);
		expect(await response.text()).toBe('Backend unavailable');
	});
});
