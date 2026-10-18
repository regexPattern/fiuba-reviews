import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Header() {
  return (
    <header className="shrink-0 border-b pt-[env(safe-area-inset-top)]">
      <nav
        className="mx-auto flex h-14 max-w-7xl items-center justify-between px-4"
        aria-label="Navegación principal"
      >
        <Link href="/" className="font-semibold">
          FIUBA Reviews
        </Link>
        <Button variant="ghost" asChild>
          <Link href="/">Materias</Link>
        </Button>
      </nav>
    </header>
  );
}
