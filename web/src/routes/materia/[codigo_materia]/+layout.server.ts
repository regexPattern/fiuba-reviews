import { db, schema } from "$lib/server/db";
import { obtenerCatedras } from "$lib/server/db/utils";
import type { LayoutServerLoad } from "./$types";
import { error } from "@sveltejs/kit";
import { eq, and } from "drizzle-orm";

export const load: LayoutServerLoad = async ({ params }) => {
  const materiasVigentes = await db
    .select({
      codigo: schema.materia.codigo,
      nombre: schema.materia.nombre,
      cuatrimestreAnio: schema.cuatrimestre.anio,
      cuatrimestreNumero: schema.cuatrimestre.numero
    })
    .from(schema.materia)
    .innerJoin(schema.planMateria, eq(schema.planMateria.codigoMateria, schema.materia.codigo))
    .innerJoin(schema.plan, eq(schema.plan.codigo, schema.planMateria.codigoPlan))
    .leftJoin(
      schema.cuatrimestre,
      eq(schema.cuatrimestre.codigo, schema.materia.cuatrimestreUltimaActualizacion)
    )
    .where(
      and(eq(schema.plan.estaVigente, true), eq(schema.materia.codigo, params.codigo_materia))
    );

  if (materiasVigentes.length === 0) {
    error(404, "Materia no encontrada en los planes vigentes.");
  }

  const materia = materiasVigentes[0];

  const equivalencias = await db
    .select({
      codigo: schema.materia.codigo,
      nombre: schema.materia.nombre
    })
    .from(schema.equivalencia)
    .innerJoin(
      schema.materia,
      eq(schema.materia.codigo, schema.equivalencia.codigoMateriaPlanAnterior)
    )
    .where(eq(schema.equivalencia.codigoMateriaPlanVigente, materia.codigo));

  let codigosMateria: string[];
  let soloActivas: boolean;

  if (materia.cuatrimestreAnio !== null) {
    codigosMateria = [materia.codigo];
    soloActivas = true;
  } else {
    codigosMateria = equivalencias.map((e) => e.codigo);
    soloActivas = false;
  }

  const catedras = await obtenerCatedras(codigosMateria, { soloActivas });

  return {
    materia: {
      codigo: materia.codigo,
      nombre: materia.nombre,
      cuatrimestre:
        materia.cuatrimestreAnio !== null
          ? { numero: materia.cuatrimestreNumero!, anio: materia.cuatrimestreAnio }
          : null,
      equivalencias
    },
    catedras
  };
};
