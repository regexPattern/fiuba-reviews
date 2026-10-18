import type { Clase, Docente, Metadata, Oferta } from "./types.js";

export type { Clase, Comision, Docente, Materia, Metadata, Oferta } from "./types.js";

const DIAS = ["Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"];

export function parseOfertas(text: string): Oferta[] {
  const lines = text.split(/\r?\n/).map((line) => line.trim());
  const carrera = lines.map((line) => line.match(/^Propuesta: (.+)$/)?.[1]).find(Boolean);
  if (!carrera) throw new Error("No se encontró la carrera (Propuesta) en el texto del SIU");

  const ofertas: Oferta[] = [];

  for (const line of lines) {
    const periodo = line.match(/^Período lectivo: (\d{4}) - (\d)/);
    const actividad = line.match(/^Actividad: (.+) \((.+)\)$/);
    const comision = line.match(/^Comisión: (.+)$/);
    const docentes = line.match(/^(?:Turno:.*? )?Docentes: (.+)$/);

    const oferta = ofertas.at(-1);
    const materia = oferta?.materias.at(-1);
    const ultimaComision = materia?.comisiones.at(-1);

    if (periodo) {
      const cuatrimestre = Number(periodo[2]) === 1 ? 1 : 2;
      ofertas.push({ carrera, anio: Number(periodo[1]), cuatrimestre, materias: [] });
    } else if (actividad && oferta) {
      oferta.materias.push({ nombre: actividad[1]!, codigo: actividad[2]!, comisiones: [] });
    } else if (comision && materia) {
      materia.comisiones.push({ nombre: comision[1]!, docentes: [], clases: [] });
    } else if (docentes && ultimaComision) {
      ultimaComision.docentes = parseDocentes(docentes[1]!);
    } else if (ultimaComision) {
      const clase = parseClase(line);
      if (clase) ultimaComision.clases.push(clase);
    }
  }

  if (ofertas.length === 0)
    throw new Error("No se encontró el período lectivo en el texto del SIU");
  return ofertas;
}

export function parseMetadata(text: string): Metadata[] {
  return parseOfertas(text).map(({ carrera, anio, cuatrimestre }) => ({
    carrera,
    anio,
    cuatrimestre
  }));
}

function parseDocentes(text: string): Docente[] {
  if (text === "Sin docentes") return [];
  return text.split(", ").map((docente) => {
    const [nombre, rol] = docente.replace(/\)$/, "").split(" (");
    return { nombre: nombre!.trim(), rol: rol ?? null };
  });
}

function parseClase(line: string): Clase | null {
  const [tipo, dia, horario, aula] = line.split("\t").map((field) => field.trim());
  if (!tipo || !dia || !horario || !DIAS.includes(dia)) return null;

  const [inicio, fin] = horario.split(" a ");
  if (!inicio || !fin) return null;

  return { tipo, dia, inicio, fin, aula: aula || null };
}
