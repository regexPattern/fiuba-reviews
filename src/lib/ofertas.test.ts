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

  it("reconoce el cuatrimestre del filtro", () => {
    expect(extraerMetadataOferta(encabezado("2025 - 2do Cuatrimestre"))?.cuatrimestre).toEqual({
      numero: 2,
      anio: 2025
    });
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
