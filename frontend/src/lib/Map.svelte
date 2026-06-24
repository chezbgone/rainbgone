<script lang="ts">
	import maplibregl from 'maplibre-gl';
	import 'maplibre-gl/dist/maplibre-gl.css';

	import temp_base_style from './temp_base_style.json';
	import precip_base_style from './precip_base_style.json';

	interface Props {
		location: {
			lat: number;
			lng: number;
		};
		precipitationSoon: boolean;
	}

	let { location, precipitationSoon }: Props = $props();

	const RADAR_METADATA_URL = 'https://api.librewxr.net/public/weather-maps.json';

	interface RadarFrame {
		time: number;
		path: string;
	}
	interface WeatherMaps {
		host: string;
		radar: { past: RadarFrame[] };
	}

	// Latest radar frame as a MapLibre raster tile URL template, or null on failure
	// (so the map degrades gracefully to just the precip basemap).
	async function latestRadarTileUrl(): Promise<string | null> {
		try {
			const res = await fetch(RADAR_METADATA_URL);
			if (!res.ok) return null;
			const data: WeatherMaps = await res.json();
			const latest = data.radar?.past?.at(-1);
			if (!latest) return null;
			// size=256, color scheme=8 (Dark Sky), smooth=1 (gaussian), snow=0
			return `${data.host}${latest.path}/256/{z}/{x}/{y}/8/1_0.png`;
		} catch {
			return null;
		}
	}

	function theMap(lat: number, lng: number) {
		return (element: HTMLDivElement) => {
			const style = (precipitationSoon
				? precip_base_style
				: temp_base_style) as unknown as maplibregl.StyleSpecification;

			let map: maplibregl.Map | null = new maplibregl.Map({
				container: element,
				style,
				center: [lng, lat],
				zoom: 6,
				maxZoom: 19
			});

			if (precipitationSoon) {
				map.on('load', () => {
					latestRadarTileUrl().then((tileUrl) => {
						if (!map || !tileUrl) return; // map torn down (cleanup nulls it) or fetch failed
						map.addSource('librewxr-radar', {
							type: 'raster',
							tiles: [tileUrl],
							tileSize: 256,
							attribution: 'Radar © <a href="https://librewxr.net/">LibreWXR</a> (CC-BY 4.0)'
						});
						map.addLayer({
							id: 'librewxr-radar',
							type: 'raster',
							source: 'librewxr-radar',
							paint: { 'raster-opacity': 0.7 }
						});
					});
				});
			}

			$effect(() => {
				map?.setCenter([lng, lat]);
			});

			return () => {
				map?.remove();
				map = null; // dereference for GC
			};
		};
	}
</script>

<div class="relative max-h-[800px] min-h-[350px] w-full bg-neutral-200 after:block after:pt-[35%]">
	<div
		class="group absolute top-0 left-0 h-full w-full select-none"
		tabindex="-1"
		onpointerdown={(e) => e.currentTarget.focus()}
	>
		<!-- MapLibre owns this inner element (it sets position:relative on it); keeping our
		     absolute-fill positioning on the parent avoids the class conflict.
		     The map starts pointer-events:none so wheel/touch pass through to the page;
		     focusing the parent (click/tap) flips it interactive via group-focus-within. -->
		<div
			{@attach theMap(location.lat, location.lng)}
			class="pointer-events-none h-full w-full group-focus-within:pointer-events-auto"
		></div>
	</div>
</div>
