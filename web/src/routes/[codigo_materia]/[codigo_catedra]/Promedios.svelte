<script lang="ts">
  import { ChevronDown, Star } from "@lucide/svelte";
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
  <Popover.Trigger
    class="flex items-center gap-2 border border-[#AB9E9C] bg-[#AB9E9C]/50 px-3 py-2"
  >
    <Star class="size-[16px] fill-yellow-500 stroke-yellow-700" />
    <span>Promedio: {promedio?.general.toFixed(1) || "-"}</span>
    {#if promedio}
      <ChevronDown class="size-[16px]" />
    {/if}
  </Popover.Trigger>

  {#if promedio}
    <Popover.Portal>
      <Popover.Content
        class="z-50 w-56 border border-border-muted bg-background/50 backdrop-blur-lg"
        align="start"
        sideOffset={6}
      >
        {JSON.stringify(promedio, null, 2)}
      </Popover.Content>
    </Popover.Portal>
  {/if}
</Popover.Root>
