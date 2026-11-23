/**
 * Tags del caché. Todas las queries se cachean con `cacheLife("max")` y solo se invalidan
 * por tag, cuando cambian los datos.
 *
 *  - `materias`: listado de materias vigentes y sus equivalencias (cambia con los planes).
 *  - `materia:<codigo>`: info de una materia vigente (nombre, última actualización).
 *  - `catedras:<codigo>`: cátedras propias de una materia y sus calificaciones. Una materia
 *    vigente que usa equivalencias depende de los `catedras:` de cada equivalencia.
 *  - `catedra:<uuid>`: detalle de una cátedra (docentes, promedios y comentarios).
 */
export const tags = {
  materias: "materias",
  materia: (codigo: string) => `materia:${codigo}`,
  catedras: (codigoMateria: string) => `catedras:${codigoMateria}`,
  catedra: (codigo: string) => `catedra:${codigo}`
} as const;
