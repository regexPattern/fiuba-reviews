import { revalidateTag } from "next/cache";
import { NextResponse } from "next/server";
import { timingSafeEqual } from "node:crypto";
import { z } from "zod";
import { env } from "@/env";

/**
 * Hay varios motivos por los que se puede invalidar el caché:
 *
 *  - `planesActualizados`: Si una actualización de los planes cambia las materias vigentes
 *    o sus equivalencias. No hace falta enviarlo en cada ejecución del actualizador del SIU
 *    si solo cambió la oferta de cátedras o docentes.
 *  - `materias`: Después de sincronizar la oferta del SIU, para las materias cuya información,
 *    cuatrimestre de actualización o listado de cátedras cambió. También si una calificación
 *    cambia el promedio o el orden de sus cátedras.
 *  - `catedras`: Si el actualizador del SIU cambia sus docentes, o se agregan, modifican o
 *    eliminan calificaciones o comentarios de esos docentes. Se envían los UUIDs de todas
 *    las cátedras afectadas, ya que un docente puede pertenecer a más de una.
 *  - `tags`: Para invalidaciones manuales o puntuales que no cubran los casos anteriores.
 *  - `inicio`: Contenido propio de la página de inicio.
 *
 * Una misma actualización puede requerir varios campos. Por ejemplo, una calificación
 * puede afectar tanto el detalle de una cátedra como el listado de cátedras de su materia.
 */

const bodySchema = z.object({
  planesActualizados: z.boolean().optional(),
  materias: z.array(z.string().min(1)).default([]),
  catedras: z.array(z.uuid()).default([]),
  tags: z.array(z.string().min(1)).default([]),
  inicio: z.boolean().optional()
});

function autorizado(request: Request) {
  const token = request.headers.get("authorization")?.replace(/^Bearer\s+/i, "") ?? "";
  const a = Buffer.from(token);
  const b = Buffer.from(env.REVALIDATE_SECRET);

  return a.length === b.length && timingSafeEqual(a, b);
}

export async function POST(request: Request) {
  if (!autorizado(request)) {
    return NextResponse.json({ error: "No autorizado" }, { status: 401 });
  }

  const parsed = bodySchema.safeParse(await request.json().catch(() => null));

  if (!parsed.success) {
    return NextResponse.json({ error: z.prettifyError(parsed.error) }, { status: 400 });
  }

  const { planesActualizados, materias, catedras, tags, inicio } = parsed.data;

  const aInvalidar = new Set<string>(tags);

  if (planesActualizados) aInvalidar.add("materias");
  if (inicio) aInvalidar.add("inicio");

  // El ranking depende de las cátedras resueltas y los comentarios de sus docentes.
  if (materias.length > 0 || catedras.length > 0) {
    aInvalidar.add("buscador-materias");
  }

  for (const codigo of materias) {
    aInvalidar.add(`materia:${codigo}`);
    aInvalidar.add(`catedras:${codigo}`);
  }

  for (const uuid of catedras) {
    aInvalidar.add(`catedra:${uuid}`);
  }

  for (const tag of aInvalidar) {
    revalidateTag(tag, "max");
  }

  return NextResponse.json({ revalidated: [...aInvalidar] });
}
