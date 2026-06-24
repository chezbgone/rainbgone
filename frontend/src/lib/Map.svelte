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
	<div class="absolute top-0 left-0 h-full w-full select-none">
		<!-- MapLibre owns this inner element (it sets position:relative on it); keeping our
		     absolute-fill positioning on the parent avoids the class conflict. -->
		<div {@attach theMap(location.lat, location.lng)} tabindex="-1" class="h-full w-full"></div>
	</div>
</div>
