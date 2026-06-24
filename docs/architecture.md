# Architecture

rainbgone is a weather forecast app with a Go HTTP backend and a SvelteKit frontend. The backend owns external weather, geocoding, and background tile provider calls. The frontend renders forecast pages and maps, and it reaches the backend through an `/api` proxy.

## Runtime Layout

```text
.
├── main.go                  # Starts the Go backend
├── server/                  # Backend routes, provider calls, cache, config
├── frontend/                # SvelteKit app
├── Dockerfile               # Backend production image
├── frontend/Dockerfile      # Frontend dev/build/production image
├── docker-compose.dev.yml   # Local two-service development stack
└── docker-compose.yml       # Production-style two-service stack
```

## Backend

`main.go` calls `server.Start()`. `server.Start()` binds the backend to port `8080` and serves the mux from `server.NewMux()`.

Current backend routes are registered in `server/server.go`:

| Route | Handler | Purpose |
| --- | --- | --- |
| `/geocode` | `GeocodeHandler` | Convert an address into one geocoded result. |
| `/forecast` | `ForecastHandler` | Fetch forecast data and attach a reverse-geocoded formatted address. |

The backend uses the standard library `net/http` package. There is no separate router framework.

## Configuration

`server/config.go` loads `.env` with `godotenv.Read()` during package initialization.

Current backend config fields:

| Env var | Used by | Purpose |
| --- | --- | --- |
| `PIRATE_WEATHER_KEY` | `server/forecast.go` | Pirate Weather API key. |
| `MAPTILER_KEY` | _(currently unused)_ | MapTiler tile API key. Retained for a possible future weather/raster overlay; the basemap now uses keyless OpenFreeMap vector tiles. |

The frontend proxy target is configured separately:

| Env var | Used by | Purpose |
| --- | --- | --- |
| `API_PROXY_TARGET` | `frontend/src/hooks.server.ts` | Runtime SvelteKit server proxy target. |
| `VITE_API_PROXY_TARGET` | `frontend/src/hooks.server.ts`, `frontend/vite.config.ts` | Runtime and Vite dev proxy target. |

Both frontend proxy variables default to `http://localhost:8080`.

`WEATHER_RADAR_BASE_URL` may exist in local environments, but the current code does not read it.

## Forecast Flow

The main forecast page is `frontend/src/routes/[lat=number],[lng=number]/+page.server.ts`.

1. The page receives `lat` and `lng` route params.
2. The SvelteKit server load fetches `/api/forecast?lat=...&lng=...`.
3. `frontend/src/hooks.server.ts` proxies `/api/forecast` to the Go backend `/forecast`.
4. `ForecastHandler` parses `lat` and `lng`.
5. `ForecastHandler` starts the Pirate Weather forecast request and Nominatim reverse geocode request concurrently.
6. The Pirate Weather response body is unmarshaled into a map.
7. The reverse geocode result is added as `formatted_address`.
8. The merged JSON response is returned to the frontend.

Forecast responses are cached in memory by latitude and longitude for one minute.

## Geocoding Flow

Address routes are handled by `frontend/src/routes/[address]/+page.server.ts`.

1. The page receives an address route param.
2. The SvelteKit server load fetches `/api/geocode?address=...`.
3. `frontend/src/hooks.server.ts` proxies `/api/geocode` to the Go backend `/geocode`.
4. `GeocodeHandler` calls Nominatim search and returns the first result.
5. The frontend redirects to `/{lat},{lng}`.

Reverse geocoding happens inside the forecast route so coordinate-based pages can still show a formatted location. Reverse geocode responses are cached in memory by latitude and longitude for 24 hours.

## Details Page Flow

Daily details live under `frontend/src/routes/details/[lat=number],[lng=number]/[date]/`.

The server load fetches the same `/api/forecast` response as the main forecast page, finds the requested daily forecast by local date, filters hourly entries for that date, and returns previous/next detail links when adjacent days exist.

## Map Flow

The map component is `frontend/src/lib/Map.svelte`.

Current map behavior:

