<script lang="ts">
  import Link from "$lib/components/link.svelte";
  import { ModeWatcher, setMode } from "mode-watcher";
  import BuscadorMaterias from "./buscador-materias.svelte";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";

  import "../app.css";
  import type { LayoutData } from "./$types";
  import { Button } from "$lib/components/ui/button";
  import { Monitor, Moon, Sun } from "lucide-svelte";

  export let data: LayoutData;
</script>

<svelte:head>
  <meta
    name="description"
    content="Encontrá calificaciones y comentarios de los docentes de la facultad, subidos por otros estudiantes de la FIUBA." />
</svelte:head>

<ModeWatcher />

<header class="sticky top-0 z-[40] border-b bg-background/75 backdrop-blur-lg">
  <div
    class="container flex h-16 items-center justify-start gap-2 p-3 sm:justify-between">
    <div class="mr-2 flex-1">
      <Link href="/" class="contents text-lg" aria-label="Página de inicio">
        <span
          class="font-serif font-bold tracking-wide text-fiuba 2xs:tracking-tight"
          >F<span class="hidden 2xs:inline">IUBA </span></span
        ><span class="font-medium tracking-tighter"
          >R<span class="hidden 2xs:inline">eviews</span></span>
      </Link>
    </div>

    <nav>
      <ul>
        <li class="flex items-center">
          <BuscadorMaterias
            label="Materias"
            materias={data.materias}
            variant="outline" />
        </li>
      </ul>
    </nav>

    <DropdownMenu.Root positioning={{ placement: "bottom-end" }}>
      <DropdownMenu.Trigger asChild let:builder>
        <Button
          builders={[builder]}
          variant="outline"
          size="icon"
          class="h-10 w-10 shrink-0"
          aria-label="Cambiar tema">
          <Sun
            class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
          <Moon
            class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
        </Button>
      </DropdownMenu.Trigger>
      <DropdownMenu.Content class="z-[50]">
        <DropdownMenu.Item on:click={() => setMode("light")}>
          Claro <Sun class="ml-auto h-4 w-4" />
        </DropdownMenu.Item>
        <DropdownMenu.Item on:click={() => setMode("dark")}>
          Oscuro <Moon class="ml-auto h-4 w-4" />
        </DropdownMenu.Item>
        <DropdownMenu.Item on:click={() => setMode("system")}>
          Dispositivo <Monitor class="ml-auto h-4 w-4" />
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Root>

    <Link
      href="https://github.com/regexPattern/fiuba-reviews"
      class="flex h-10 w-10 shrink-0 items-center justify-center"
      aria-label="Repositorio del proyecto">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        x="0px"
        y="0px"
        viewBox="0 0 24 24"
        class="h-[26px] w-[26px] fill-black dark:fill-white">
        <path
          d="M10.9,2.1c-4.6,0.5-8.3,4.2-8.8,8.7c-0.5,4.7,2.2,8.9,6.3,10.5C8.7,21.4,9,21.2,9,20.8v-1.6c0,0-0.4,0.1-0.9,0.1 c-1.4,0-2-1.2-2.1-1.9c-0.1-0.4-0.3-0.7-0.6-1C5.1,16.3,5,16.3,5,16.2C5,16,5.3,16,5.4,16c0.6,0,1.1,0.7,1.3,1c0.5,0.8,1.1,1,1.4,1 c0.4,0,0.7-0.1,0.9-0.2c0.1-0.7,0.4-1.4,1-1.8c-2.3-0.5-4-1.8-4-4c0-1.1,0.5-2.2,1.2-3C7.1,8.8,7,8.3,7,7.6C7,7.2,7,6.6,7.3,6 c0,0,1.4,0,2.8,1.3C10.6,7.1,11.3,7,12,7s1.4,0.1,2,0.3C15.3,6,16.8,6,16.8,6C17,6.6,17,7.2,17,7.6c0,0.8-0.1,1.2-0.2,1.4 c0.7,0.8,1.2,1.8,1.2,3c0,2.2-1.7,3.5-4,4c0.6,0.5,1,1.4,1,2.3v2.6c0,0.3,0.3,0.6,0.7,0.5c3.7-1.5,6.3-5.1,6.3-9.3 C22,6.1,16.9,1.4,10.9,2.1z" />
      </svg>
    </Link>
  </div>
</header>

<slot />
