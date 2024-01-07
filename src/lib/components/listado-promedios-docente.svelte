<script lang="ts">
	export let cantidadCalificaciones: number;
	export let aceptaCritica: number;
	export let asistencia: number;
	export let buenTrato: number;
	export let claridad: number;
	export let claseOrganizada: number;
	export let cumpleHorarios: number;
	export let fomentaParticipacion: number;
	export let panoramaAmplio: number;
	export let respondeMails: number;
</script>

<table class="mb-2">
	<tbody class="[&>tr>td:last-child]:pl-2">
		<tr>
			<td>Acepta Crítica</td>
			<td>{aceptaCritica.toFixed(1)}</td>
		</tr>
		<tr>
			<td>Asistencia</td>
			<td>{asistencia.toFixed(1)}</td>
		</tr>
		<tr>
			<td>Buen Trato</td>
			<td>{buenTrato.toFixed(1)}</td>
		</tr>
		<tr>
			<td>Claridad</td>
			<td>{claridad.toFixed(1)}</td>
		</tr>
		<tr>
			<td>Clase Organizada</td>
			<td>{claseOrganizada.toFixed(1)}</td>
		</tr>
		<tr>
			<td>Cumple Horario</td>
			<td>{cumpleHorarios.toFixed(1)}</td>
		</tr>
		<tr>
			<td>Fomenta Participación</td>
			<td>{fomentaParticipacion.toFixed(1)}</td>
		</tr>
		<tr>
			<td>Panorama Amplio</td>
			<td>{panoramaAmplio.toFixed(1)}</td>
		</tr>
		<tr>
			<td>Responde Mails</td>
			<td>{respondeMails.toFixed(1)}</td>
		</tr>
	</tbody>
</table>
<span class="block text-center text-sm text-muted-foreground"
	>{cantidadCalificaciones} calificaci{cantidadCalificaciones > 1 ? "ones" : "ón"}</span
>
