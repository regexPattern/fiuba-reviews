"use client";

import type { OfertaCarrera } from "@/server/db/queries/ofertas";
import type { ReactNode } from "react";
import { LoaderCircle, Send } from "@withicons/react/line";
import {
  AlertCircle,
  AlertTriangle,
  CircleCheckBig,
  Handshake,
  PartyPopper,
  Paste
} from "@withicons/react/plush";
import { useActionState, useMemo, useState } from "react";
import { type Cuatrimestre, esPosterior, formatearCuatrimestre } from "@/lib/cuatrimestres";
import { parsearMetadata } from "@fiuba-reviews/parser-siu";
import { type ResultadoOferta, enviarOferta } from "./actions";
import { buscarCarrera } from "./utils";

/**
 * - `neutro`: todavía no se pegó nada.
 * - `ok`: la oferta se puede enviar.
 * - `aviso`: lo pegado no es una oferta del SIU o la carrera no está registrada.
 * - `error`: el servidor rechazó la oferta o falló al guardarla.
 * - `gracias`: ya tenemos esa oferta o una más nueva.
 */
type VarianteEstado = "neutro" | "ok" | "aviso" | "error" | "gracias";

const VARIANTES_ESTADO: Record<VarianteEstado, { icono: typeof Paste; clases: string }> = {
  neutro: { icono: Paste, clases: "text-stone-500" },
  ok: { icono: CircleCheckBig, clases: "text-emerald-700" },
  aviso: { icono: AlertTriangle, clases: "text-amber-800" },
  error: { icono: AlertCircle, clases: "text-red-700" },
  gracias: { icono: Handshake, clases: "text-stone-600" }
};

const INICIAL: ResultadoOferta = { estado: "inicial" };

/**
 * Valida lo pegado para avisar antes de enviar si la carrera, el cuatrimestre o la oferta tienen
 * algún problema. El servidor lo vuelve a validar en `enviarOferta`.
 */
function validarOferta(
  contenido: string,
  ofertas: OfertaCarrera[]
): { variante: VarianteEstado; mensaje: ReactNode; enviable: boolean } {
  if (contenido.trim() === "") {
    return { variante: "neutro", mensaje: "Pegá la página copiada del SIU.", enviable: false };
  }

  const [metadata] = parsearMetadata(contenido);

  if (!metadata) {
    return {
      variante: "aviso",
      mensaje:
        "No reconocemos esto como la oferta del SIU. Copiá la página entera, incluido el encabezado con tu carrera.",
      enviable: false
    };
  }

  const cuatrimestre = <Cuatri cuatrimestre={metadata.cuatrimestre} />;
  const carrera = buscarCarrera(ofertas, metadata.carrera);

  if (!carrera) {
    return {
      variante: "aviso",
      mensaje: <>No tenemos registrada la carrera «{metadata.carrera}».</>,
      enviable: false
    };
  }

  /* `recibida` solo está si es más nueva que `procesada`. */
  const ultima = carrera.recibida ?? carrera.procesada;
  const nombre = <span className="font-medium">{carrera.nombre}</span>;

  if (ultima && !esPosterior(metadata.cuatrimestre, ultima)) {
    return {
      variante: "gracias",
      mensaje: (
        <>
          Es la oferta de {nombre} del {cuatrimestre}, pero ya tenemos la del{" "}
          <Cuatri cuatrimestre={ultima} />. ¡Gracias igual!
        </>
      ),
      enviable: false
    };
  }

  return {
    variante: "ok",
    mensaje: (
      <>
        Oferta de {nombre} del {cuatrimestre}
        {ultima ? (
          <>
            , más nueva que la que tenemos (<Cuatri cuatrimestre={ultima} />)
          </>
        ) : null}
        .
      </>
    ),
    enviable: true
  };
}

function Cuatri({ cuatrimestre }: { cuatrimestre: Cuatrimestre }) {
  return (
    <span className="font-medium">
      <span className="font-mono tabular-nums">{formatearCuatrimestre(cuatrimestre)}</span>
    </span>
  );
}

function Estado({ variante, children }: { variante: VarianteEstado; children: ReactNode }) {
  const { icono: Icono, clases } = VARIANTES_ESTADO[variante];

  return (
    <p className={`flex items-start gap-2 text-sm ${clases}`}>
      <Icono size={18} aria-hidden="true" className="mt-px shrink-0" />
      <span className="min-w-0">{children}</span>
    </p>
  );
}

export function FormularioOferta({ ofertas }: { ofertas: OfertaCarrera[] }) {
  const [contenido, setContenido] = useState("");
  /* El resultado del servidor solo se muestra hasta que se vuelve a editar lo pegado. */
  const [editado, setEditado] = useState(false);
  const [resultado, formAction, enviando] = useActionState(
    async (anterior: ResultadoOferta, formData: FormData) => {
      setEditado(false);
      const resultado = await enviarOferta(anterior, formData);
      if (resultado.estado === "enviada") setContenido("");
      return resultado;
    },
    INICIAL
  );

  const validacion = useMemo(() => validarOferta(contenido, ofertas), [contenido, ofertas]);
  const respuesta = editado ? INICIAL : resultado;

  return (
    <form action={formAction} className="space-y-3">
      {respuesta.estado === "enviada" && (
        <output className="flex gap-2 rounded-md border-[1.5px] border-emerald-200 bg-emerald-50 px-3 py-2.5 text-sm text-emerald-800">
          <PartyPopper size={18} aria-hidden="true" className="mt-px shrink-0" />
          <p>
            <span className="font-medium">¡Gracias!</span> Recibimos la oferta de{" "}
            {respuesta.carrera} del <Cuatri cuatrimestre={respuesta.cuatrimestre} />. Las ofertas no
            se procesan al instante, así que puede tardar unos días en verse en las materias.
          </p>
        </output>
      )}

      <label className="block">
        <span className="sr-only">Página de la oferta de comisiones copiada del SIU</span>
        <textarea
          name="contenido"
          required
          value={contenido}
          onChange={(event) => {
            setContenido(event.target.value);
            setEditado(true);
          }}
          rows={6}
          spellCheck={false}
          autoComplete="off"
          placeholder="Pegá acá la página del SIU…"
          className="block w-full resize-y campo py-2 font-mono placeholder:font-sans md:text-xs/5 md:placeholder:text-sm"
        />
      </label>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div aria-live="polite" className="min-w-0 flex-1">
          {respuesta.estado === "error" ? (
            <Estado variante="error">{respuesta.mensaje}</Estado>
          ) : (
            <Estado variante={validacion.variante}>{validacion.mensaje}</Estado>
          )}
        </div>
        <button
          type="submit"
          disabled={!validacion.enviable || enviando}
          className="boton-primario self-end disabled:pointer-events-none disabled:bg-stone-100 disabled:text-stone-400 disabled:[--relieve:var(--color-stone-300)] sm:self-start"
        >
          {enviando ? (
            <LoaderCircle size={16} aria-hidden="true" className="shrink-0 animate-spin" />
          ) : (
            <Send size={16} aria-hidden="true" className="shrink-0" />
          )}
          {enviando ? "Enviando…" : "Enviar oferta"}
        </button>
      </div>
    </form>
  );
}
