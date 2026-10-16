import { revalidateTag } from "next/cache";
import { NextResponse } from "next/server";
import { timingSafeEqual } from "node:crypto";
import { z } from "zod";
import { env } from "@/env";

const bodySchema = z.object({
  /** Lista de materias vigentes o nuevas: cambió el listado general. */
  materiasCambiaron: z.boolean().optional(),
  /** Códigos de materia que cambiaron (info de materia y/o cátedras). */
  materias: z.array(z.string().min(1)).default([]),
  /** UUIDs de cátedras cuyo detalle cambió. */
  catedras: z.array(z.uuid()).default([]),
  /** Tags arbitrarios, por si hace falta algo puntual. */
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

  const { materiasCambiaron, materias, catedras, tags, inicio } = parsed.data;

  const aInvalidar = new Set<string>(tags);

  if (materiasCambiaron) aInvalidar.add("materias");
  if (inicio) aInvalidar.add("inicio");

  for (const codigo of materias) {
    aInvalidar.add(`materia:${codigo}`);
    aInvalidar.add(`catedras:${codigo}`);
  }

  for (const uuid of catedras) {
    aInvalidar.add(`catedra:${uuid}`);
  }

  // Sirve la versión vieja y regenera en segundo plano.
  for (const tag of aInvalidar) {
    revalidateTag(tag, "max");
  }

  return NextResponse.json({ revalidated: [...aInvalidar] });
}
