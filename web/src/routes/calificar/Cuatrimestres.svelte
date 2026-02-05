<script lang="ts">
  import { getCuatrimestres } from "./data.remote";

  let cuatrimestres = await getCuatrimestres();
</script>

Hola, ya estoy aca
{#each cuatrimestres as { codigo, numero, anio } (codigo)}
  <p>{numero}C{anio}</p>
{/each}
