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
    <header className="fixed top-0 right-(--compensacion-scrollbar) left-0 z-40 border-b border-stone-300 bg-white/75 pt-[env(safe-area-inset-top)] backdrop-blur-xl">
      <nav
        aria-label="Navegación principal"
        className="mx-auto flex h-14 max-w-6xl items-center justify-between gap-2 px-4 md:px-6"
      >
        <Link
          href="/"
          className="flex min-h-11 items-center text-lg font-medium tracking-tight text-stone-900"
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
