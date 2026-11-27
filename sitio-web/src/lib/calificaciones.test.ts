import { describe, expect, it } from "vitest";
import {
  BORRADOR_VACIO,
  CRITERIOS,
  LARGO_MAXIMO_COMENTARIO,
  borradorCompleto,
  borradorTocado,
  criteriosFaltantes,
  formatearPuntaje,
  PUNTAJES,
  validarCalificacion,
  type BorradorCalificacion
} from "./calificaciones";

const completo: BorradorCalificacion = {
  puntajes: Object.fromEntries(CRITERIOS.map(({ nombre }) => [nombre, 4])),
  cuatrimestre: { numero: 1, anio: 2026 },
  comentario: ""
};

describe("PUNTAJES", () => {
  it("va de 1 a 5 de a medio punto", () => {
    expect(PUNTAJES).toEqual([1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5]);
  });
});

describe("formatearPuntaje", () => {
  it("usa coma decimal y siempre un decimal", () => {
    expect(formatearPuntaje(4.5)).toBe("4,5");
    expect(formatearPuntaje(4)).toBe("4,0");
    expect(formatearPuntaje(3.66)).toBe("3,7");
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

describe("validarCalificacion", () => {
  const permitidos = [
    { numero: 1, anio: 2026 },
    { numero: 2, anio: 2025 }
  ];
  const envio = { ...completo, codigoDocente: "5f3c1b2e-8a4d-4e6f-9b7a-1c2d3e4f5a6b" };

  it("acepta un envío completo y recorta el comentario", () => {
    const resultado = validarCalificacion({ ...envio, comentario: "  Muy bueno \n" }, permitidos);
    expect(resultado).toEqual({
      ok: true,
      calificacion: {
        codigoDocente: envio.codigoDocente,
        puntajes: envio.puntajes,
        cuatrimestre: { numero: 1, anio: 2026 },
        comentario: "Muy bueno"
      }
    });
  });

  it("guarda un comentario vacío como null", () => {
    const resultado = validarCalificacion({ ...envio, comentario: "   " }, permitidos);
    expect(resultado.ok && resultado.calificacion.comentario).toBeNull();
  });

  it("rechaza criterios faltantes o puntajes fuera de la escala", () => {
    const { claridad: _, ...resto } = envio.puntajes;
    expect(validarCalificacion({ ...envio, puntajes: resto }, permitidos).ok).toBe(false);
    for (const puntaje of [0.5, 4.2, 6]) {
      const puntajes = { ...envio.puntajes, claridad: puntaje };
      expect(validarCalificacion({ ...envio, puntajes }, permitidos).ok).toBe(false);
    }
  });

  it("rechaza un cuatrimestre que no se ofrece", () => {
    const cuatrimestre = { numero: 1, anio: 2025 };
    expect(validarCalificacion({ ...envio, cuatrimestre }, permitidos)).toEqual({
      ok: false,
      mensaje: "No se puede calificar el 1C2025."
    });
  });

  it("rechaza un docente inválido o un comentario demasiado largo", () => {
    expect(validarCalificacion({ ...envio, codigoDocente: "x" }, permitidos).ok).toBe(false);
    const comentario = "a".repeat(LARGO_MAXIMO_COMENTARIO + 1);
    expect(validarCalificacion({ ...envio, comentario }, permitidos).ok).toBe(false);
  });
});
