<script lang="ts">
	import { formatWeekdayDate } from '$lib/common/time';

	interface Props {
		previous: string | null;
		previousLabel: string | null;
		next: string | null;
		nextLabel: string | null;
		time: number;
		timezone: string;
	}

	let { previous, previousLabel, next, nextLabel, time, timezone }: Props = $props();
</script>

<div class="mt-2 mb-10 flex items-baseline justify-center gap-6 text-center">
	{#if previous}
		<a href={previous} class="flex-1 text-right text-base text-sky-600 hover:underline">
			← {previousLabel}
		</a>
	{:else}
		<span class="flex-1"></span>
	{/if}
	<div class="text-lg text-black">{formatWeekdayDate(time, timezone)}</div>
	{#if next}
		<a href={next} class="flex-1 text-left text-base text-sky-600 hover:underline">
			{nextLabel} →
		</a>
	{:else}
		<span class="flex-1"></span>
	{/if}
</div>
