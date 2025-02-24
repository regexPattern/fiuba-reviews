<script lang="ts">
  import type { PageData } from "./$types";

  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  import sinCatedras from "$lib/assets/sin-catedras.webp";
  import SkeletonCatedra from "./skeleton-catedra.svelte";
  import Link from "$lib/components/link.svelte";

  export let data: PageData;

  onMount(async () => {
    const catedras = await data.catedras;
    if (catedras.length > 0) {
      goto(`/materias/${data.materia.codigo}/${catedras[0].codigo}`);
    }
  });
</script>

{#await data.catedras}
  <SkeletonCatedra />
{:then catedras}
  {#if catedras.length > 0}
    <SkeletonCatedra />
  {:else}
    <div class="space-y-6 text-center">
      <img
        alt="Steve de Minecraft."
        src={sinCatedras}
        class="mx-auto"
        height={337.5}
        width={150} />
      <p class="mx-auto max-w-lg pb-4">
        Aún no tenemos información de las cátedras de esta materia. Podés
        ayudarnos a actualizar los listados enviándonos tu plan de estudio.
        <Link href="/planes" class="underline" external>Más Información.</Link>
      </p>
    </div>
  {/if}
{/await}
