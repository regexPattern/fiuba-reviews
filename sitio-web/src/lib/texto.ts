/** Para comparar textos ignorando tildes, mayúsculas y espacios de más. */
export function normalizarTexto(texto: string) {
  return texto.normalize("NFD").replace(/\p{M}/gu, "").replace(/\s+/g, " ").trim().toLowerCase();
}
