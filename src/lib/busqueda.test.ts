import { describe, expect, it } from "vitest";
import { crearBusqueda } from "./busqueda";

describe("búsqueda de cátedras", () => {
  const catedras = [
    { nombre: "Comas-Canepa-Medina" },
    { nombre: "López-Canga-Silvestre-Vicedo-Vorobioff-Zhang" },
    { nombre: "Acero-López-García-Endelli-Portocarrero" },
    { nombre: "González-Amoreo-Petrovich" }
  ];
  const buscar = crearBusqueda(catedras, ["nombre"]);
  const nombres = (query: string) => buscar(query).map(({ item }) => item.nombre);

  it("sin búsqueda devuelve todas en el orden original", () => {
    expect(nombres("")).toEqual(catedras.map((catedra) => catedra.nombre));
  });

  it("ignora tildes y mayúsculas", () => {
    expect(nombres("GONZALEZ")).toEqual(["González-Amoreo-Petrovich"]);
  });

  it("tolera errores de tipeo", () => {
    expect(nombres("gonsalez")).toEqual(["González-Amoreo-Petrovich"]);
  });

  it("encuentra las palabras en cualquier orden", () => {
    expect(nombres("petrovich gonzalez")).toEqual(["González-Amoreo-Petrovich"]);
  });

  it("tiene que encontrar todas las palabras", () => {
    expect(nombres("lopez zhang")).toEqual(["López-Canga-Silvestre-Vicedo-Vorobioff-Zhang"]);
    expect(nombres("zhang endelli")).toEqual([]);
  });
});

describe("búsqueda de materias", () => {
  const materias = [
    { codigo: "CB001", nombre: "Análisis Matemático II" },
    { codigo: "CB002", nombre: "Álgebra Lineal" },
    { codigo: "CB100", nombre: "Algoritmos y Estructuras de Datos" }
  ];
  const buscar = crearBusqueda(materias, ["codigo", "nombre"]);
  const codigos = (query: string) => buscar(query).map(({ item }) => item.codigo);

  it("encuentra por código", () => {
    expect(codigos("cb100")[0]).toBe("CB100");
  });

  it("encuentra por nombre", () => {
    expect(codigos("algebra")).toEqual(["CB002"]);
    expect(codigos("datos algoritmos")).toEqual(["CB100"]);
  });

  it("encuentra mientras se escribe, aunque la última palabra tenga una sola letra", () => {
    expect(codigos("analisis m")).toEqual(["CB001"]);
  });

  it("entre las coincidencias exactas respeta el orden original", () => {
    const buscar = crearBusqueda(
      [
        { codigo: "TA001", nombre: "Criptoanálisis" },
        { codigo: "CB001", nombre: "Análisis Matemático II" },
        { codigo: "CB005", nombre: "Análisis Numérico" }
      ],
      ["codigo", "nombre"]
    );
    expect(buscar("analisis").map(({ item }) => item.codigo)).toEqual(["TA001", "CB001", "CB005"]);
  });

  it("pone las coincidencias exactas antes que las aproximadas", () => {
    const buscar = crearBusqueda(
      [{ nombre: "Gonzales-Pérez" }, { nombre: "González-Amoreo" }, { nombre: "Fernández" }],
      ["nombre"]
    );
    // "Gonzales" está antes en la lista, pero es solo parecido.
    expect(buscar("gonzalez").map(({ item }) => item.nombre)).toEqual([
      "González-Amoreo",
      "Gonzales-Pérez"
    ]);
  });

  it("devuelve los rangos a subrayar sobre el texto original", () => {
    const [resultado] = buscar("mate anal");
    expect(resultado.matches?.map(({ key, indices }) => ({ key, indices }))).toEqual([
      {
        key: "nombre",
        indices: [
          [0, 3],
          [9, 12]
        ]
      }
    ]);
  });
});
