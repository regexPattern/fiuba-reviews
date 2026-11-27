import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";
import { parsearMetadata, parsearOfertas } from "./parser";

describe("ofertas reales", () => {
  it.each([
    ["civil-1C-2026", "Ingeniería Civil", 2026, 1, 73, 183, 407],
    ["industrial-1C-2026", "Ingeniería Industrial", 2026, 1, 57, 231, 486],
    ["informatica-1C-2026", "Ingeniería en Informática", 2026, 1, 41, 166, 359],
    ["informatica-2C-2026", "Ingeniería en Informática", 2026, 2, 43, 166, 323],
    ["mecanica-1C-2026", "Ingeniería Mecánica", 2026, 1, 61, 246, 461]
  ])("%s", (name, carrera, anio, cuatrimestre, materias, comisiones, clases) => {
    const texto = readFileSync(new URL(`./fixtures/${name}.txt`, import.meta.url), "utf8");

    const ofertas = parsearOfertas(texto);
    const resumen = ofertas.map((oferta) => {
      const comisiones = oferta.materias.flatMap((m) => m.comisiones);
      return {
        carrera: oferta.carrera,
        anio: oferta.cuatrimestre.anio,
        cuatrimestre: oferta.cuatrimestre.numero,
        materias: oferta.materias.length,
        comisiones: comisiones.length,
        clases: comisiones.reduce((sum, c) => sum + c.clases.length, 0)
      };
    });

    expect(resumen).toEqual([{ carrera, anio, cuatrimestre, materias, comisiones, clases }]);
  });

  it("petroleo-2C-2025-y-1C-2026", () => {
    const texto = readFileSync(
      new URL("./fixtures/petroleo-2C-2025-y-1C-2026.txt", import.meta.url),
      "utf8"
    );
    const carrera = "Ingeniería en Petróleo";

    const ofertas = parsearOfertas(texto);
    const resumen = ofertas.map((oferta) => {
      const comisiones = oferta.materias.flatMap((m) => m.comisiones);
      return {
        carrera: oferta.carrera,
        anio: oferta.cuatrimestre.anio,
        cuatrimestre: oferta.cuatrimestre.numero,
        materias: oferta.materias.length,
        comisiones: comisiones.length,
        clases: comisiones.reduce((sum, c) => sum + c.clases.length, 0)
      };
    });
    const metadata = parsearMetadata(texto);

    expect(resumen).toEqual([
      { carrera, anio: 2025, cuatrimestre: 2, materias: 44, comisiones: 160, clases: 358 },
      { carrera, anio: 2026, cuatrimestre: 1, materias: 46, comisiones: 179, clases: 394 }
    ]);
    expect(metadata).toEqual([
      { carrera, cuatrimestre: { numero: 2, anio: 2025 } },
      { carrera, cuatrimestre: { numero: 1, anio: 2026 } }
    ]);
  });

  it("parsea una comisión con varias clases", () => {
    const texto = readFileSync(
      new URL("./fixtures/informatica-1C-2026.txt", import.meta.url),
      "utf8"
    );

    const [oferta] = parsearOfertas(texto);
    const algebra = oferta?.materias.find((m) => m.codigo === "CB002");

    expect(algebra?.nombre).toBe("ÁLGEBRA LINEAL");
    expect(algebra?.comisiones[0]).toEqual({
      nombre: "CURSO: 01A",
      docentes: [
        { nombre: "VARGAS GABRIELA", rol: "Profesor/a Adjunto/a" },
        { nombre: "PALACIOS PATRICIA ELIZABETH", rol: "Jefe/a Trabajos Practicos" },
        { nombre: "SOSA MARINA", rol: "Jefe/a Trabajos Practicos" }
      ],
      clases: [
        { tipo: "Teórica", dia: "Lunes", inicio: "09:00", fin: "11:00", aula: "403 - PC" },
        { tipo: "Práctica", dia: "Lunes", inicio: "11:00", fin: "13:00", aula: "403 - PC" },
        { tipo: "Teórica", dia: "Miércoles", inicio: "09:00", fin: "11:00", aula: "403 - PC" },
        { tipo: "Práctica", dia: "Miércoles", inicio: "11:00", fin: "13:00", aula: "403 - PC" }
      ]
    });
  });

  it("parsea una comisión sin horario definido (formato en una línea)", () => {
    const texto = readFileSync(new URL("./fixtures/mecanica-1C-2026.txt", import.meta.url), "utf8");

    const [oferta] = parsearOfertas(texto);
    const algebra = oferta?.materias.find((m) => m.codigo === "CB002");

    expect(algebra?.comisiones[0]).toEqual({
      nombre: "CONDICIONALES",
      docentes: [
        { nombre: "GIGOLA SILVIA VIVIANA", rol: "Profesor/a Adjunto/a" },
        { nombre: "GRYNBERG SEBASTIAN PABLO", rol: "Profesor/a Adjunto/a" },
        { nombre: "VARGAS GABRIELA", rol: "Profesor/a Adjunto/a" }
      ],
      clases: []
    });
  });

  it("parsea una comisión sin docentes", () => {
    const texto = readFileSync(new URL("./fixtures/civil-1C-2026.txt", import.meta.url), "utf8");

    const [oferta] = parsearOfertas(texto);
    const materia = oferta?.materias.find((m) => m.codigo === "TA040");

    expect(materia?.comisiones[0]).toEqual({ nombre: "CONDICIONALES", docentes: [], clases: [] });
  });

  it("ignora las observaciones de las comisiones", () => {
    const texto = readFileSync(new URL("./fixtures/civil-1C-2026.txt", import.meta.url), "utf8");

    const [oferta] = parsearOfertas(texto);
    const comisiones = oferta!.materias.flatMap((m) => m.comisiones);

    expect(JSON.stringify(comisiones)).not.toContain("Observaciones");
  });
});
