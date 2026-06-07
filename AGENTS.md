# Agent Instructions

## Project Overview

rainbgone is a small weather app with two runtimes:

- Go backend in `server/`, started by `main.go`.
- SvelteKit frontend in `frontend/`.

Read `docs/architecture.md` before changing forecast, geocoding, map, tile proxy, runtime config, Docker, or API behavior.

## Common Commands

Backend:

```sh
go run .
go test ./...
```

Frontend:

```sh
cd frontend
npm run dev
npm run check
npm run lint
npm run build
```

Docker development:

```sh
docker compose -f docker-compose.dev.yml up --build
```

## Runtime Configuration

The Go backend reads `.env` with `godotenv.Read()` at startup. Current backend config keys are:

- `PIRATE_WEATHER_KEY`
- `MAPTILER_KEY`

The SvelteKit server-side API proxy uses:

- `API_PROXY_TARGET`
- `VITE_API_PROXY_TARGET`

Both proxy variables default to `http://localhost:8080` when unset.

Do not add, rename, or remove configuration keys without updating `server/config.go`, Docker configuration, and `docs/architecture.md`.

## Secrets

Never print, quote, commit, or embed secrets from `.env` or provider-generated files. Treat local provider exports such as `style.json` as unsafe to commit until they have been checked for embedded keys.

If a task requires checking whether an env var exists, report only whether it is set, not its value.

## Backend Notes

Routes are registered in `server/server.go`:

- `/geocode`
- `/forecast`
- `/map/background-tiles/`

Follow the existing `net/http` style unless a broader backend rewrite is explicitly requested.

The backend uses simple in-memory caches in `server/cache.go`. Cache entries are process-local and are lost on restart.

## Frontend Notes

SvelteKit server loads call backend routes through `/api/...`; the proxy is implemented in `frontend/src/hooks.server.ts` and configured for Vite dev in `frontend/vite.config.ts`.

The map UI lives in `frontend/src/lib/Map.svelte` and uses OpenLayers plus `ol-mapbox-style`.

When changing API routes, update every frontend call site and the architecture document in the same change.

## Map And Tile Changes

Do not assume external background tile IDs or future weather layer names are stable. Verify MapTiler or other provider URLs before depending on hard-coded IDs.

Keep background/base map tiles and weather/radar overlays conceptually separate. A missing weather overlay should not remove base geography unless the design explicitly requires that behavior.

## Validation Expectations

Run the narrowest relevant check for the files changed:

- Go backend changes: `go test ./...`
- Frontend TypeScript/Svelte changes: `cd frontend && npm run check`
- Formatting-only frontend changes: `cd frontend && npm run lint`
- Docker/runtime changes: validate the compose file or build the touched service when practical

If the user asks to skip tests, state that validation was skipped.
