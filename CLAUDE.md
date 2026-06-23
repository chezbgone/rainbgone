# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

rainbgone is a weather app with two runtimes: a Go HTTP backend (`server/`, started by `main.go`) and a SvelteKit frontend (`frontend/`). Read `@docs/architecture.md` before changing forecast, geocoding, map, tile proxy, config, Docker, or API behavior.

## Commands

Backend:
```sh
go run .
```

Frontend:
```sh
cd frontend
npm run dev       # Vite dev server
npm run check     # TypeScript + Svelte type check
npm run lint      # Prettier format check
npm run format    # Apply Prettier formatting
npm run build
```

Full dev stack (both services):
```sh
docker compose -f docker-compose.dev.yml up --build
```

## Validation

Run the narrowest relevant check for the files changed:

- Frontend TypeScript/Svelte changes: `cd frontend && npm run check`
- Formatting-only frontend changes: `cd frontend && npm run lint`
- Docker/runtime changes: validate compose syntax or build the touched service

There are no backend tests yet. If the user asks to skip validation, state that it was skipped.

## Configuration

The Go backend loads `.env` with `godotenv.Read()` at package init — it will **panic** on startup if `.env` is missing or unreadable. Current keys: `PIRATE_WEATHER_KEY`, `MAPTILER_KEY`.

The SvelteKit proxy uses `API_PROXY_TARGET` and `VITE_API_PROXY_TARGET` (both default to `http://localhost:8080`).

Do not add, rename, or remove config keys without updating `server/config.go`, Docker config, and `docs/architecture.md`.

## Secrets

Never print, quote, commit, or embed secrets from `.env` or provider-generated files (e.g., `style.json` may contain embedded MapTiler keys). If checking whether an env var is set, report only whether it is set — not its value.

## Frontend Code Style

Prettier is configured with non-default settings:
- Indentation: **tabs** (not spaces)
- Quotes: **single**
- Trailing commas: **none**
- Print width: **100**

TypeScript strict mode is enabled.

## API Routes

Routes are registered in `server/server.go`. When changing a backend route, update `server/server.go`, all frontend `/api/...` call sites, and `docs/architecture.md` in the same change.

## Commits

Use `type: short description`. Body is optional — only include one if it adds meaningful context beyond the subject line. Types in use: `feat`, `fix`, `refactor`, `perf`, `build`, `docs`.

## Map and Tile Changes

Background/base map tiles (MapTiler) and weather/radar overlays are conceptually separate — do not mix them. MapTiler tile IDs are provider-controlled; verify URLs before hard-coding. Do not commit provider-generated style files that contain embedded API keys.
