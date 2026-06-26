<script lang="ts">
	import type { PageProps } from './$types';
	import Forecast from '$lib/Forecast/Forecast.svelte';

	let { data }: PageProps = $props();
</script>

<Forecast geocode={data.geocode} forecast={data.forecast} />
