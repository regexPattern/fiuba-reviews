"use client";

import { usePathname } from "next/navigation";
import { BuscadorTrigger } from "@/components/buscador/buscador-trigger";

/** En el inicio ya hay un buscador grande, así que el de la navbar se oculta. */
export function BuscadorTriggerNavbar() {
  const pathname = usePathname();

  if (pathname === "/") return null;

  return <BuscadorTrigger variante="navbar" />;
}
