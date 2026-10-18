export const CRITERIOS = [
  { name: "aceptaCritica", label: "Acepta crítica" },
  { name: "asistencia", label: "Asistencia" },
  { name: "buenTrato", label: "Buen trato" },
  { name: "claridad", label: "Claridad" },
  { name: "claseOrganizada", label: "Clase organizada" },
  { name: "cumpleHorarios", label: "Cumple horarios" },
  { name: "fomentaParticipacion", label: "Fomenta participación" },
  { name: "panoramaAmplio", label: "Panorama amplio" },
  { name: "respondeMails", label: "Responde mails" }
] as const;

export type CriterioName = (typeof CRITERIOS)[number]["name"];

export const CALIFICACIONES_VALIDAS = [0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5] as const;
