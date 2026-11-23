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

  it("ignora tildes y tolera errores de tipeo", () => {
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
});
