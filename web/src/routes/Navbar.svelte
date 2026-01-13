<script lang="ts">
  import { setMode, resetMode } from "mode-watcher";
  import { DropdownMenu } from "bits-ui";

  interface Props {
    materias: {
      codigo: string;
      nombre: string;
    }[];
  }

  let { materias }: Props = $props();
</script>

<header
  class="fixed top-0 left-0 z-[100] h-[56px] w-full border-b border-border-muted bg-background/50 px-4 backdrop-blur-lg"
>
  <ul>
    <a href="/">Inicio</a>
  </ul>

  <DropdownMenu.Root>
    <DropdownMenu.Trigger>Tema</DropdownMenu.Trigger>
    <DropdownMenu.Portal>
      <DropdownMenu.Content>
        <DropdownMenu.Item onclick={() => setMode("light")}>Light</DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => setMode("dark")}>Dark</DropdownMenu.Item>
        <DropdownMenu.Item onclick={resetMode}>System</DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Portal>
  </DropdownMenu.Root>
</header>
