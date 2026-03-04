<script lang="ts">
  import { ChevronRight, Hourglass } from "@lucide/svelte";
</script>

{#snippet tecla(tecla: string)}
  <kbd class="border border-button-border bg-button-background px-1.25 py-0.5 text-sm">{tecla}</kbd>
{/snippet}

<section class="space-y-4">
  <h1 class="text-2xl font-semibold">¿Cómo enviar tu oferta de comisiones?</h1>
  <p>
    Si querés enviar una versión más reciente de la oferta de comisiones de tu carrera, lo podés
    hacer siguiendo los pasos listado a continuación:
  </p>
  <ol class="list-decimal px-[2ch]">
    <li>
      En el SIU, andá a
      <a
        href="https://guaraniautogestion.fi.uba.ar/g3w/oferta_comisiones"
        target="_blank"
        class="inline-flex items-center border border-button-border bg-button-background px-1 py-0.5 text-sm underline hover:bg-button-hover hover:transition-colors"
      >
        Reportes <ChevronRight class="size-4" /> Oferta de comisiones
      </a>
    </li>
    <li>
      Seleccioná todo el contenido de la página con {@render tecla("Ctrl+K")} en Windows/Linux o {@render tecla(
        "⌘K"
      )} en MacOS.
    </li>
    <li>
      Copia todo con {@render tecla("Ctrl+C")} en Windows/Linux o {@render tecla("⌘C")} en MacOS.
    </li>
    <li>
      Pega todo en este formulario con {@render tecla("Ctrl+V")} en Windows/Linux o {@render tecla(
        "⌘V"
      )} en MacOS.
    </li>
  </ol>
  <div
    class="flex items-center gap-2 border border-orange-500 bg-[#fdbb6f] p-3 text-sm text-orange-800 dark:border-orange-400 dark:bg-orange-700 dark:text-orange-200"
  >
    <Hourglass class="size-[18px] shrink-0" />
    <span>
      La actualización no va a ser instantánea porque todavía hay algunos errores que corregir en el
      programa que aplica los cambios.
    </span>
  </div>
</section>
