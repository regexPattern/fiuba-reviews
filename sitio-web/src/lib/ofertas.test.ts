import { describe, expect, it } from "vitest";
import { extraerMetadataOferta, normalizarNombre, ordenCuatrimestre } from "./ofertas";

/** Encabezado real de la página del SIU (con el nombre del alumno inventado). */
const encabezado = (periodoFiltro: string) => `
 0
 NOMBRE APELLIDO
Propuesta: Ingeniería en Informática
 70: 25/02/2026 14:30
Inscripción a Materias
Reportes
Perfil: Alumno
Oferta de comisiones
Filtros
LunesMartesMiércolesJuevesViernesSábadoDomingoSin definir
Período lectivo: 
${periodoFiltro}
Actividad: 
Todas
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: ÁLGEBRA LINEAL (CB002)
Comisión: CURSO: 01A
`;

describe("extraerMetadataOferta", () => {
  it("reconoce la carrera y el cuatrimestre con el filtro en Todos", () => {
    expect(extraerMetadataOferta(encabezado("Todos"))).toEqual({
      carrera: "Ingeniería en Informática",
      cuatrimestre: { numero: 1, anio: 2026 }
    });
  });

  it("toma el primer período lectivo si hay varios", () => {
    const contenido = `${encabezado("Todos")}Período lectivo: 2026 - 2do Cuatrimestre\n`;

    expect(extraerMetadataOferta(contenido)?.cuatrimestre).toEqual({ numero: 1, anio: 2026 });
  });

  it("devuelve null con contenido que no es del SIU", () => {
    expect(extraerMetadataOferta("")).toBeNull();
    expect(extraerMetadataOferta("hola")).toBeNull();
  });
});

describe("normalizarNombre", () => {
  it("ignora tildes, mayúsculas y espacios", () => {
    expect(normalizarNombre("  Ingeniería en  Informática ")).toBe(
      normalizarNombre("INGENIERIA EN INFORMATICA")
    );
  });
});

describe("ordenCuatrimestre", () => {
  it("ordena cronológicamente", () => {
    expect(ordenCuatrimestre({ numero: 2, anio: 2025 })).toBeLessThan(
      ordenCuatrimestre({ numero: 1, anio: 2026 })
    );
  });
});
