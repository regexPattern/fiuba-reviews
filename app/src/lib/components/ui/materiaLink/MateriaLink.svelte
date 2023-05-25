<script lang="ts">
	import type { materia } from "@prisma/client";

	let classes: string = "";
	export { classes as class };

  export let materia: materia;

  let codigo = materia.codigo_equivalencia || materia.codigo;
</script>

<a href={`/materias/${codigo}`} class={classes}>{materia.codigo} - {materia.nombre}</a>
