<script lang="ts">
  import "ol/ol.css";
  import Map from 'ol/Map';
  import View from 'ol/View';
  import TileLayer from 'ol/layer/Tile';
  import XYZ from 'ol/source/XYZ';
  import { useGeographic } from 'ol/proj';
  import { apply } from "ol-mapbox-style";

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
      useGeographic();

      const temperatureBackgroundLayer = new TileLayer({
        source: new XYZ({
          url: '/api/map/background-tiles/temp/{z}/{x}/{y}.png',
          tileSize: 256,
          maxZoom: 9,
        }),
      });

      const precipitationBackgroundLayer = new TileLayer({
        source: new XYZ({
          url: '/api/map/background-tiles/precipitation/{z}/{x}/{y}.png',
          tileSize: 256,
          maxZoom: 9,
        }),
      });

      let map: Map | null = new Map({
        target: element,
        layers: [
          precipitationSoon
            ? precipitationBackgroundLayer
            : temperatureBackgroundLayer
        ],
        view: new View({
          center: [lng, lat],
          zoom: 6,
          maxZoom: 19
        }),
      });
      
      apply(map, precipitationSoon ? precip_base_style : temp_base_style);

      $effect(() => {
        map?.getView().setCenter([lng, lat]);

      })
      return () => {
        map?.setTarget(undefined);
        map = null; // dereference for GC
      }
    }
  }
</script>

<div class="relative min-h-[350px] max-h-[800px] w-full bg-neutral-200 after:pt-[35%] after:block">
  <div {@attach theMap(location.lat, location.lng)} tabindex="-1" class="absolute top-0 left-0 w-full h-full select-none"></div>
</div>
