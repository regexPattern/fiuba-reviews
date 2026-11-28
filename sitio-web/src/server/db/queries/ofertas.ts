import "server-only";
import { desc, eq } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { type Cuatrimestre, esMismoCuatrimestre, esPosterior } from "@/lib/cuatrimestres";
import { tags } from "@/server/cache";
import { db, schema } from "@/server/db";

/**
 * - `al-dia`: la oferta procesada es la del último cuatrimestre registrado.
 * - `en-proceso`: llegó una oferta más nueva que la procesada y todavía no se aplicó.
 * - `desactualizada`: la oferta procesada es de un cuatrimestre anterior.
 * - `sin-oferta`: nunca se procesó una oferta de la carrera.
 */
export type EstadoOferta = "al-dia" | "en-proceso" | "desactualizada" | "sin-oferta";

export type OfertaCarrera = {
  codigo: number;
  nombre: string;
  planVigente: number | null;
  procesada: Cuatrimestre | null;
  /** Oferta recibida más nueva que la procesada, pendiente de procesar. */
  recibida: Cuatrimestre | null;
  estado: EstadoOferta;
};

const c = schema.cuatrimestre;
const cuatrimestre = { numero: c.numero, anio: c.anio };

async function ultimoPorCarrera(
  tabla: typeof schema.ofertaComisiones | typeof schema.ofertaComisionesRaw
) {
  const filas = await db
    .selectDistinctOn([tabla.codigoCarrera], {
      codigoCarrera: tabla.codigoCarrera,
      cuatrimestre
    })
    .from(tabla)
    .innerJoin(c, eq(c.codigo, tabla.codigoCuatrimestre))
    .orderBy(tabla.codigoCarrera, desc(c.anio), desc(c.numero));

  return new Map(filas.map((f) => [f.codigoCarrera, f.cuatrimestre]));
}

/** Devuelve el estado de la oferta de comisiones de cada carrera, ordenadas por nombre. */
export async function getOfertasCarreras() {
  "use cache";
  cacheLife("max");
  cacheTag(tags.ofertas);

  const [carreras, planes, ultimasProcesadas, ultimasRecibidas, [ultimoCuatrimestre]] =
    await Promise.all([
      db.select().from(schema.carrera),
      db
        .select({
          codigoCarrera: schema.plan.codigoCarrera,
          anio: schema.plan.anio
        })
        .from(schema.plan)
        .where(eq(schema.plan.vigente, true)),
      ultimoPorCarrera(schema.ofertaComisiones),
      ultimoPorCarrera(schema.ofertaComisionesRaw),
      db.select(cuatrimestre).from(c).orderBy(desc(c.anio), desc(c.numero)).limit(1)
    ]);

  const planesVigentes = new Map(planes.map((p) => [p.codigoCarrera, p.anio]));

  const ofertas = carreras
    .map(({ codigo, nombre }): OfertaCarrera => {
      const procesada = ultimasProcesadas.get(codigo) ?? null;
      const ultimaRecibida = ultimasRecibidas.get(codigo) ?? null;
      const recibida =
        ultimaRecibida && (!procesada || esPosterior(ultimaRecibida, procesada))
          ? ultimaRecibida
          : null;

      let estado: EstadoOferta;
      if (recibida) estado = "en-proceso";
      else if (!procesada) estado = "sin-oferta";
      else if (ultimoCuatrimestre && esMismoCuatrimestre(procesada, ultimoCuatrimestre))
        estado = "al-dia";
      else estado = "desactualizada";

      return {
        codigo,
        nombre,
        planVigente: planesVigentes.get(codigo) ?? null,
        procesada,
        recibida,
        estado
      };
    })
    .toSorted((a, b) => a.nombre.localeCompare(b.nombre, "es"));

  return { ofertas, ultimoCuatrimestre: ultimoCuatrimestre ?? null };
}
