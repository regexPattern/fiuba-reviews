import Link from "next/link";
import { Suspense } from "react";
import { BotonBuscar, CommandPalette } from "@/components/busqueda/command-palette";
import { getMateriasVigentes } from "@/server/db/queries/materias";

async function BusquedaGlobal() {
  const materias = await getMateriasVigentes();
  return <CommandPalette materias={materias} />;
}

export function Header() {
  return (
    <header className="h-(--header-height) border-b">
      <nav
        aria-label="Navegación principal"
        className="mx-auto flex h-full max-w-6xl items-center justify-between gap-2 px-4 md:px-6"
      >
        <Link href="/" className="flex min-h-11 items-center font-pixel text-lg">
          FIUBA Reviews
        </Link>
        <Suspense fallback={<BotonBuscar disabled />}>
          <BusquedaGlobal />
        </Suspense>
      </nav>
    </header>
  );
}
