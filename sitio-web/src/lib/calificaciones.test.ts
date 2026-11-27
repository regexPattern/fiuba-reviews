import { describe, expect, it } from "vitest";
import {
  BORRADOR_VACIO,
  CRITERIOS,
  LARGO_MAXIMO_COMENTARIO,
  borradorCompleto,
  borradorTocado,
  criteriosFaltantes,
  cuatrimestresRecientes,
  formatearPuntaje,
  PUNTAJES,
  type BorradorCalificacion
} from "./calificaciones";

const completo: BorradorCalificacion = {
  puntajes: Object.fromEntries(CRITERIOS.map(({ nombre }) => [nombre, 4])),
  cuatrimestre: { numero: 1, anio: 2026 },
  comentario: ""
};

describe("PUNTAJES", () => {
  it("va de 0,5 a 5 de a medio punto", () => {
    expect(PUNTAJES).toEqual([0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5]);
  });
});

describe("formatearPuntaje", () => {
  it("usa coma decimal y omite el ,0", () => {
    expect(formatearPuntaje(4.5)).toBe("4,5");
    expect(formatearPuntaje(4)).toBe("4");
  });
});

describe("criteriosFaltantes", () => {
  it("devuelve todos los criterios de un borrador vacío, en orden", () => {
    expect(criteriosFaltantes(BORRADOR_VACIO)).toEqual(CRITERIOS.map(({ nombre }) => nombre));
  });

  it("devuelve solo los que no tienen puntaje", () => {
    const { claridad: _, ...resto } = completo.puntajes;
    expect(criteriosFaltantes({ ...completo, puntajes: resto })).toEqual(["claridad"]);
  });
});

describe("borradorCompleto", () => {
  it("pide todos los criterios y el cuatrimestre, pero no el comentario", () => {
    expect(borradorCompleto(completo)).toBe(true);
    expect(borradorCompleto({ ...completo, cuatrimestre: null })).toBe(false);
    expect(borradorCompleto({ ...completo, puntajes: { claridad: 5 } })).toBe(false);
  });

  it("rechaza comentarios demasiado largos", () => {
    const comentario = "a".repeat(LARGO_MAXIMO_COMENTARIO + 1);
    expect(borradorCompleto({ ...completo, comentario })).toBe(false);
  });
});

describe("borradorTocado", () => {
  it("ignora un comentario de solo espacios", () => {
    expect(borradorTocado(BORRADOR_VACIO)).toBe(false);
    expect(borradorTocado({ ...BORRADOR_VACIO, comentario: "  \n" })).toBe(false);
    expect(borradorTocado({ ...BORRADOR_VACIO, puntajes: { asistencia: 3 } })).toBe(true);
  });
});

describe("cuatrimestresRecientes", () => {
  it("arranca por el 2C de agosto en adelante", () => {
    expect(cuatrimestresRecientes(new Date(2026, 7, 1), 3)).toEqual([
      { numero: 2, anio: 2026 },
      { numero: 1, anio: 2026 },
      { numero: 2, anio: 2025 }
    ]);
  });

  it("arranca por el 1C de marzo a julio", () => {
    expect(cuatrimestresRecientes(new Date(2026, 2, 1), 2)).toEqual([
      { numero: 1, anio: 2026 },
      { numero: 2, anio: 2025 }
    ]);
  });

  it("en el verano arranca por el 2C del año anterior", () => {
    expect(cuatrimestresRecientes(new Date(2026, 0, 15), 1)).toEqual([{ numero: 2, anio: 2025 }]);
  });
});
