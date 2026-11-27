import "server-only";
import { eq } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { type Cuatrimestre, ordenCuatrimestre } from "@/lib/ofertas";
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

/** Último cuatrimestre de cada carrera entre los pares `(carrera, cuatrimestre)`. */
function ultimoPorCarrera(filas: { codigoCarrera: number; cuatrimestre: Cuatrimestre }[]) {
  const ultimos = new Map<number, Cuatrimestre>();

  for (const { codigoCarrera, cuatrimestre } of filas) {
    const actual = ultimos.get(codigoCarrera);
    if (!actual || ordenCuatrimestre(cuatrimestre) > ordenCuatrimestre(actual)) {
      ultimos.set(codigoCarrera, cuatrimestre);
    }
  }

  return ultimos;
}

const c = schema.cuatrimestre;
const cuatrimestre = { numero: c.numero, anio: c.anio };

/** Estado de la oferta de comisiones de cada carrera, ordenadas por nombre. */
export async function getOfertasCarreras() {
  "use cache";
  cacheLife("max");
  cacheTag(tags.ofertas);

  const [carreras, planes, procesadas, recibidas, cuatrimestres] = await Promise.all([
    db.select().from(schema.carrera),
    db
      .select({ codigoCarrera: schema.plan.codigoCarrera, anio: schema.plan.anio })
      .from(schema.plan)
      .where(eq(schema.plan.vigente, true)),
    db
      .select({ codigoCarrera: schema.ofertaComisiones.codigoCarrera, cuatrimestre })
      .from(schema.ofertaComisiones)
      .innerJoin(c, eq(c.codigo, schema.ofertaComisiones.codigoCuatrimestre)),
    db
      .select({ codigoCarrera: schema.ofertaComisionesRaw.codigoCarrera, cuatrimestre })
      .from(schema.ofertaComisionesRaw)
      .innerJoin(c, eq(c.codigo, schema.ofertaComisionesRaw.codigoCuatrimestre)),
    db.select(cuatrimestre).from(c)
  ]);

  const ultimoCuatrimestre =
    cuatrimestres.toSorted((a, b) => ordenCuatrimestre(b) - ordenCuatrimestre(a))[0] ?? null;
  const ultimasProcesadas = ultimoPorCarrera(procesadas);
  const ultimasRecibidas = ultimoPorCarrera(recibidas);
  const planesVigentes = new Map(planes.map((p) => [p.codigoCarrera, p.anio]));

  const ofertas = carreras
    .map(({ codigo, nombre }): OfertaCarrera => {
      const procesada = ultimasProcesadas.get(codigo) ?? null;
      const ultimaRecibida = ultimasRecibidas.get(codigo) ?? null;
      const recibida =
        ultimaRecibida &&
        (!procesada || ordenCuatrimestre(ultimaRecibida) > ordenCuatrimestre(procesada))
          ? ultimaRecibida
          : null;

      let estado: EstadoOferta;
      if (recibida) estado = "en-proceso";
      else if (!procesada) estado = "sin-oferta";
      else if (
        ultimoCuatrimestre &&
        ordenCuatrimestre(procesada) >= ordenCuatrimestre(ultimoCuatrimestre)
      )
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

  return { ofertas, ultimoCuatrimestre };
}
