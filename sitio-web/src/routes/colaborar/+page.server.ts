import type { PageServerLoad } from "./$types";

import { db, schema } from "$lib/server/db";

import { asc, eq } from "drizzle-orm";

export const prerender = true;

export const load: PageServerLoad = async ({ setHeaders }) => {
  setHeaders({
    "x-robots-tag": "noindex, nofollow"
  });

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
