"use client";

import Link from "next/link";
import { Menu, Monitor, Moon, SunMedium } from "lucide-react";
import { useTheme } from "next-themes";
import { GithubIcon } from "@/components/icons/github";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";
import { GITHUB_URL, STATS_URL } from "@/lib/constants";

export function MobileMenu() {
  const { setTheme } = useTheme();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        aria-label="Abrir menú"
        className="inline-flex size-9 items-center justify-center border border-layout-border bg-background outline-hidden md:hidden"
      >
        <Menu className="size-5" aria-hidden="true" />
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" sideOffset={8} className="w-56">
        <DropdownMenuLabel>Navegación</DropdownMenuLabel>
        <DropdownMenuItem asChild>
          <Link href="/">Inicio</Link>
        </DropdownMenuItem>
        <DropdownMenuItem asChild>
          <a href={STATS_URL} target="_blank" rel="noreferrer noopener">
            Estadísticas
          </a>
        </DropdownMenuItem>

        <DropdownMenuSeparator />
        <DropdownMenuLabel>Tema</DropdownMenuLabel>
        <DropdownMenuItem onSelect={() => setTheme("light")}>
          <SunMedium /> Claro
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => setTheme("dark")}>
          <Moon /> Oscuro
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => setTheme("system")}>
          <Monitor /> Sistema
        </DropdownMenuItem>

        <DropdownMenuSeparator />
        <DropdownMenuLabel>Contacto</DropdownMenuLabel>
        <DropdownMenuItem asChild>
          <a href={GITHUB_URL} target="_blank" rel="noreferrer noopener">
            <GithubIcon /> GitHub
          </a>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
