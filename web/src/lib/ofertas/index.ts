const METADATA_RE =
  /Propuesta:\s*([^\r\n]*?)(?=\s+\d+\s*:\s*\d{2}\/\d{2}\/\d{4})[\s\S]*?per[ií]odo lectivo:\s+(\d{4}).*?(\d)(?:er|do)/i;

export function extraerMetadataOferta(contenido: string) {
  if (contenido === "") {
    return null;
  }

  const matches = METADATA_RE.exec(contenido);

  if (matches === null) {
    return null;
  }

  const carrera = matches[1];
  const anio = parseInt(matches[2], 10);
  const numero = parseInt(matches[3], 10);

  return { carrera, cuatrimestre: { anio, numero } };
}
