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

<div class="relative min-h-[350px] max-h-[800px] w-full bg-neutral-200 after:pt-[35%] after:block">
  <div {@attach theMap(location.lat, location.lng)} tabindex="-1" class="absolute top-0 left-0 w-full h-full select-none"></div>
</div>
