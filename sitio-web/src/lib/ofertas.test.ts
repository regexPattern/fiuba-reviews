import { describe, expect, it } from "vitest";
import { normalizarNombre, ordenCuatrimestre } from "./ofertas";

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
