export type Cuatrimestre = { numero: number; anio: number };

/**
 * Encabezado de la página "Oferta de comisiones" del SIU copiada entera: la carrera sale de
 * `Propuesta: <carrera>` (seguida de `<n>: dd/mm/aaaa`) y el cuatrimestre del primer
 * `Período lectivo: aaaa - 1er|2do Cuatrimestre` (el filtro puede decir `Todos`).
 */
const METADATA_RE =
  /Propuesta:\s*([^\r\n]*?)(?=\s+\d+\s*:\s*\d{2}\/\d{2}\/\d{4})[\s\S]*?per[ií]odo lectivo:\s+(\d{4}).*?(\d)(?:er|do)/i;

/** Carrera y cuatrimestre de una oferta pegada desde el SIU, o `null` si no se reconoce. */
export function extraerMetadataOferta(contenido: string) {
  const matches = METADATA_RE.exec(contenido);

  if (!matches) return null;

  const carrera = matches[1].trim();
  const cuatrimestre = { anio: Number(matches[2]), numero: Number(matches[3]) };

  if (!carrera || (cuatrimestre.numero !== 1 && cuatrimestre.numero !== 2)) return null;

  return { carrera, cuatrimestre };
}

/** Para comparar nombres de carrera sin importar tildes, mayúsculas ni espacios. */
export function normalizarNombre(nombre: string) {
  return nombre
    .normalize("NFD")
    .replace(/\p{Diacritic}/gu, "")
    .replace(/\s+/g, " ")
    .trim()
    .toLowerCase();
}

/** Orden cronológico de un cuatrimestre, para compararlos. */
export function ordenCuatrimestre({ numero, anio }: Cuatrimestre) {
  return anio * 2 + numero;
}

export function formatearCuatrimestre({ numero, anio }: Cuatrimestre) {
  return `${numero}C${anio}`;
}
