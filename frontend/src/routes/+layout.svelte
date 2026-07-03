<script lang="ts">
	import './app.css';
	import Header from '$lib/Header.svelte';
	import SearchBar from '$lib/SearchBar.svelte';

	import { page } from '$app/state';

	let { children } = $props();
</script>

<Header />
<SearchBar
	defaultAddress={page.data.geocode?.formatted_address ?? page.data.formattedAddress ?? ''}
	back={page.data.links?.back}
/>
{@render children()}
