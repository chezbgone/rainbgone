<script lang="ts">
	import { formatWeekdayDate } from '$lib/common/time';

	interface Props {
		back: string;
		previous: string | null;
		previousLabel: string | null;
		next: string | null;
		nextLabel: string | null;
		time: number;
		timezone: string;
	}

	let { back, previous, previousLabel, next, nextLabel, time, timezone }: Props = $props();
</script>

<div class="my-5 text-center text-xl font-light">
	<a href={back} class="text-neutral-600 hover:text-neutral-900">← Go Back</a>
</div>

<div class="mt-2 mb-10 flex items-baseline justify-center gap-6 text-center">
	{#if previous}
		<a href={previous} class="flex-1 text-right text-base text-blue-600 hover:text-blue-800">
			← {previousLabel}
		</a>
	{:else}
		<span class="flex-1"></span>
	{/if}
	<div class="text-lg text-black">{formatWeekdayDate(time, timezone)}</div>
	{#if next}
		<a href={next} class="flex-1 text-left text-base text-blue-600 hover:text-blue-800">
			{nextLabel} →
		</a>
	{:else}
		<span class="flex-1"></span>
	{/if}
</div>
