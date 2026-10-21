import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Header() {
  return (
    <header className="fixed inset-x-0 top-0 z-40 border-b bg-background/75 pt-[env(safe-area-inset-top)] backdrop-blur-surface">
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
