<script lang="ts">
	import { goto, replaceState } from '$app/navigation';

	interface Props {
		defaultAddress: string;
		back?: string | null;
	}

	let props: Props = $props();
	let address: string = $derived(props.defaultAddress);

	async function onSubmit(event: Event) {
		event.preventDefault();
		if (!address.trim()) {
			return;
		}
		goto(`/${address}`);
	}
</script>

<nav class="relative flex justify-center gap-4 bg-neutral-200 py-2">
	{#if props.back}
		<a
			href={props.back}
			class="absolute top-1/2 left-4 -translate-y-1/2 text-sky-600 hover:underline"
		>
			← Go Back
		</a>
	{/if}
	<form onsubmit={onSubmit} class="w-2/5">
		<input
			type="text"
			class="w-full rounded-full bg-white px-4 text-center text-lg focus:outline-none"
			bind:value={address}
		/>
	</form>
</nav>
