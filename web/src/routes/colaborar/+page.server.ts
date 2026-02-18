import type { PageServerLoad } from "./$types";
import { asc, eq } from "drizzle-orm";
import { db, schema } from "$lib/server/db";

export const prerender = true;

export const load: PageServerLoad = async () => {
  const actualizacionesOfertas = await db
    .select({
      carrera: { codigo: schema.carrera.codigo, nombre: schema.carrera.nombre },
      cuatrimestre: {
        codigo: schema.cuatrimestre.codigo,
        numero: schema.cuatrimestre.numero,
        anio: schema.cuatrimestre.anio
      }
    })
    .from(schema.ofertaComisiones)
    .innerJoin(schema.carrera, eq(schema.carrera.codigo, schema.ofertaComisiones.codigoCarrera))
    .innerJoin(
      schema.cuatrimestre,
      eq(schema.cuatrimestre.codigo, schema.ofertaComisiones.codigoCuatrimestre)
    )
    .orderBy(asc(schema.carrera.nombre));

  return { actualizacionesOfertas };
};
