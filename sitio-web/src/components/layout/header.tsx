import Link from "next/link";
import { GithubIcon } from "@/components/icons/github";
import { GITHUB_URL, STATS_URL } from "@/lib/constants";
import { BuscadorTriggerNavbar } from "./buscador-trigger-navbar";
import { MobileMenu } from "./mobile-menu";
import { ThemeToggle } from "./theme-toggle";

export function Header() {
  return (
    <header className="fixed top-0 left-0 z-100 h-[calc(56px+env(safe-area-inset-top))] w-full border-b border-layout-border bg-background/80 px-4 pt-[env(safe-area-inset-top)] backdrop-blur-md">
      <div className="container mx-auto flex h-full items-center gap-2">
        <Link
          href="/"
          className="mr-auto shrink-0 font-serif text-xl font-semibold tracking-tight"
          aria-label="Ir al inicio"
        >
          <span className="text-fiuba">FIUBA</span> Reviews
        </Link>

        <BuscadorTriggerNavbar />

        <nav className="hidden items-center gap-5 md:mx-3 md:flex" aria-label="Navegación">
          <Link href="/" className="text-sm hover:text-fiuba">
            Inicio
          </Link>
          <a
            href={STATS_URL}
            target="_blank"
            rel="noreferrer noopener"
            className="text-sm hover:text-fiuba"
          >
            Estadísticas
          </a>
        </nav>

        <ThemeToggle className="hidden md:inline-flex" />

        <a
          href={GITHUB_URL}
          target="_blank"
          rel="noreferrer noopener"
          aria-label="Abrir GitHub"
          className="hidden size-9 items-center justify-center md:inline-flex"
        >
          <GithubIcon className="size-5" aria-hidden="true" />
        </a>

        <MobileMenu />
      </div>
    </header>
  );
}
