"use client";

import { ErrorView } from "@/components/error-view";
import { Button } from "@/components/ui/button";

export default function Error({
  reset
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <ErrorView codigo={500} mensaje="Ocurrió un error inesperado.">
      <Button variant="outline" onClick={reset} className="mt-2">
        Reintentar
      </Button>
    </ErrorView>
  );
}
