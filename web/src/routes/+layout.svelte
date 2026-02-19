<script lang="ts">
  import { resolve } from "$app/paths";
  import { page } from "$app/state";
  import favicon from "$lib/assets/favicon.svg";
  import BuscadorMateriaDialog from "$lib/componentes/buscador/BuscadorMateriaDialog.svelte";
  import BuscadorMateriaTrigger from "$lib/componentes/buscador/BuscadorMateriaTrigger.svelte";
  import { buscadorMaterias } from "$lib/componentes/buscador/materias.svelte";
  import { Github, Menu, Monitor, Moon, SunMedium } from "@lucide/svelte";
  import { DropdownMenu } from "bits-ui";
  import { mode, ModeWatcher, resetMode, setMode } from "mode-watcher";
  import "@fontsource-variable/google-sans-code";
  import "@fontsource-variable/inter";
  import "@fontsource-variable/source-serif-4";
  import "./layout.css";

  const META_TITULO = "FIUBA Reviews";
  const META_DESCRIPCION =
    "Encontrá calificaciones y comentarios de los docentes de la facultad, subidos por otros estudiantes. Basado en el legendario Dolly FIUBA.";
  const META_URL_BASE = "https://fiuba-reviews.com";

  let { children, data } = $props();

  $effect(() => {
    buscadorMaterias.setMaterias(data.materias);
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
  <title>{META_TITULO}</title>
  <meta name="description" content={META_DESCRIPCION} />
  <meta name="author" content="Carlos Eduardo Castillo Pereira" />
  <meta property="og:type" content="website" />
  <meta property="og:site_name" content={META_TITULO} />
  <meta property="og:locale" content="es_AR" />
  <meta property="og:url" content={page.url.href} />
  <meta property="og:title" content={META_TITULO} />
  <meta property="og:description" content={META_DESCRIPCION} />
  <!-- TODO: <meta property="og:image" content={`${META_URL_BASE}/og/home.png`} /> -->

  <meta name="twitter:card" content="summary" />
  <meta name="twitter:title" content={META_TITULO} />
  <meta name="twitter:description" content={META_DESCRIPCION} />
  <!-- TODO: <meta name="twitter:image" content={`${META_URL_BASE}/og/home.png`} /> -->
</svelte:head>

<ModeWatcher />

<header
  class="fixed top-0 left-0 z-100 h-[calc(56px+env(safe-area-inset-top))] w-full border-b border-layout-border bg-background/80 px-4 pt-[env(safe-area-inset-top)] backdrop-blur-md"
>
  <div class="container mx-auto flex h-full items-center gap-2">
    <a
      href={resolve("/")}
      class="mr-auto shrink-0 font-serif text-xl font-semibold tracking-tight"
      aria-label="Ir al inicio"
    >
      <span class="text-fiuba">FIUBA</span> Reviews
    </a>

    {#if page.url.pathname !== resolve("/")}
      <BuscadorMateriaTrigger variante="navbar" />
    {/if}

    <nav class="hidden items-center gap-5 md:mx-3 md:flex" aria-label="Navegación">
      <a href={resolve("/")} class="text-sm hover:text-fiuba">Inicio</a>
      <a
        href="https://us.posthog.com/shared/c8cbP4SoFDIll_Niw7z2MaUbMRqEyA"
        target="_blank"
        rel="noreferrer"
        class="text-sm hover:text-fiuba">Estadísticas</a
      >
      <a href={resolve("/colaborar")} class="text-sm hover:text-fiuba">Colaborar</a>
    </nav>

    <DropdownMenu.Root>
      <DropdownMenu.Trigger
        class="hidden size-9 items-center justify-center text-sm font-medium md:inline-flex"
      >
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
          align="end"
        >
          <DropdownMenu.Item
            class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
            onSelect={() => setMode("light")}
          >
            <SunMedium class="size-4" aria-hidden="true" />
            Claro
          </DropdownMenu.Item>
          <DropdownMenu.Item
            class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
            onSelect={() => setMode("dark")}
          >
            <Moon class="size-4" aria-hidden="true" />
            Oscuro
          </DropdownMenu.Item>
          <DropdownMenu.Item
            class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
            onSelect={resetMode}
          >
            <Monitor class="size-4" aria-hidden="true" />
            Sistema
          </DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>

    <button
      type="button"
      class="hidden size-9 items-center justify-center text-sm font-medium md:inline-flex"
      aria-label="Abrir GitHub"
      onclick={() => window.open(META_URL_BASE, "_blank")}
    >
      <Github class="size-5" aria-hidden="true" />
    </button>

    <DropdownMenu.Root>
      <DropdownMenu.Trigger
        class="inline-flex size-9 items-center justify-center border border-layout-border bg-background md:hidden"
        aria-label="Abrir menú"
      >
        <Menu class="size-5" aria-hidden="true" />
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          class="z-500 w-56 border border-layout-border bg-background p-1 shadow-lg data-[state=closed]:animate-out data-[state=closed]:fade-out data-[state=open]:animate-in data-[state=open]:fade-in"
          sideOffset={8}
          align="end"
        >
          <DropdownMenu.Group>
            <DropdownMenu.GroupHeading class="text-muted-foreground px-3 py-2 text-xs font-medium">
              Navegación
            </DropdownMenu.GroupHeading>
            <DropdownMenu.Item textValue="Inicio">
              {#snippet child({ props })}
                <a
                  {...props}
                  href={resolve("/")}
                  class="data-highlighted:bg-muted block px-3 py-2 text-sm outline-hidden"
                >
                  Inicio
                </a>
              {/snippet}
            </DropdownMenu.Item>
            <!-- <DropdownMenu.Item -->
            <!--   class="data-highlighted:bg-muted block px-3 py-2 text-sm outline-hidden"> -->
            <!--   <a href="/estadisticas">Estadísticas</a> -->
            <!-- </DropdownMenu.Item> -->
            <DropdownMenu.Item textValue="Colaborar">
              {#snippet child({ props })}
                <a
                  {...props}
                  href={resolve("/colaborar")}
                  class="data-highlighted:bg-muted block px-3 py-2 text-sm outline-hidden"
                >
                  Colaborar
                </a>
              {/snippet}
            </DropdownMenu.Item>
          </DropdownMenu.Group>

          <DropdownMenu.Separator class="my-2 h-px bg-layout-border" />

          <DropdownMenu.Group>
            <DropdownMenu.GroupHeading class="text-muted-foreground px-3 py-2 text-xs font-medium">
              Tema
            </DropdownMenu.GroupHeading>
            <DropdownMenu.Item
              class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
              onSelect={() => setMode("light")}
            >
              <SunMedium class="size-4" aria-hidden="true" />
              Claro
            </DropdownMenu.Item>
            <DropdownMenu.Item
              class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
              onSelect={() => setMode("dark")}
            >
              <Moon class="size-4" aria-hidden="true" />
              Oscuro
            </DropdownMenu.Item>
            <DropdownMenu.Item
              class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
              onSelect={resetMode}
            >
              <Monitor class="size-4" aria-hidden="true" />
              Sistema
            </DropdownMenu.Item>
          </DropdownMenu.Group>

          <DropdownMenu.Separator class="my-2 h-px bg-layout-border" />

          <DropdownMenu.Group>
            <DropdownMenu.GroupHeading class="text-muted-foreground px-3 py-2 text-xs font-medium">
              Contacto
            </DropdownMenu.GroupHeading>
            <DropdownMenu.Item
              class="data-highlighted:bg-muted flex items-center gap-2 px-3 py-2 text-sm outline-hidden"
              onSelect={() => window.open(META_URL_BASE, "_blank")}
            >
              <Github class="size-4" aria-hidden="true" />
              GitHub
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  </div>
</header>
<BuscadorMateriaDialog />

<div class="pt-[calc(56px+env(safe-area-inset-top))]">
  {@render children()}
</div>
