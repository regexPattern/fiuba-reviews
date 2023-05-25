<script lang="ts">
	import type { PageData } from "./$types";
	import "../app.css";

	export let data: PageData;

	let materiaBuscada = "";
	let materiasFiltradas = data.materias;

	function filtrarMaterias() {
		materiasFiltradas = data.materias.filter((m) => m.codigo === 6103);
		console.log(data.materias.filter((m) => m.codigo === 6109));
	}
</script>

<form>
	<input bind:value={materiaBuscada} on:submit|preventDefault={filtrarMaterias} />
	<input type="submit" />
</form>
<ul>
	{#each materiasFiltradas as materia}
		<li>{materia.codigo} - {materia.nombre}</li>
	{/each}
</ul>
