<script lang="ts">
  import "./layout.css";
  import favicon from "$lib/assets/favicon.svg";
  import Navbar from "./Navbar.svelte";
  import { ModeWatcher } from "mode-watcher";
  import Fuse from "fuse.js";
  import { afterNavigate } from "$app/navigation";
  import "@fontsource-variable/inter";
  import "@fontsource-variable/source-serif-4";

  let { children, data } = $props();

  afterNavigate((navigation) => {
    if (navigation.type === "popstate" || navigation.to?.url.hash) {
      return;
    }

    const contenedor = document.querySelector<HTMLElement>('[data-scroll-container="main"]');
    if (contenedor) {
      contenedor.scrollTo({ top: 0, left: 0 });
    } else {
      window.scrollTo({ top: 0, left: 0 });
    }
  });

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

{@render children()}
