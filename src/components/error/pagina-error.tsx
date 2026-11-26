import type { StaticImageData } from "next/image";
import Image from "next/image";
import colectivoFalla from "@/assets/colectivo-99.webp";
import colectivo400 from "@/assets/colectivo-400.webp";
import colectivo404 from "@/assets/colectivo-404.webp";
import colectivo500 from "@/assets/colectivo-500.webp";

const COLECTIVOS: Partial<Record<number, StaticImageData>> = {
  400: colectivo400,
  404: colectivo404,
  500: colectivo500
};

/**
 * Página de error con la foto del colectivo de la línea del código (o la del 99 si no hay una),
 * enmarcada como una foto impresa que cuelga de una chinche.
 */
export function PaginaError({
  codigo,
  descripcion,
  children
}: {
  /** Sin código (error desconocido) el título es "Error". */
  codigo?: number;
  descripcion: React.ReactNode;
  /** Acciones (links, botones) debajo del mensaje. */
  children?: React.ReactNode;
}) {
  const colectivo = (codigo !== undefined && COLECTIVOS[codigo]) || colectivoFalla;

  return (
    <main className="flex min-h-[calc(100dvh-var(--header-height))] flex-col items-center justify-center gap-10 px-6 py-12">
      <figure className="relative flex w-full max-w-xs flex-col items-center sm:max-w-sm">
        {/* Chinche vista desde atrás, un poco inclinada hacia arriba a la derecha: el canto, como el
            del `relieve`, asoma abajo a la izquierda. Tapa el vértice del hilo. */}
        <span
          aria-hidden
          className="absolute top-0.5 left-1/2 z-10 size-5 -translate-1/2 rounded-full border-[1.5px] border-fiuba-oscuro bg-fiuba shadow-[-1.5px_2.5px_0_var(--color-fiuba-oscuro)]"
        />
        {/* Hilo de la chinche al marco, a un 12% de cada borde. Se mete 4px debajo del marco (que se
            dibuja encima) para que no se vea dónde termina contra el borde. */}
        <svg
          aria-hidden
          className="-mb-1 h-14 w-[76%]"
          viewBox="0 0 100 100"
          preserveAspectRatio="none"
        >
          <path
            d="M0,100 L46,8 Q50,0 54,8 L100,100"
            fill="none"
            strokeWidth={2.5}
            strokeLinecap="round"
            vectorEffect="non-scaling-stroke"
            className="stroke-stone-300"
          />
        </svg>
        {/* Marco en relieve como las cards (canto abajo, por fuera) y la foto hundida en él (canto
            arriba, por dentro). */}
        <div className="relative w-full relieve bg-white p-3">
          <Image
            src={colectivo}
            alt={`Colectivo de la línea ${colectivo === colectivoFalla ? 99 : codigo}`}
            placeholder="blur"
            preload
            sizes="(min-width: 640px) 360px, 90vw"
            className="h-auto w-full rounded-sm"
          />
          <span
            aria-hidden
            className="pointer-events-none absolute inset-3 rounded-sm border-[1.5px] border-stone-300 shadow-[inset_0_3px_0_var(--color-stone-300)]"
          />
        </div>
      </figure>

      <div className="flex max-w-md flex-col items-center text-center">
        <h1 className="text-5xl font-medium tracking-tight text-stone-900">
          {codigo !== undefined ? `Error ${codigo}` : "Error"}
        </h1>
        <p className="mt-4 text-base/relaxed text-balance text-stone-600">{descripcion}</p>
        {children && <div className="mt-6 flex flex-wrap justify-center gap-3">{children}</div>}
      </div>
    </main>
  );
}
