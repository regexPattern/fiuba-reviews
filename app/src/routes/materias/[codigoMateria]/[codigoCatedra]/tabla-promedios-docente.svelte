<script lang="ts">
  export let cantidadCalificaciones: number;
  export let acepta_critica: number;
  export let asistencia: number;
  export let buen_trato: number;
  export let claridad: number;
  export let clase_organizada: number;
  export let cumple_horarios: number;
  export let fomenta_participacion: number;
  export let panorama_amplio: number;
  export let responde_mails: number;
</script>

<table class="mb-2">
  <tbody class="[&>tr>td:last-child]:pl-2">
    <tr>
      <td>Acepta Crítica</td>
      <td>{acepta_critica.toFixed(1)}</td>
    </tr>
    <tr>
      <td>Asistencia</td>
      <td>{asistencia.toFixed(1)}</td>
    </tr>
    <tr>
      <td>Buen Trato</td>
      <td>{buen_trato.toFixed(1)}</td>
    </tr>
    <tr>
      <td>Claridad</td>
      <td>{claridad.toFixed(1)}</td>
    </tr>
    <tr>
      <td>Clase Organizada</td>
      <td>{clase_organizada.toFixed(1)}</td>
    </tr>
    <tr>
      <td>Cumple Horario</td>
      <td>{cumple_horarios.toFixed(1)}</td>
    </tr>
    <tr>
      <td>Fomenta Participación</td>
      <td>{fomenta_participacion.toFixed(1)}</td>
    </tr>
    <tr>
      <td>Panorama Amplio</td>
      <td>{panorama_amplio.toFixed(1)}</td>
    </tr>
    <tr>
      <td>Responde Mails</td>
      <td>{responde_mails.toFixed(1)}</td>
    </tr>
  </tbody>
</table>
<span class="block text-center text-sm text-muted-foreground"
  >{cantidadCalificaciones} calificaci{cantidadCalificaciones > 1
    ? "ones"
    : "ón"}</span>
