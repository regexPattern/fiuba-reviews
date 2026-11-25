import Link from "next/link";
import { Suspense } from "react";
import { BotonBuscar, CommandPalette } from "@/components/busqueda/command-palette";
import { getMateriasVigentes } from "@/server/db/queries/materias";

async function BusquedaGlobal() {
  const materias = await getMateriasVigentes();
  return <CommandPalette materias={materias} />;
}

/** Fijo y translúcido, el contenido scrollea por debajo. Suma el safe area superior por `viewport-fit=cover`. */
export function Header() {
  return (
    <header className="fixed inset-x-0 top-0 z-40 border-b border-stone-300 bg-white/75 pt-[env(safe-area-inset-top)] backdrop-blur-xl">
      <nav
        aria-label="Navegación principal"
        className="mx-auto flex h-14 max-w-6xl items-center justify-between gap-2 px-4 md:px-6"
      >
        <Link
          href="/"
          className="flex min-h-11 items-center font-acento text-lg font-medium text-tinta"
        >
          FIUBA Reviews
        </Link>
        <Suspense fallback={<BotonBuscar disabled />}>
          <BusquedaGlobal />
        </Suspense>
      </nav>
    </header>
  );
}
