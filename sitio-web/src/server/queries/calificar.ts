import "server-only";
import { desc, eq } from "drizzle-orm";
import { db, schema } from "@/server/db";

/** Datos para el formulario de calificar. No se cachea: depende del request. */
export async function getDocenteParaCalificar(codigo: string) {
  const [docente] = await db
    .select({
      codigo: schema.docente.codigo,
      codigoMateria: schema.docente.codigoMateria,
      nombre: schema.docente.nombre,
      nombreSiu: schema.docente.nombreSiu,
      rol: schema.docente.rol
    })
    .from(schema.docente)
    .where(eq(schema.docente.codigo, codigo))
    .limit(1);

  return docente ?? null;
}

export function getCuatrimestresCalificables() {
  return db
    .select({
      codigo: schema.cuatrimestre.codigo,
      numero: schema.cuatrimestre.numero,
      anio: schema.cuatrimestre.anio
    })
    .from(schema.cuatrimestre)
    .orderBy(desc(schema.cuatrimestre.codigo))
    .limit(4);
}
