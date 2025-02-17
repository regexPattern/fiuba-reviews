<script lang="ts">
  import { goto } from "$app/navigation";
  import PlaceholderCatedra from "$lib/components/placeholder-catedra.svelte";
  import type { PageData } from "./$types";
  import { onMount } from "svelte";

  export let data: PageData;

  onMount(async () => {
    goto(`/materias/${data.materia.codigo}/${data.catedras[0].codigo}`);
  });
</script>

<!-- TODO: Handle en el caso de que no haya catedras -->
<PlaceholderCatedra />
