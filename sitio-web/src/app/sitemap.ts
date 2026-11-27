import type { MetadataRoute } from "next";
import { env } from "@/env";
import { getMateriasVigentes } from "@/server/db/queries/materias";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const materias = await getMateriasVigentes();
  const url = (path: string) => new URL(path, env.NEXT_PUBLIC_BASE_URL).toString();

  return [
    { url: url("/") },

    /*
     * Las páginas individuales de cada cátedra no se indexan, solo las páginas
     * de las materias. Esto es porque las cátedras pueden variar de
     * cuatrimestre a cuatrimestre, dejando los índices con referencias
     * inválidas.
     */
    ...materias.map(({ codigo }) => ({ url: url(`/materias/${codigo}`) }))
  ];
}
