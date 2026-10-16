import type { MetadataRoute } from "next";
import { env } from "@/env";
import { getMateriasVigentes } from "@/server/queries/materias-vigentes";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const materias = await getMateriasVigentes();
  const base = env.NEXT_PUBLIC_BASE_URL;

  return [
    { url: new URL("/", base).toString() },
    ...materias.map(({ codigo }) => ({ url: new URL(`/materia/${codigo}`, base).toString() }))
  ];
}
