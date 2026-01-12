import { db } from "$lib/server/db";
import * as schema from "$lib/server/db/schema";
import type { PageServerLoad } from "./$types";
import { and, eq, sql } from "drizzle-orm";
import { exprPromedioDollyPorFila } from "$lib/server/db/utils";

export const load: PageServerLoad = async ({ params }) => {
  const docentesDeCatedra = db.$with("docentes_de_catedra").as(
    db
      .select({
        codigo: schema.docente.codigo,
        nombre: schema.docente.nombre,
        resumenComentarios: schema.docente.resumenComentarios,
        comentariosUltimoResumen: schema.docente.comentariosUltimoResumen,
        nombreSiu: schema.docente.nombreSiu,
        rol: schema.docente.rol
      })
      .from(schema.catedra)
      .innerJoin(
        schema.catedraDocente,
        eq(schema.catedraDocente.codigoCatedra, schema.catedra.codigo)
      )
      .innerJoin(schema.docente, eq(schema.docente.codigo, schema.catedraDocente.codigoDocente))
      .where(
        and(
          eq(schema.catedra.codigo, params.codigo_catedra),
          eq(schema.catedra.codigoMateria, params.codigo_materia)
        )
      )
  );

  const promedioFilaDollyExpr = exprPromedioDollyPorFila(schema.calificacionDolly);

  const metricasDocente = db.$with("metricas_docente").as(
    db
      .select({
        codigoDocente: schema.calificacionDolly.codigoDocente,
        promedioCalificaciones: sql<number>`avg(${promedioFilaDollyExpr})::float8`.as(
          "promedio_calificaciones"
        ),
        cantidadCalificaciones: sql<number>`count(${schema.calificacionDolly.codigo})::int`.as(
          "cantidad_calificaciones"
        )
      })
      .from(schema.calificacionDolly)
      .groupBy(schema.calificacionDolly.codigoDocente)
  );

  const comentariosDocente = db.$with("comentarios_docente").as(
    db
      .select({
        codigoDocente: schema.docente.codigo,
        comentarios: sql<
          Array<{
            codigo: number;
            contenido: string;
            anio: number;
            numero: number;
            fechaCreacion: string;
            esDeDolly: boolean;
          }>
        >`coalesce(
          json_agg(
            json_build_object(
              'codigo', ${schema.comentario.codigo},
              'contenido', ${schema.comentario.contenido},
              'anio', ${schema.cuatrimestre.anio},
              'numero', ${schema.cuatrimestre.numero},
              'fechaCreacion', ${schema.comentario.fechaCreacion},
              'esDeDolly', ${schema.comentario.esDeDolly}
            )
            order by ${schema.cuatrimestre.anio} desc, ${schema.cuatrimestre.numero} desc, ${schema.comentario.fechaCreacion} desc
          ),
          '[]'::json
        )`.as("comentarios")
      })
      .from(docentesDeCatedra)
      .innerJoin(schema.docente, eq(schema.docente.codigo, docentesDeCatedra.codigo))
      .leftJoin(schema.comentario, eq(schema.comentario.codigoDocente, schema.docente.codigo))
      .leftJoin(
        schema.cuatrimestre,
        eq(schema.cuatrimestre.codigo, schema.comentario.codigoCuatrimestre)
      )
      .groupBy(schema.docente.codigo)
  );

  const docentes = await db
    .with(docentesDeCatedra, metricasDocente, comentariosDocente)
    .select({
      codigo: docentesDeCatedra.codigo,
      nombre: docentesDeCatedra.nombre,
      resumenComentarios: docentesDeCatedra.resumenComentarios,
      comentariosUltimoResumen: docentesDeCatedra.comentariosUltimoResumen,
      nombreSiu: docentesDeCatedra.nombreSiu,
      rol: docentesDeCatedra.rol,
      promedioCalificaciones: metricasDocente.promedioCalificaciones,
      cantidadCalificaciones: metricasDocente.cantidadCalificaciones,
      comentarios: comentariosDocente.comentarios
    })
    .from(docentesDeCatedra)
    .leftJoin(metricasDocente, eq(metricasDocente.codigoDocente, docentesDeCatedra.codigo))
    .leftJoin(comentariosDocente, eq(comentariosDocente.codigoDocente, docentesDeCatedra.codigo));

  return { docentes };
};
