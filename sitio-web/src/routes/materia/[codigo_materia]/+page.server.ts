import type { EntryGenerator, PageServerLoad } from "./$types";
import { error } from "@sveltejs/kit";
import { db, schema } from "$lib/server/db";
import {
  getMateriasDisponibles,
  obtenerCatedrasMaterias,
  obtenerDetallesCatedras
} from "$lib/server/db/utils";
import { and, eq } from "drizzle-orm";

export const prerender = true;

export const entries: EntryGenerator = async () => {
  return (await getMateriasDisponibles()).map((m) => ({ codigo_materia: m.codigo }));
};

export const load: PageServerLoad = async ({ params }) => {
  const materiasRows = await db
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

  if (materiasRows.length === 0) {
    error(404, "Materia no encontrada en los planes vigentes.");
  }

  const materia = materiasRows[0];

  const equivalenciasRows = await db
    .select({ codigo: schema.materia.codigo, nombre: schema.materia.nombre })
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
    codigosMateria = equivalenciasRows.map((e) => e.codigo);
    soloActivas = false;
  }

  const catedras = await obtenerCatedrasMaterias(codigosMateria, soloActivas);
  const docentesPorCatedra =
    catedras.length > 0 ? await obtenerDetallesCatedras(catedras.map((c) => c.codigo)) : {};

  const catedrasConDocentes = catedras.map((catedra) => ({
    ...catedra,
    docentes: docentesPorCatedra[catedra.codigo] ?? []
  }));

  return {
    materia: {
      codigo: materia.codigo,
      nombre: materia.nombre,
      cuatrimestre:
        materia.cuatrimestreAnio !== null
          ? { numero: materia.cuatrimestreNumero!, anio: materia.cuatrimestreAnio }
          : null,
      equivalencias: equivalenciasRows
    },
    catedras: catedrasConDocentes
  };
};
