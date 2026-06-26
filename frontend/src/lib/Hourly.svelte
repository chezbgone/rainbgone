<script lang="ts">
	import Stripes from '$lib/common/Stripes.svelte';
	import type { Forecast } from './Forecast/types';

	interface Props {
		hourly: Forecast['hourly'];
	}

	let { hourly }: Props = $props();
	const hours = $derived(hourly.data.slice(0, 24));
</script>

<Stripes {hours} />
