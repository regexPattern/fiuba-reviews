import { describe, expect, it } from "vitest";
import {
  CRITERIOS,
  criteriosFaltantes,
  formatearPuntaje,
  schemaCalificacion,
  type BorradorCalificacion
} from "./calificaciones";

const completo: BorradorCalificacion = {
  puntajes: Object.fromEntries(CRITERIOS.map(({ nombre }) => [nombre, 4])),
  cuatrimestre: { numero: 1, anio: 2026 },
  comentario: ""
};

describe("formatearPuntaje", () => {
  it("usa coma decimal y siempre un decimal", () => {
    expect(formatearPuntaje(4.5)).toBe("4,5");
    expect(formatearPuntaje(4)).toBe("4,0");
    expect(formatearPuntaje(3.66)).toBe("3,7");
  });
});

describe("criteriosFaltantes", () => {
  it("devuelve solo los que no tienen puntaje", () => {
    const { claridad: _, ...resto } = completo.puntajes;
    expect(criteriosFaltantes({ ...completo, puntajes: resto })).toEqual(["claridad"]);
  });
});

describe("schemaCalificacion", () => {
  const schema = schemaCalificacion([
    { codigo: 10, numero: 1, anio: 2026 },
    { codigo: 9, numero: 2, anio: 2025 }
  ]);
  const envio = { ...completo, codigoDocente: "5f3c1b2e-8a4d-4e6f-9b7a-1c2d3e4f5a6b" };

  it("acepta un envío completo, recorta el comentario y devuelve el cuatrimestre permitido", () => {
    const resultado = schema.safeParse({ ...envio, comentario: "  Muy bueno \n" });
    expect(resultado.data).toEqual({
      codigoDocente: envio.codigoDocente,
      puntajes: envio.puntajes,
      cuatrimestre: { codigo: 10, numero: 1, anio: 2026 },
      comentario: "Muy bueno"
    });
  });

  it("guarda un comentario vacío como null", () => {
    expect(schema.safeParse({ ...envio, comentario: "   " }).data?.comentario).toBeNull();
  });

  it("rechaza puntajes fuera de la escala de media estrella", () => {
    for (const puntaje of [0.5, 4.2, 6]) {
      const puntajes = { ...envio.puntajes, claridad: puntaje };
      expect(schema.safeParse({ ...envio, puntajes }).success).toBe(false);
    }
  });

  it("rechaza un cuatrimestre que no se ofrece, con un mensaje para mostrar", () => {
    const cuatrimestre = { numero: 1, anio: 2025 };
    expect(schema.safeParse({ ...envio, cuatrimestre }).error?.issues).toMatchObject([
      { message: "No se puede calificar el 1C2025." }
    ]);
  });
});
