import type { MetadataRoute } from "next";
import { env } from "@/env";
import { getMateriasVigentes } from "@/server/db/queries/materias";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const materias = await getMateriasVigentes();
  const url = (path: string) => new URL(path, env.NEXT_PUBLIC_BASE_URL).toString();

  return [
    { url: url("/") },

    /*
     * Solo se indexan las materias y no sus cátedras, porque las cátedras
     * cambian de un cuatrimestre a otro y los buscadores quedarían con links
     * rotos.
     */
    ...materias.map(({ codigo }) => ({ url: url(`/materias/${codigo}`) }))
  ];
}
