import { revalidateTag } from "next/cache";
import { timingSafeEqual } from "node:crypto";
import { z } from "zod";
import { env } from "@/env";
import { tags } from "@/server/cache";

/**
 * Invalida el caché después de una actualización de datos. Se pueden mandar varios campos a la
 * vez, porque un mismo cambio puede afectar a más de uno.
 *
 *  - `planes`: cambiaron las materias vigentes o sus equivalencias.
 *  - `materias`: cambió la info o el listado de cátedras (incluidas sus calificaciones) de estas
 *    materias.
 *  - `catedras`: cambiaron los docentes, calificaciones o comentarios de estas cátedras.
 *  - `ofertas`: se procesaron ofertas de comisiones (estado por carrera en `/planes`).
 *  - `tags`: tags arbitrarios, para invalidaciones manuales.
 */
const SchemaRevalidacion = z.object({
  planes: z.boolean().default(false),
  ofertas: z.boolean().default(false),
  materias: z.array(z.string().min(1)).default([]),
  catedras: z.array(z.uuid()).default([]),
  tags: z.array(z.string().min(1)).default([])
});

function estaAutorizado(req: Request) {
  const token = req.headers.get("authorization")?.replace(/^Bearer\s+/i, "") ?? "";
  const recibido = Buffer.from(token);
  const esperado = Buffer.from(env.REVALIDATE_SECRET);

  return recibido.length === esperado.length && timingSafeEqual(recibido, esperado);
}

export async function POST(req: Request) {
  if (!estaAutorizado(req)) {
    return Response.json({ error: "No autorizado" }, { status: 401 });
  }

  const body = SchemaRevalidacion.safeParse(await req.json().catch(() => null));

  if (!body.success) {
    return Response.json({ error: z.prettifyError(body.error) }, { status: 400 });
  }

  const { planes, ofertas, materias, catedras } = body.data;

  const invalidados = new Set([
    ...body.data.tags,
    ...(planes ? [tags.materias] : []),
    ...(ofertas ? [tags.ofertas] : []),
    ...materias.flatMap((codigo) => [tags.materia(codigo), tags.catedras(codigo)]),
    ...catedras.map(tags.catedra)
  ]);

  for (const tag of invalidados) {
    revalidateTag(tag, "max");
  }

  return Response.json({ invalidados: [...invalidados] });
}
