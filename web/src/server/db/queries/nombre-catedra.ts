type DocenteParaNombre = { nombre: string; prioridad: number };

export const PRIORIDAD_SIN_ROL = Number.MAX_SAFE_INTEGER;

export function calcularNombreCatedra(docentes: DocenteParaNombre[]) {
  return [...docentes]
    .sort((a, b) => a.prioridad - b.prioridad || a.nombre.localeCompare(b.nombre))
    .map((docente) => docente.nombre)
    .join("-");
}
