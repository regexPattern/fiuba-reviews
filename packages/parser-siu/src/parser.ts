import type { Clase, Docente, Metadata, Oferta } from "./tipos";

const DIAS = ["Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"];

/**
 * Una oferta por cada período lectivo de la página, en el orden en que aparecen (con el filtro
 * en `Todos` puede haber varios). Vacío si el texto no es una oferta del SIU.
 */
export function parsearOfertas(text: string): Oferta[] {
  const lines = text.split(/\r?\n/).map((line) => line.trim());
  const carrera = lines.map((line) => line.match(/^Propuesta: (.+)$/)?.[1]).find(Boolean);
  if (!carrera) return [];

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
      const numero = Number(periodo[2]) === 1 ? 1 : 2;
      ofertas.push({ carrera, cuatrimestre: { numero, anio: Number(periodo[1]) }, materias: [] });
    } else if (actividad && oferta) {
      oferta.materias.push({ nombre: actividad[1]!, codigo: actividad[2]!, comisiones: [] });
    } else if (comision && materia) {
      materia.comisiones.push({ nombre: comision[1]!, docentes: [], clases: [] });
    } else if (docentes && ultimaComision) {
      ultimaComision.docentes = parsearDocentes(docentes[1]!);
    } else if (ultimaComision) {
      const clase = parsearClase(line);
      if (clase) ultimaComision.clases.push(clase);
    }
  }

  return ofertas;
}

/** Como `parsearOfertas`, pero solo la carrera y el cuatrimestre de cada período. */
export function parsearMetadata(text: string): Metadata[] {
  return parsearOfertas(text).map(({ carrera, cuatrimestre }) => ({ carrera, cuatrimestre }));
}

function parsearDocentes(text: string): Docente[] {
  if (text === "Sin docentes") return [];
  return text.split(", ").map((docente) => {
    const [nombre, rol] = docente.replace(/\)$/, "").split(" (");
    return { nombre: nombre!.trim(), rol: rol ?? null };
  });
}

function parsearClase(line: string): Clase | null {
  const [tipo, dia, horario, aula] = line.split("\t").map((field) => field.trim());
  if (!tipo || !dia || !horario || !DIAS.includes(dia)) return null;

  const [inicio, fin] = horario.split(" a ");
  if (!inicio || !fin) return null;

  return { tipo, dia, inicio, fin, aula: aula || null };
}
