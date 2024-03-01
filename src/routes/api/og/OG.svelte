<script lang="ts">
	export let materia: {
		codigo: number;
		nombre: string;
	};
</script>

<div
	style="height: 100%; width: 100%; padding: 4rem; background-color: white; font-family: Geist; font-size: 2.5rem"
>
	<span>{materia.codigo}</span>
	<span>{materia.nombre}</span>
</div>
