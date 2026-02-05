<script lang="ts">
  import favicon from "$lib/assets/favicon.svg";
  import Navbar from "./Navbar.svelte";
  import "./layout.css";
  import "@fontsource-variable/google-sans-code";
  import "@fontsource-variable/inter";
  import "@fontsource-variable/source-serif-4";
  import Fuse from "fuse.js";
  import { ModeWatcher } from "mode-watcher";

  let { children, data } = $props();

  const DEBOUNCE_TIMEOUT_MS = 300;

  let queryValue = $state("");
  let queryDebounced = $state("");
  let fuse = $derived(
    new Fuse(data.materias, {
      ignoreDiacritics: true,
      shouldSort: false,
      includeScore: false,
      threshold: 0.5,
      keys: ["codigo", "nombre"]
    })
  );
  let materias = $derived.by(() => {
    if (queryDebounced.trim() === "") {
      return data.materias;
    }
    return fuse.search(queryDebounced).map((r) => r.item);
  });

  $effect(() => {
    if (queryValue.trim() === "") {
      queryDebounced = "";
      return;
    }

    const handler = setTimeout(() => {
      queryDebounced = queryValue;
    }, DEBOUNCE_TIMEOUT_MS);

    return () => clearTimeout(handler);
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<ModeWatcher />
<Navbar {materias} />

<div class="pt-[calc(56px+env(safe-area-inset-top))]">
  {@render children()}
</div>
