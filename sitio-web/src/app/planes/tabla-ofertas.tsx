import type { EstadoOferta, OfertaCarrera } from "@/server/db/queries/ofertas";
import { type Cuatrimestre, formatearCuatrimestre } from "@/lib/cuatrimestres";

const ETIQUETAS: Record<EstadoOferta, { texto: string; clases: string }> = {
  "al-dia": { texto: "Al día", clases: "bg-emerald-50 text-emerald-700 border-emerald-200" },
  "en-proceso": { texto: "En proceso", clases: "bg-sky-50 text-fiuba-oscuro border-sky-200" },
  desactualizada: {
    texto: "Desactualizada",
    clases: "bg-amber-50 text-amber-800 border-amber-200"
  },
  "sin-oferta": {
    texto: "Sin oferta",
    clases: "border-dashed border-stone-300 text-stone-500"
  }
};

export function TablaOfertas({
  ofertas,
  ultimoCuatrimestre
}: {
  ofertas: OfertaCarrera[];
  ultimoCuatrimestre: Cuatrimestre | null;
}) {
  const alDia = ofertas.filter((o) => o.estado === "al-dia").length;

  return (
    <div className="overflow-hidden relieve bg-white">
      <table className="w-full text-left text-sm">
        <caption className="border-b border-stone-200 px-4 py-3 text-left text-sm text-stone-500">
          <span className="font-mono tabular-nums">{alDia}</span> de{" "}
          <span className="font-mono tabular-nums">{ofertas.length}</span> carreras al día
          {ultimoCuatrimestre && (
            <>
              {" "}
              con el{" "}
              <span className="text-stone-700">
                <span className="font-mono tabular-nums">
                  {formatearCuatrimestre(ultimoCuatrimestre)}
                </span>
              </span>
            </>
          )}
          .
        </caption>
        <thead className="border-b border-stone-200 bg-stone-50 rotulo text-stone-500">
          <tr>
            <th scope="col" className="px-4 py-2 font-normal">
              Carrera
            </th>
            <th scope="col" className="px-4 py-2 text-right font-normal">
              Última oferta
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-stone-200">
          {ofertas.map((oferta) => (
            <Fila key={oferta.codigo} oferta={oferta} />
          ))}
        </tbody>
      </table>
    </div>
  );
}

function Fila({ oferta }: { oferta: OfertaCarrera }) {
  const etiqueta = ETIQUETAS[oferta.estado];

  return (
    <tr>
      <th scope="row" className="px-4 py-2.5 font-normal">
        <span className="block text-base/6 tracking-tight text-stone-900">{oferta.nombre}</span>
        {oferta.planVigente && (
          <span className="block text-xs text-stone-500">
            Plan <span className="font-mono tabular-nums">{oferta.planVigente}</span>
          </span>
        )}
      </th>
      {/* oxlint-disable-next-line jsx-a11y/control-has-associated-label -- Es una celda de datos, no un control: el texto es dinámico. */}
      <td className="px-4 py-2.5 text-right align-top">
        <span className="flex flex-col items-end gap-1 sm:flex-row sm:items-center sm:justify-end sm:gap-2">
          <span className="font-mono text-stone-700 tabular-nums">
            {oferta.procesada ? formatearCuatrimestre(oferta.procesada) : "–"}
          </span>
          <span
            className={`rounded-full border-[1.5px] px-2 py-px text-xs whitespace-nowrap ${etiqueta.clases}`}
            title={
              oferta.recibida
                ? `Recibimos la oferta del ${formatearCuatrimestre(oferta.recibida)} y se está procesando.`
                : undefined
            }
          >
            {etiqueta.texto}
            {oferta.recibida && (
              <>
                {" "}
                <span className="font-mono tabular-nums">
                  {formatearCuatrimestre(oferta.recibida)}
                </span>
              </>
            )}
          </span>
        </span>
      </td>
    </tr>
  );
}
