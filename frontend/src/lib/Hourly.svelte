<script lang="ts">
	import Stripes from '$lib/common/Stripes.svelte';
	import type { Forecast } from './Forecast/types';

	interface Props {
		hourly: Forecast['hourly'];
	}

	let { hourly }: Props = $props();
	const hours = $derived(hourly.data.slice(0, 25));
</script>

<div class="text-center text-2xl font-light">
	{hourly.summary}
</div>
<Stripes {hours} />
