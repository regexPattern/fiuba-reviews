<script lang="ts">
  import { ChevronDown, Star } from "@lucide/svelte";
  import { Popover } from "bits-ui";

  const MAX_CALIFICACION = 5;

  type Promedio = {
    general: number;
    aceptaCritica: number;
    asistencia: number;
    buenTrato: number;
    claridad: number;
    claseOrganizada: number;
    cumpleHorarios: number;
    fomentaParticipacion: number;
    panoramaAmplio: number;
    respondeMails: number;
  };

  interface Props {
    promedio: Promedio | null;
    cantidadCalificaciones: number;
  }

  let { promedio, cantidadCalificaciones }: Props = $props();
</script>

{#snippet Criterio(label: string, valor: number)}
  {@const MIN_CALIFICACION = 1}
  {@const RANGO = MAX_CALIFICACION - MIN_CALIFICACION}
  {@const porcentaje = Math.max(0, Math.min(100, ((valor - MIN_CALIFICACION) / RANGO) * 100))}

  <div class="relative overflow-hidden px-2 py-1">
    <div class="absolute inset-0 bg-fiuba/5"></div>
    <div class="absolute inset-y-0 left-0 bg-fiuba/20" style:width={`${porcentaje}%`}></div>

    <div
      class="relative z-10 grid grid-cols-[1fr_auto] items-center gap-3 text-foreground tabular-nums"
    >
      <span>{label}</span>
      <span>{valor.toFixed(1)}</span>
    </div>
  </div>
{/snippet}

<Popover.Root>
  <Popover.Trigger
    class="flex items-center gap-2 border border-button-border bg-button-background px-3 py-2
      {promedio ? 'hover:bg-button-hover' : 'pointer-events-none'}"
  >
    <Star class="size-[16px] fill-yellow-500 stroke-yellow-700" />
    <span>Promedio: {promedio?.general.toFixed(1) || "–"}</span>
    {#if promedio}
      <ChevronDown class="size-[16px]" />
    {/if}
  </Popover.Trigger>

  {#if promedio}
    <Popover.Portal>
      <Popover.Content
        class="z-50 w-56 border border-button-border/30 bg-background/95 p-3 text-sm shadow-md backdrop-blur-xl"
        align="start"
        sideOffset={6}
      >
        <div class="divide-y-2 divide-background">
          {@render Criterio("Acepta Crítica", promedio.aceptaCritica)}
          {@render Criterio("Asistencia", promedio.asistencia)}
          {@render Criterio("Buen Trato", promedio.buenTrato)}
          {@render Criterio("Claridad", promedio.claridad)}
          {@render Criterio("Clase Organizada", promedio.claseOrganizada)}
          {@render Criterio("Cumple Horario", promedio.cumpleHorarios)}
          {@render Criterio("Fomenta Participación", promedio.fomentaParticipacion)}
          {@render Criterio("Panorama Amplio", promedio.panoramaAmplio)}
          {@render Criterio("Responde Mails", promedio.respondeMails)}
        </div>

        <div class="text-secondary-foreground pt-3 text-center">
          {cantidadCalificaciones} calificacion{cantidadCalificaciones === 1 ? "" : "es"}
        </div>
      </Popover.Content>
    </Popover.Portal>
  {/if}
</Popover.Root>
