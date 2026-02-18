<script lang="ts">
  import "@fontsource-variable/google-sans-code";
  import "@fontsource-variable/inter";
  import "@fontsource-variable/source-serif-4";
  import { ModeWatcher } from "mode-watcher";
  import Navbar from "./Navbar.svelte";
  import favicon from "$lib/assets/favicon.svg";
  import BuscadorMateriaDialog from "$lib/componentes/buscador/BuscadorMateriaDialog.svelte";
  import { buscadorMaterias } from "$lib/componentes/buscador/materias.svelte";
  import "./layout.css";

  let { children, data } = $props();

  $effect(() => {
    buscadorMaterias.setMaterias(data.materias);
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
  <title>FIUBA Reviews</title>
</svelte:head>

<ModeWatcher />
<Navbar />
<BuscadorMateriaDialog />

<div class="pt-[calc(56px+env(safe-area-inset-top))]">
  {@render children()}
</div>
