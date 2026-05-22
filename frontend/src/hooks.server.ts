import type { Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

const apiProxyTarget = env.API_PROXY_TARGET ?? env.VITE_API_PROXY_TARGET ?? 'http://localhost:8080';

export const handle: Handle = async ({ event, resolve }) => {
  if (!event.url.pathname.startsWith('/api/')) {
    return resolve(event);
  }

  const proxyUrl = new URL(event.url.pathname.replace(/^\/api/, '') + event.url.search, apiProxyTarget);
  const headers = new Headers(event.request.headers);
  headers.delete('host');

  return fetch(proxyUrl, {
    method: event.request.method,
    headers,
    body: ['GET', 'HEAD'].includes(event.request.method) ? undefined : await event.request.arrayBuffer(),
    redirect: 'manual',
  });
};
