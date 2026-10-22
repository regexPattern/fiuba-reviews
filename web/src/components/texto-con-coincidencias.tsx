export type Rangos = readonly (readonly [number, number])[];
export function TextoConCoincidencias({
  texto,
  rangos = [],
  className
}: {
  texto: string;
  rangos?: Rangos;
  className?: string;
}) {
  const unidos: [number, number][] = [];
  for (const [inicio, fin] of [...rangos].sort((a, b) => a[0] - b[0])) {
    const ultimo = unidos[unidos.length - 1];
    if (ultimo && inicio <= ultimo[1] + 1) {
      ultimo[1] = Math.max(ultimo[1], fin);
    } else {
      unidos.push([inicio, fin]);
    }
  }

  let cursor = 0;
  return (
    <span className={className}>
      {unidos.map(([inicio, fin]) => {
        const anterior = texto.slice(cursor, inicio);
        cursor = fin + 1;
        return (
          <span key={inicio}>
            {anterior}
            <span className="font-medium underline decoration-2 underline-offset-2">
              {texto.slice(inicio, fin + 1)}
            </span>
          </span>
        );
      })}
      {texto.slice(cursor)}
    </span>
  );
}

const normalizar = (caracter: string) =>
  caracter.normalize("NFD").replace(/\p{M}/gu, "").toLowerCase()[0] ?? caracter;

/** Dentro de los rangos aproximados de Fuse, deja solo los caracteres que realmente coinciden con la búsqueda (subsecuencia común más larga). */
export function rangosExactos(
  texto: string,
  busqueda: string,
  rangos: Rangos = []
): [number, number][] {
  const palabra = [...busqueda].filter((c) => !/\s/.test(c)).map(normalizar);
  const resultado: [number, number][] = [];

  for (const [inicio, fin] of rangos) {
    const sub = Array.from(texto.slice(inicio, fin + 1), normalizar);
    const n = palabra.length;
    const m = sub.length;
    const tabla = Array.from({ length: n + 1 }, () => Array.from({ length: m + 1 }, () => 0));
    for (let i = 1; i <= n; i++) {
      for (let j = 1; j <= m; j++) {
        tabla[i]![j] =
          palabra[i - 1] === sub[j - 1]
            ? tabla[i - 1]![j - 1]! + 1
            : Math.max(tabla[i - 1]![j]!, tabla[i]![j - 1]!);
      }
    }
    let i = n;
    let j = m;
    while (i > 0 && j > 0) {
      if (palabra[i - 1] === sub[j - 1]) {
        resultado.push([inicio + j - 1, inicio + j - 1]);
        i--;
        j--;
      } else if (tabla[i - 1]![j]! >= tabla[i]![j - 1]!) i--;
      else j--;
    }
  }
  return resultado;
}