1. `LazyMap.svelte` client-only-imports `Map.svelte` on mount (MapLibre is WebGL/`window`-dependent, and this keeps the large `maplibre-gl` bundle out of the entry chunk).
2. MapLibre GL JS initializes a map from one of two Mapbox-GL style-spec documents: `frontend/src/lib/temp_base_style.json` or `frontend/src/lib/precip_base_style.json`.
3. Each style draws the basemap from **OpenFreeMap** vector tiles (`tiles.openfreemap.org/planet`, no API key) — water, roads, borders, and (precip only) city labels. There is no raster background layer.
4. `Forecast.svelte` derives `precipitationSoon` from the next 12 hours of forecast data; `Map.svelte` loads the precipitation style when it is true, otherwise the temperature style.

There is currently no weather/radar overlay and no map tile proxy on the backend. Future weather data (e.g. a LibreWXR radar overlay) should be rendered as a separate layer above the OpenFreeMap basemap.

> Map container note: MapLibre adds a `maplibregl-map` class (`position: relative`) to its container element. Do not also put a Tailwind `absolute`/positioning utility on that same element — the classes have equal specificity and MapLibre's wins, collapsing the canvas. Keep layout positioning on a parent wrapper and give MapLibre its own inner `<div>` (see `Map.svelte`).

## External Services

| Service | Code | Notes |
| --- | --- | --- |
| Pirate Weather | `server/forecast.go` | Forecast source. Uses `extend=hourly`. |
| Nominatim | `server/geocode.go` | Search and reverse geocoding. Sets `User-Agent: rainbgone/1.0`. |
| OpenFreeMap | map style JSON files | Vector basemap tiles + glyphs (roads, borders, water, city labels). No API key required. |

## Caching

All caches use `server/cache.go`.

| Cache | Key | TTL |
| --- | --- | --- |
| Forecast | `lat,lng` | 1 minute |
| Reverse geocode | `lat,lng` | 24 hours |
| Background tiles | `variant/z/x/y` | 5 minutes |

The cache is in memory only. It is not shared across processes or persisted across restarts.

## Frontend Routes

| Route | File | Purpose |
| --- | --- | --- |
| `/` | `frontend/src/routes/+page.server.ts` | Redirects to Seattle coordinates. |
| `/{address}` | `frontend/src/routes/[address]/+page.server.ts` | Geocodes an address and redirects to coordinates. |
| `/{lat},{lng}` | `frontend/src/routes/[lat=number],[lng=number]/` | Main forecast page. |
| `/details/{lat},{lng}/{date}` | `frontend/src/routes/details/[lat=number],[lng=number]/[date]/` | Daily forecast detail page. |

## API Proxy

The frontend uses `/api/...` for backend calls.

There are two proxy paths:

- `frontend/vite.config.ts` proxies `/api` during Vite development.
- `frontend/src/hooks.server.ts` proxies `/api` during SvelteKit server runtime.

Both strip the `/api` prefix before forwarding to the Go backend. For example:

```text
/api/forecast?lat=47.6062&lng=-122.3321
```

forwards to:

```text
http://localhost:8080/forecast?lat=47.6062&lng=-122.3321
```

## Docker

Development compose file:

- `docker-compose.dev.yml`
- Backend exposes `8080`.
- Frontend exposes Vite on `5173`.
- Frontend proxy target points at `http://backend:8080`.

Production-style compose file:

- `docker-compose.yml`
- Backend exposes `8080`.
- Frontend exposes Node adapter output on `3000`.
- Frontend proxy target points at `http://backend:8080`.

## Change Guidance

When changing backend API routes:

1. Update `server/server.go`.
2. Update affected frontend `/api/...` call sites.
3. Update `frontend/src/hooks.server.ts` only if the proxy contract changes.
4. Update this document.

When changing map layers:

1. Keep base geography/background tiles and weather overlays separate in the mental model and in route naming.
2. Verify provider tile URLs and zoom bounds before hard-coding layer IDs.
3. Confirm temperature and precipitation background variants both still render base geography.
4. Do not commit provider-generated style files that contain embedded API keys.
5. Render future weather observations, forecasts, radar, precipitation, or temperature data as a separate overlay above the background tiles.

When changing config:

1. Update `server/config.go` or frontend env usage.
2. Update Docker compose environment where needed.
3. Update `AGENTS.md` and this document.
