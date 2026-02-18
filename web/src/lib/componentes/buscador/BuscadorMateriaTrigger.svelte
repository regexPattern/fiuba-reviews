<script lang="ts">
  import { Search } from "@lucide/svelte";
  import { buscadorMaterias } from "./materias.svelte";

  interface Props {
    variante?: "navbar" | "hero";
  }

  let { variante = "navbar" }: Props = $props();

  // TODO: cambiar esto para usar inline mejor.
  const clasesPorVariante = {
    navbar:
      "flex items-center gap-2 rounded-full border border-button-border bg-button-background px-3 py-2 text-sm text-button-foreground hover:bg-button-hover hover:transition-colors",
    hero: "flex items-center gap-2 rounded-full border border-button-border bg-button-background px-4 py-2.5 text-sm text-button-foreground hover:bg-button-hover hover:transition-colors"
  } as const;
</script>

<button
  type="button"
  class={clasesPorVariante[variante]}
  onclick={() => buscadorMaterias.abrirBuscador()}
  aria-label="Abrir buscador de materias"
>
  {#if variante === "hero"}
    <Search class="size-4" />
    <span>Buscar materias</span>
    <span class="hidden md:inline">⌘K</span>
  {:else}
    <span class="hidden items-center gap-2 md:flex">
      <Search class="size-4" />
      <span>Buscar materias</span>
    </span>
    <span class="flex items-center gap-2 md:hidden">
      <span>Buscar</span>
      <Search class="size-4" />
    </span>
    <span class="hidden md:inline">⌘K</span>
  {/if}
</button>
