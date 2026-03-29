<script lang="ts">
  import { resolve } from "$app/paths";
  import { Info, MessageSquarePlus } from "@lucide/svelte";
  import Comentarios from "./Comentarios.svelte";
  import Promedios from "./Promedios.svelte";

  type PromedioCalificaciones = {
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
    catedra: {
      codigo: string;
      codigoMateria: string;
      nombre: string;
      calificacion: number;
      docentes: {
        nombre: string;
        codigo: string;
        rol: string | null;
        cantidadCalificaciones: number;
        promedioCalificaciones: PromedioCalificaciones | null;
        resumenComentario: string | null;
        comentarios: {
          codigo: number;
          contenido: string;
          cuatrimestre: { numero: number; anio: number };
          esDeDolly: boolean;
        }[];
      }[];
    };
  }

  let { catedra }: Props = $props();
</script>

<div class="m-4 space-y-8 md:m-6">
  {#if catedra.docentes.length === 0}
    <p class="text-foreground-muted py-2 text-sm">Esta catedra no tiene docentes cargados.</p>
  {/if}

  {#each catedra.docentes as docente (docente.codigo)}
    <section id={docente.codigo} class="scroll-mt-17 space-y-3">
      <div>
        <h1 class="w-fit text-4xl font-semibold tracking-tight">
          {docente.nombre}
        </h1>
        {#if docente.rol}
          <small class="text-sm">({docente.rol})</small>
        {/if}
      </div>

      {#if docente.resumenComentario}
        <div class="divide-y divide-fiuba border border-fiuba bg-fiuba/45">
          <p class={`p-3 before:content-['"'] after:content-['"']`}>
            {docente.resumenComentario}
          </p>
          <div class="flex items-center gap-1 p-3 text-button-foreground select-none">
            <Info class="size-[16px]" />
            <span class="text-sm">Resumen generado con IA.</span>
          </div>
        </div>
      {/if}

      <div class="flex gap-2 text-sm text-button-foreground">
        <Promedios
          promedio={docente.promedioCalificaciones}
          cantidadCalificaciones={docente.cantidadCalificaciones}
        />

        <a
          href={resolve(`/calificar?docente=${docente.codigo}&catedra=${catedra.codigo}`)}
          class="flex items-center gap-2 border border-button-border bg-button-background px-3 py-2 hover:bg-button-hover hover:transition-colors"
        >
          <span>Calificar</span>
          <MessageSquarePlus
            class="size-4 fill-fiuba/50 stroke-[#665889] dark:fill-[#D1BCE3]"
          />
        </a>
      </div>

      {#if docente.comentarios.length > 0}
        <Comentarios comentarios={docente.comentarios} />
      {:else}
        <p class="text-foreground-muted py-2 text-sm">Docente sin comentarios</p>
      {/if}
    </section>
  {/each}
</div>
