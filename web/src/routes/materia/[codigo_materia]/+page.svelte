<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  let { data } = $props();

  onMount(() => {
    if (data.catedras.length > 0) {
      goto(`/materia/${data.materia.codigo}/${data.catedras[0].codigo}`, { replaceState: true });
    }
  });
</script>

<svelte:head>
  <title>FIUBA Reviews • {data.materia.codigo} • {data.materia.nombre}</title>
</svelte:head>
