/** Normaliza el texto removiendo tildes, mayúsculas y espacios de más. */
export function normalizar(texto: string) {
  return texto.normalize("NFD").replace(/\p{M}/gu, "").replace(/\s+/g, " ").trim().toLowerCase();
}

const FORMATO_PUNTAJE = new Intl.NumberFormat("es-AR", {
  minimumFractionDigits: 1,
  maximumFractionDigits: 1
});

/** Formatea un puntaje con coma decimal y siempre un decimal. */
export function formatearPuntaje(puntaje: number) {
  return FORMATO_PUNTAJE.format(puntaje);
}
