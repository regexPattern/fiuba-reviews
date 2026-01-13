<script lang="ts">
  import { Popover } from "bits-ui";

  interface Props {
    promedio: {
      general: number;
      aceptaCritica: number;
      asistencia: number;
      buenTrato: number;
      claridad: number;
      claseOrganizada: number;
      cumpleHorarios: number;
      fomentaParticipacion: number;
      panoramaAmplio: number;
      respondeMails: number;
    } | null;
  }

  let { promedio }: Props = $props();
</script>

<Popover.Root>
  <Popover.Trigger>
    Promedio: {promedio?.general.toFixed(1) || "-"}
  </Popover.Trigger>

  {#if promedio}
    <Popover.Portal>
      <Popover.Content class="w-56 bg-white">
        {JSON.stringify(promedio, null, 2)}
      </Popover.Content>
    </Popover.Portal>
  {/if}
</Popover.Root>
