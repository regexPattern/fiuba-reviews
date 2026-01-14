<script lang="ts">
  import { Collapsible } from "bits-ui";

  type Comentario = {
    codigo: number;
    contenido: string;
    cuatrimestre: {
      numero: number;
      anio: number;
    };
    esDeDolly: boolean;
  };

  interface Props {
    comentarios: Comentario[];
  }

  let { comentarios }: Props = $props();

  const MAX_VISIBLES = 10;

  let comentariosVisibles = $derived(comentarios.slice(0, MAX_VISIBLES));
  let comentariosOcultos = $derived(comentarios.slice(MAX_VISIBLES));
  let cantidadOcultos = $derived(comentariosOcultos.length);

  let estaExpandido = $state(false);
</script>

{#snippet ComentarioItem(comentario: Comentario)}
  <div class="space-x-0.5 py-3">
    <p class={`inline before:content-['"'] after:content-['"']`}>
      {comentario.contenido.trim()}
    </p>

    <div class="inline-flex items-center gap-2 text-xs text-[#495883] select-none">
      <span class="border border-fiuba/30 bg-fiuba/10 px-1.5 py-0.5">
        {comentario.cuatrimestre.numero}C{comentario.cuatrimestre.anio}
      </span>
      {#if comentario.esDeDolly}
        <span class="border border-fiuba/30 bg-fiuba/10 px-1.5 py-0.5">Dolly FIUBA</span>
      {/if}
    </div>
  </div>
{/snippet}

<Collapsible.Root open={estaExpandido} onOpenChange={(open) => (estaExpandido = open)}>
  <div class="divide-y divide-layout-border">
    {#each comentariosVisibles as comentario (comentario.codigo)}
      {@render ComentarioItem(comentario)}
    {/each}

    {#if cantidadOcultos > 0}
      <Collapsible.Content class="divide-y divide-border-muted/75">
        {#each comentariosOcultos as comentario (comentario.codigo)}
          {@render ComentarioItem(comentario)}
        {/each}
      </Collapsible.Content>

      <div class="w-full text-center">
        <Collapsible.Trigger
          class="md:w-content mt-2 w-full cursor-pointer p-2 text-sm text-foreground-muted hover:underline md:w-fit"
        >
          {estaExpandido
            ? "Mostrar menos"
            : `Mostrar ${cantidadOcultos} comentario${cantidadOcultos === 1 ? "" : "s"} más`}
        </Collapsible.Trigger>
      </div>
    {/if}
  </div>
</Collapsible.Root>
