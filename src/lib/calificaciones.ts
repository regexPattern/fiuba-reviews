export const CRITERIOS = [
  { nombre: "aceptaCritica", etiqueta: "Acepta crítica" },
  { nombre: "asistencia", etiqueta: "Asistencia" },
  { nombre: "buenTrato", etiqueta: "Buen trato" },
  { nombre: "claridad", etiqueta: "Claridad" },
  { nombre: "claseOrganizada", etiqueta: "Clase organizada" },
  { nombre: "cumpleHorarios", etiqueta: "Cumple horarios" },
  { nombre: "fomentaParticipacion", etiqueta: "Fomenta participación" },
  { nombre: "panoramaAmplio", etiqueta: "Panorama amplio" },
  { nombre: "respondeMails", etiqueta: "Responde mails" }
] as const;

export type NivelPuntaje = "muyBueno" | "bueno" | "regular" | "malo";

export function nivelPuntaje(puntaje: number): NivelPuntaje {
  if (puntaje >= 4.5) return "muyBueno";
  if (puntaje >= 3.5) return "bueno";
  if (puntaje >= 2.5) return "regular";
  return "malo";
}
