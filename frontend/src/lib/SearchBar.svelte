<script lang="ts">
  import { goto, replaceState } from '$app/navigation';
  let props = $props();
  let address: string = $derived(props.defaultAddress);

  async function onSubmit(event: Event) {
    event.preventDefault();
    if (!address.trim()) {
      return;
    }
    goto(`/${address}`);
  }
</script>

<nav class="flex justify-center gap-4 bg-neutral-200 py-2">
  <form onsubmit={onSubmit} class="w-2/5">
    <input
      type="text"
      class="text-center bg-white focus:outline-none rounded-full px-4 w-full text-lg"
      bind:value={address}
    />
  </form>
</nav>