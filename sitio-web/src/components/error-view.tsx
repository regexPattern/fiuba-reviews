import Image from "next/image";
import colectivo99 from "@/assets/colectivo-99.webp";
import colectivo400 from "@/assets/colectivo-400.webp";
import colectivo404 from "@/assets/colectivo-404.webp";
import colectivo500 from "@/assets/colectivo-500.webp";
import { BackgroundBlob } from "@/components/background-blob";
import { GITHUB_URL } from "@/lib/constants";

const IMAGENES: Record<number, typeof colectivo99> = {
  400: colectivo400,
  404: colectivo404,
  500: colectivo500
};

export function ErrorView({
  codigo,
  mensaje,
  children
}: {
  codigo: number;
  mensaje: string;
  children?: React.ReactNode;
}) {
  return (
    <div className="relative isolate">
      <BackgroundBlob upperLeft lowerRight />

      <div className="relative z-10 mx-auto space-y-6 p-6 sm:p-8">
        <Image
          src={IMAGENES[codigo] ?? colectivo99}
          alt={`Error ${codigo}`}
          className="mx-auto max-h-[468px] w-auto"
          priority
        />
        <div className="mx-auto max-w-[512px] space-y-4 text-center">
          <div className="space-y-2">
            <h1 className="text-4xl font-semibold sm:text-6xl">Error {codigo}</h1>
            <p>{mensaje}</p>
            {children}
          </div>
          <div className="mx-auto h-px w-[72px] bg-foreground/30" aria-hidden="true" />
          <div className="text-sm text-foreground/75">
            Si considerás que este error no debería haber ocurrido podes reportarlo en el{" "}
            <a
              href={GITHUB_URL}
              target="_blank"
              rel="noreferrer"
              className="underline underline-offset-2"
            >
              repositorio
            </a>{" "}
            o escribirme un mail{" "}
            <a href="mailto:ccastillo@fi.uba.ar" className="underline underline-offset-2">
              acá
            </a>
            .
          </div>
        </div>
      </div>
    </div>
  );
}
