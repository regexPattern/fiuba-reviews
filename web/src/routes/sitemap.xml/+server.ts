import type { RequestHandler } from "./$types";
import { db, schema } from "$lib/server/db";
import { eq } from "drizzle-orm";

const URL_BASE = "https://fiuba-reviews.com";

const escapeXml = (value: string) =>
  value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&apos;");

export const GET: RequestHandler = async () => {
  const materias = await db
    .select({ codigo: schema.materia.codigo })
    .from(schema.materia)
    .innerJoin(schema.planMateria, eq(schema.planMateria.codigoMateria, schema.materia.codigo))
    .innerJoin(schema.plan, eq(schema.plan.codigo, schema.planMateria.codigoPlan))
    .where(eq(schema.plan.estaVigente, true))
    .groupBy(schema.materia.codigo);

  const rutasEstaticas = ["/", "/colaborar"];
  const rutasMaterias = materias.map(({ codigo }) => `/materia/${codigo}`);
  const urls = [...rutasEstaticas, ...rutasMaterias];

  const body = `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">${urls
    .map((ruta) => {
      const loc = new URL(ruta, URL_BASE).toString();
      return `<url><loc>${escapeXml(loc)}</loc></url>`;
    })
    .join("")}</urlset>`;

  return new Response(body, {
    headers: {
      "content-type": "application/xml; charset=utf-8",
      "cache-control": "max-age=0, s-maxage=86400"
    }
  });
};
