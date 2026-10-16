import type { LucideIcon } from "lucide-react";
import type { Metadata } from "next";
import { Database, GraduationCap, HatGlasses, LayersPlus } from "lucide-react";
import Link from "next/link";
import { BackgroundBlob } from "@/components/background-blob";
import { BuscadorTrigger } from "@/components/buscador/buscador-trigger";
import { ExternalLink } from "@/components/external-link";
import { ComentariosDestacados } from "@/components/inicio/comentarios-destacados";
import { getInicio } from "@/server/queries/inicio";

const title = "FIUBA Reviews";
const description =
  "Encontrá calificaciones y comentarios de los docentes de la facultad, subidos por otros estudiantes. Basado en el legendario Dolly FIUBA.";

export const metadata: Metadata = {
  title,
  description,
  alternates: { canonical: "/" },
  robots: { index: true, follow: true, "max-snippet": -1, "max-image-preview": "large" },
  openGraph: { title, description, images: [{ url: "/opengraph-image", alt: title }] },
  twitter: { card: "summary", title, description }
};

const FEATURES: { icon: LucideIcon; titulo: string; descripcion: string }[] = [
  {
    icon: HatGlasses,
    titulo: "Reviews anónimas",
    descripcion: "Las calificaciones y comentarios agregados son totalmente anónimos."
  },
  {
    icon: Database,
    titulo: "Mismos datos de Dolly",
    descripcion: "Usamos los datos originales de Dolly recolectados durante más de 10 años."
  },
  {
    icon: LayersPlus,
    titulo: "Nuevos planes",
    descripcion:
      "Se agregaron las materias de los nuevos planes manteniendo comentarios de sus equivalencias."
  },
  {
    icon: GraduationCap,
    titulo: "Todas las ingenierías",
    descripcion: "Están disponibles todas las materias de las 11 carreras de ingeniería."
  }
];

export default async function Page() {
  const { comentarios, materiasPopulares } = await getInicio();

  return (
    <div className="relative isolate">
      <BackgroundBlob upperLeft lowerRight={false} />

      <main className="relative z-10 container mx-auto mb-4 grid gap-10 p-6 lg:grid-cols-2 lg:gap-12">
        <section id="hero" className="space-y-8 self-center">
          <div className="mx-auto space-y-4 text-center lg:max-w-lg">
            <h1 className="text-5xl font-semibold tracking-tighter sm:text-6xl">
              <span className="text-fiuba">FIUBA</span> <span>Reviews</span>
            </h1>
            <p className="mx-auto max-w-117 font-medium">{description}</p>
          </div>

          <div className="flex justify-center">
            <BuscadorTrigger variante="hero" />
          </div>

          <div className="mx-auto grid max-w-155 gap-6 sm:grid-cols-2">
            {FEATURES.map(({ icon: Icon, titulo, descripcion }) => (
              <article key={titulo} className="flex flex-1 items-start gap-3">
                <div className="p-2 select-none">
                  <Icon className="size-5.5 stroke-fiuba" />
                </div>
                <div className="space-y-2">
                  <p className="font-medium">{titulo}</p>
                  <p>{descripcion}</p>
                </div>
              </article>
            ))}
          </div>
        </section>

        <ComentariosDestacados comentarios={comentarios} />

        <section id="materias-populares" className="space-y-4 text-center">
          <h2 className="text-3xl font-semibold">Materias más populares</h2>
          <p>
            Las {materiasPopulares.length} más populares según la cantidad carreras que la cursan.
            Con los nuevos planes se unificaron bastantes materias que antes eran distinguidas por
            diferente código, pero realmente eran comunes a varias carreras.
          </p>

          <div className="mt-4 flex flex-wrap justify-center gap-2">
            {materiasPopulares.map((materia) => (
              <Link
                key={materia.codigo}
                href={`/materia/${materia.codigo}`}
                className="line-clamp-1 rounded-full border border-button-border bg-button-background px-3 py-1 text-sm"
              >
                {materia.codigo} • {materia.nombre}
              </Link>
            ))}
          </div>
        </section>

        <section id="acerca-del-proyecto" className="space-y-4 text-center">
          <h2 className="text-3xl font-semibold">Acerca del proyecto</h2>
          <p>
            <ExternalLink href="https://github.com/lugfi/dolly">Dolly FIUBA</ExternalLink> era el
            sitio original en donde los estudiantes de FIUBA publicaban calificaciones y comentarios
            de los docentes con los que cursaban. Como alumno, desde que entré a la facultad fue un
            recurso invaluable al momento de elegir cátedras al iniciar cada cuatrimestre.
          </p>
          <p>
            Ahora que Dolly ya no está en funcionamiento, me parece necesario mantener activa una
            plataforma donde los alumnos puedan comentar sobre sus experiencias con las diferente
            cátedras.
          </p>
        </section>
      </main>
    </div>
  );
}
