import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { Suspense } from "react";
import { CatedraSlot } from "@/components/materia/catedra-slot";
import { MateriaHeader } from "@/components/materia/materia-header";
import { MateriaShell } from "@/components/materia/materia-shell";
import { SidebarCatedras } from "@/components/materia/sidebar-catedras";
import { CatedraSkeleton, SidebarSkeleton } from "@/components/materia/skeletons";
import { getMateria } from "@/server/queries/materia";
import { getMateriasVigentes } from "@/server/queries/materias-vigentes";

type Props = {
  params: Promise<{ codigo: string }>;
  searchParams: Promise<{ catedra?: string | string[] }>;
};

export async function generateStaticParams() {
  const materias = await getMateriasVigentes();

  return materias.map(({ codigo }) => ({ codigo }));
}

export async function generateMetadata({ params }: Pick<Props, "params">): Promise<Metadata> {
  const { codigo } = await params;
  const materia = await getMateria(codigo);

  if (!materia) return {};

  const title = `${materia.codigo} • ${materia.nombre} | FIUBA Reviews`;
  const description = `Visitá la página de ${materia.nombre} para ver calificaciones y comentarios de las cátedras.`;
  const imageAlt = `FIUBA Reviews Materia ${materia.codigo}`;
  const url = `/materia/${materia.codigo}`;

  return {
    title,
    description,
    alternates: { canonical: url },
    robots: { index: true, follow: true, "max-snippet": -1, "max-image-preview": "large" },
    openGraph: {
      title,
      description,
      url,
      images: [{ url: `${url}/opengraph-image`, alt: imageAlt }]
    },
    twitter: {
      card: "summary_large_image",
      title,
      description,
      images: [{ url: `${url}/opengraph-image`, alt: imageAlt }]
    }
  };
}

export default async function Page({ params, searchParams }: Props) {
  const { codigo } = await params;
  const materia = await getMateria(codigo);

  if (!materia) notFound();

  const sidebar = (
    <aside className="flex h-full min-h-0 w-full flex-col border-r border-layout-border md:pt-14">
      <MateriaHeader codigo={codigo} />
      <Suspense fallback={<SidebarSkeleton />}>
        <SidebarCatedras codigo={codigo} />
      </Suspense>
    </aside>
  );

  return (
    <MateriaShell nombreMateria={materia.nombre} sidebar={sidebar}>
      <Suspense fallback={<CatedraSkeleton />}>
        <CatedraSlot codigo={codigo} searchParams={searchParams} />
      </Suspense>
    </MateriaShell>
  );
}
