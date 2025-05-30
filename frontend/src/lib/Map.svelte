<script lang="ts">
  import "ol/ol.css";
  import Map from 'ol/Map';
  import View from 'ol/View';
  import TileLayer from 'ol/layer/Tile';
  import OSM from 'ol/source/OSM';
  import { useGeographic } from 'ol/proj';

  interface Props {
    location: {
      lat: number;
      lng: number;
    };
  }

  let { location }: Props = $props();

  function theMap(lat: number, lng: number) {
    return (element: HTMLDivElement) => {
      useGeographic();
      let map: Map = new Map({
        target: element,
        layers: [
          new TileLayer({
            source: new OSM(),
          }),
        ],
        view: new View({
          center: [lng, lat],
          zoom: 6,
        }),
      });

      $effect(() => {
        map.getView().setCenter([lng, lat]);
      })
      return () => {
        map.setTarget(undefined);
        map = null; // dereference for GC
      }
    }
  }
</script>

<div class="relative min-h-[350px] max-h-[800px] w-full bg-neutral-200 after:pt-[35%] after:block">
  <div {@attach theMap(location.lat, location.lng)} tabindex="-1" class="absolute top-0 left-0 w-full h-full select-none"></div>
</div>