import type { MetadataRoute } from "next";
import { env } from "@/env";
import { getMateriasVigentes } from "@/server/db/queries/materias";

/** Solo la home y las materias: las páginas de cátedra son `noindex`. */
export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const materias = await getMateriasVigentes();
  const url = (path: string) => new URL(path, env.NEXT_PUBLIC_BASE_URL).toString();

  return [
    { url: url("/") },
    ...materias.map(({ codigo }) => ({ url: url(`/materias/${codigo}`) }))
  ];
}
