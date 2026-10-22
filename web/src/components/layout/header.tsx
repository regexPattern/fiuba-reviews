import { Contact } from "@withicons/react/chrome";
import { ArrowUpRight } from "lucide-react";
import Link from "next/link";
import { Suspense } from "react";
import { CopySupportEmail } from "@/components/layout/copy-support-email";
import { CommandPalette } from "@/components/materia/command-palette";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { env } from "@/env";
import { getMateriasBuscador } from "@/server/db/queries/materias-buscador";

async function BusquedaGlobal() {
  const materias = await getMateriasBuscador();
  return <CommandPalette materias={materias} />;
}

export function Header() {
  return (
    <header className="fixed inset-x-0 top-0 z-40 border-b bg-background/75 pt-[env(safe-area-inset-top)] backdrop-blur-surface">
      <nav
        className="mx-auto flex h-14 max-w-7xl items-center justify-between gap-2 px-4"
        aria-label="Navegación principal"
      >
        <Link href="/" className="shrink-0 font-semibold" aria-label="FIUBA Reviews, inicio">
          <span className="sm:hidden">FR</span>
          <span className="hidden sm:inline">FIUBA Reviews</span>
        </Link>
        <div className="flex items-center gap-1 sm:gap-2">
          <Suspense
            fallback={
              <Button
                variant="ghost"
                className="h-11 sm:h-9"
                disabled
                aria-label="Cargando buscador"
              >
                Buscar
              </Button>
            }
          >
            <BusquedaGlobal />
          </Suspense>
          <Button asChild variant="ghost" size="icon">
            <a
              href={env.GITHUB_URL}
              target="_blank"
              rel="noopener noreferrer"
              aria-label="GitHub (abre en otra pestaña)"
              title="GitHub"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                <path d="M12 .297a12 12 0 0 0-3.793 23.385c.6.111.82-.261.82-.577v-2.234c-3.338.726-4.043-1.416-4.043-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.108-.775.419-1.305.762-1.605-2.665-.304-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.536-1.524.117-3.176 0 0 1.008-.322 3.301 1.23a11.52 11.52 0 0 1 6.006 0c2.291-1.552 3.297-1.23 3.297-1.23.655 1.652.243 2.873.119 3.176.77.84 1.235 1.91 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.216.694.825.576A12.001 12.001 0 0 0 12 .297Z" />
              </svg>
            </a>
          </Button>
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" className="px-2 sm:px-3" aria-label="Contacto">
                <Contact size={24} className="size-6" aria-hidden="true" />
                <span className="hidden sm:inline">Contacto</span>
              </Button>
            </PopoverTrigger>
            <PopoverContent
              align="end"
              className="flex w-auto flex-col gap-1 bg-popover/75 p-2 whitespace-nowrap shadow-glass backdrop-blur-surface"
              aria-label="Contacto"
            >
              <Button
                asChild
                variant="ghost"
                className="flex h-auto w-full justify-start p-2 font-normal has-[>svg]:px-2"
              >
                <a
                  href={env.LINKEDIN_URL}
                  target="_blank"
                  rel="noopener noreferrer"
                  aria-label="LinkedIn (abre en otra pestaña)"
                >
                  <svg
                    className="text-linkedin"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                    aria-hidden="true"
                  >
                    <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.049c.476-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433a2.062 2.062 0 1 1 0-4.124 2.062 2.062 0 0 1 0 4.124zM7.119 20.452H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z" />
                  </svg>
                  LinkedIn
                  <ArrowUpRight className="ml-auto" aria-hidden="true" />
                </a>
              </Button>
              <CopySupportEmail email={env.MAIL_SOPORTE} />
            </PopoverContent>
          </Popover>
        </div>
      </nav>
    </header>
  );
}
