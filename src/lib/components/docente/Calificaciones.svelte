<script lang="ts">
	import StarIcon from "~icons/material-symbols/star";

	export let docente: {
		promedio: number;
		respuestas: number;
		calificaciones: {
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
	};
</script>

<div class="flex flex-wrap space-x-1 text-sm">
	<div class="flex space-x-1">
		<StarIcon
			class="text-yellow-500"
			fill="currentColor"
			stroke="currentColor"
			aria-hidden="true"
		/>

		<span class="font-semibold">{docente.promedio.toFixed(1)}</span>
	</div>

	<span>&bull;</span>

	<span class="underline underline-offset-2">
		{docente.respuestas}
		{#if docente.respuestas == 1}
			Calificación
		{:else}
			Calificaciones
		{/if}
	</span>
</div>
