<script lang="ts">
  import Buscador from "./Buscador.svelte";
  import { Github, Menu, Monitor, Moon, SunMedium } from "@lucide/svelte";
  import { DropdownMenu } from "bits-ui";
  import { mode, resetMode, setMode } from "mode-watcher";

  interface Props {
    materias: {
      codigo: string;
      nombre: string;
    }[];
  }

  const LINK_GITHUB = "https://github.com/regexPattern/fiuba-reviews";

  let { materias }: Props = $props();
</script>

<header
  class="fixed top-0 left-0 z-100 h-[calc(56px+env(safe-area-inset-top))] w-full border-b border-layout-border bg-background/50 px-4 pt-[env(safe-area-inset-top)] backdrop-blur-lg">
  <div class="container mx-auto flex h-full items-center gap-2">
    <a
      href="/"
      class="mr-auto shrink-0 text-xl font-semibold tracking-tight"
      aria-label="Ir al inicio">
      <span class="font-serif text-fiuba">FIUBA</span> Reviews
    </a>

    <Buscador {materias} />

    <nav class="hidden items-center gap-5 md:mx-3 md:flex" aria-label="Navegación">
      <a href="/" class="text-sm hover:text-fiuba">Inicio</a>
      <!-- <a href="/estadisticas" class="text-sm hover:text-fiuba">Estadísticas</a> -->
      <a href="/colaborar" class="text-sm hover:text-fiuba">Colaborar</a>
    </nav>

    <DropdownMenu.Root>
      <DropdownMenu.Trigger
        class="hidden size-9 items-center justify-center text-sm font-medium md:inline-flex">
        {#if mode.current === "light"}
          <SunMedium class="size-[22px]" />
        {:else}
          <Moon class="size-[20px]" />
        {/if}
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          class="z-500 w-fit border border-layout-border bg-background p-1 shadow-lg data-[state=closed]:animate-out data-[state=closed]:fade-out data-[state=open]:animate-in data-[state=open]:fade-in"
          sideOffset={8}
          align="end">
          <DropdownMenu.Item
            class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
            onSelect={() => setMode("light")}>
            <SunMedium class="size-4" aria-hidden="true" />
            Claro
          </DropdownMenu.Item>
          <DropdownMenu.Item
            class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
            onSelect={() => setMode("dark")}>
            <Moon class="size-4" aria-hidden="true" />
            Oscuro
          </DropdownMenu.Item>
          <DropdownMenu.Item
            class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
            onSelect={resetMode}>
            <Monitor class="size-4" aria-hidden="true" />
            Sistema
          </DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>

    <a
      class="hidden size-9 items-center justify-center text-sm font-medium md:inline-flex"
      href={LINK_GITHUB}
      target="_blank"
      rel="noreferrer">
      <Github class="size-5" aria-hidden="true" />
    </a>

    <DropdownMenu.Root>
      <DropdownMenu.Trigger
        class="inline-flex size-9 items-center justify-center border border-layout-border bg-background md:hidden"
        aria-label="Abrir menú">
        <Menu class="size-5" aria-hidden="true" />
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          class="z-500 w-56 divide-layout-border border border-layout-border bg-background p-1 shadow-lg data-[state=closed]:animate-out data-[state=closed]:fade-out data-[state=open]:animate-in data-[state=open]:fade-in divide-y"
          sideOffset={8}
          align="end">
          <DropdownMenu.Group>
            <DropdownMenu.GroupHeading class="text-muted-foreground px-3 py-2 text-xs font-medium">
              Navegación
            </DropdownMenu.GroupHeading>
            <DropdownMenu.Item
              class="data-highlighted:bg-muted block px-3 py-2 text-sm outline-hidden">
              <a href="/" class="block">Inicio</a>
            </DropdownMenu.Item>
            <!-- <DropdownMenu.Item -->
            <!--   class="data-highlighted:bg-muted block px-3 py-2 text-sm outline-hidden"> -->
            <!--   <a href="/estadisticas">Estadísticas</a> -->
            <!-- </DropdownMenu.Item> -->
            <DropdownMenu.Item
              class="data-highlighted:bg-muted block px-3 py-2 text-sm outline-hidden">
              <a href="/colaborar">Colaborar</a>
            </DropdownMenu.Item>
          </DropdownMenu.Group>

          <DropdownMenu.Separator class="bg-border-muted my-1 h-px" />

          <DropdownMenu.Group>
            <DropdownMenu.GroupHeading class="text-muted-foreground px-3 py-2 text-xs font-medium">
              Tema
            </DropdownMenu.GroupHeading>
            <DropdownMenu.Item
              class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
              onSelect={() => setMode("light")}>
              <SunMedium class="size-4" aria-hidden="true" />
              Claro
            </DropdownMenu.Item>
            <DropdownMenu.Item
              class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
              onSelect={() => setMode("dark")}>
              <Moon class="size-4" aria-hidden="true" />
              Oscuro
            </DropdownMenu.Item>
            <DropdownMenu.Item
              class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
              onSelect={resetMode}>
              <Monitor class="size-4" aria-hidden="true" />
              Sistema
            </DropdownMenu.Item>
          </DropdownMenu.Group>

          <DropdownMenu.Separator class="bg-border-muted my-1 h-px" />

          <DropdownMenu.Group>
            <DropdownMenu.GroupHeading class="text-muted-foreground px-3 py-2 text-xs font-medium">
              Contacto
            </DropdownMenu.GroupHeading>
            <DropdownMenu.Item
              class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
              onSelect={() => window.open(LINK_GITHUB, "_blank")}>
              <Github class="size-4" aria-hidden="true" />
              GitHub
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  </div>
</header>
