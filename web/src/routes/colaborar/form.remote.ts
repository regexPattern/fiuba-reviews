import { error, invalid } from "@sveltejs/kit";
import { and, eq, max, sql } from "drizzle-orm";
import * as v from "valibot";
import { form } from "$app/server";
import { db, schema } from "$lib/server/db";
import { validateToken } from "$lib/server/turnstile";

export const submitForm = form(
  v.object({
    metadata: v.object({
      carrera: v.string(),
      cuatrimestre: v.object({ numero: v.number(), anio: v.number() })
    }),
    contenido: v.string(),
    cfTurnstileResponse: v.string()
  }),
  async ({ metadata, contenido, cfTurnstileResponse }) => {
    const { success } = await validateToken(cfTurnstileResponse);

    if (!success) {
      invalid("CAPTCHA inválido.");
    }

    const carreraNuevaOferta = await db
      .select({ codigo: schema.carrera.codigo })
      .from(schema.carrera)
      .where(
        sql`lower(unaccent(trim(${schema.carrera.nombre}))) = lower(unaccent(trim(${metadata.carrera})))`
      );

    const mensajesErrores = [];

    if (carreraNuevaOferta.length === 0) {
      mensajesErrores.push("La carrera especificada no existe en la base de datos.");
    }

    const cuatrimestreNuevaOferta = await db
      .select({ codigo: schema.cuatrimestre.codigo })
      .from(schema.cuatrimestre)
      .where(
        and(
          eq(schema.cuatrimestre.numero, metadata.cuatrimestre.numero),
          eq(schema.cuatrimestre.anio, metadata.cuatrimestre.anio)
        )
      );

    if (cuatrimestreNuevaOferta.length === 0) {
      mensajesErrores.push("El cuatrimestre especificado no existe en la base de datos.");
    }

    if (mensajesErrores.length > 0) {
      invalid(...mensajesErrores);
    }

    const ofertaActualCarrera = await db
      .select({ maxCodigoCuatrimestre: max(schema.ofertaComisiones.codigoCuatrimestre) })
      .from(schema.ofertaComisiones)
      .where(eq(schema.ofertaComisiones.codigoCarrera, carreraNuevaOferta[0].codigo));

    const cuatrimestreOfertaActual = ofertaActualCarrera[0]?.maxCodigoCuatrimestre;

    if (
      cuatrimestreOfertaActual !== null &&
      cuatrimestreOfertaActual >= cuatrimestreNuevaOferta[0].codigo
    ) {
      invalid("Ya existe una oferta más reciente para esta carrera");
    }

    try {
      await db
        .insert(schema.ofertaComisionesRaw)
        .values({
          codigoCarrera: carreraNuevaOferta[0].codigo,
          codigoCuatrimestre: cuatrimestreNuevaOferta[0].codigo,
          contenido
        });
    } catch (e) {
      console.error("[enviarOferta] Error al insertar oferta.", {
        codigoCarrera: carreraNuevaOferta[0].codigo,
        codigoCuatrimestre: cuatrimestreNuevaOferta[0].codigo,
        error: e
      });

      error(500, "No pudimos guardar tu oferta. Intentá nuevamente en unos minutos.");
    }

    console.info("[enviarOferta] Oferta insertada correctamente", {
      codigoCarrera: carreraNuevaOferta[0].codigo,
      codigoCuatrimestre: cuatrimestreNuevaOferta[0].codigo
    });

    return { success: true };
  }
);
