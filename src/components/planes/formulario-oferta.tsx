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
import { type ResultadoEnvio, enviarOferta } from "@/app/planes/actions";
import { Numero } from "@/components/ui/numero";
import {
  type Cuatrimestre,
  extraerMetadataOferta,
  formatearCuatrimestre,
  normalizarNombre,
  ordenCuatrimestre
} from "@/lib/ofertas";

type Tono = "neutro" | "ok" | "aviso" | "error" | "gracias";

const TONOS: Record<Tono, { icono: typeof Paste; clases: string }> = {
  neutro: { icono: Paste, clases: "text-stone-500" },
  ok: { icono: CircleCheckBig, clases: "text-emerald-700" },
  aviso: { icono: AlertTriangle, clases: "text-amber-800" },
  error: { icono: AlertCircle, clases: "text-red-700" },
  gracias: { icono: Handshake, clases: "text-stone-600" }
};

const INICIAL: ResultadoEnvio = { estado: "inicial" };

/** La más nueva entre la procesada y la recibida: contra esa se compara lo que se pega. */
function ultimaConocida({ procesada, recibida }: OfertaCarrera) {
  return recibida ?? procesada;
}

/**
 * Qué se detectó en lo pegado y si se puede enviar. Lo mismo se vuelve a validar en el servidor;
 * acá es para avisar antes de enviar.
 */
function diagnosticar(
  contenido: string,
  ofertas: OfertaCarrera[]
): { tono: Tono; mensaje: ReactNode; enviable: boolean } {
  if (contenido.trim() === "") {
    return { tono: "neutro", mensaje: "Pegá la página copiada del SIU.", enviable: false };
  }

  const metadata = extraerMetadataOferta(contenido);

  if (!metadata) {
    return {
      tono: "aviso",
      mensaje:
        "No reconocemos esto como la oferta del SIU. Copiá la página entera, incluido el encabezado con tu carrera.",
      enviable: false
    };
  }

  const cuatrimestre = <Cuatri cuatrimestre={metadata.cuatrimestre} />;
  const carrera = ofertas.find(
    (o) => normalizarNombre(o.nombre) === normalizarNombre(metadata.carrera)
  );

  if (!carrera) {
    return {
      tono: "aviso",
      mensaje: <>No tenemos registrada la carrera «{metadata.carrera}».</>,
      enviable: false
    };
  }

  const ultima = ultimaConocida(carrera);
  const nombre = <span className="font-medium">{carrera.nombre}</span>;

  if (ultima && ordenCuatrimestre(metadata.cuatrimestre) <= ordenCuatrimestre(ultima)) {
    return {
      tono: "gracias",
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
    tono: "ok",
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
      <Numero>{formatearCuatrimestre(cuatrimestre)}</Numero>
    </span>
  );
}

function Estado({ tono, children }: { tono: Tono; children: ReactNode }) {
  const { icono: Icono, clases } = TONOS[tono];

  return (
    <p className={`flex items-start gap-2 text-sm ${clases}`}>
      <Icono size={18} aria-hidden="true" className="mt-px shrink-0" />
      <span className="min-w-0">{children}</span>
    </p>
  );
}

export function FormularioOferta({ ofertas }: { ofertas: OfertaCarrera[] }) {
  const [contenido, setContenido] = useState("");
  // El resultado del servidor solo se muestra hasta que se vuelve a editar lo pegado.
  const [editado, setEditado] = useState(false);
  const [resultado, accion, enviando] = useActionState(
    async (anterior: ResultadoEnvio, formData: FormData) => {
      setEditado(false);
      const resultado = await enviarOferta(anterior, formData);
      if (resultado.estado === "enviada") setContenido("");
      return resultado;
    },
    INICIAL
  );

  const diagnostico = useMemo(() => diagnosticar(contenido, ofertas), [contenido, ofertas]);
  const respuesta = editado ? INICIAL : resultado;

  return (
    <form action={accion} className="space-y-3">
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
        {/* 16px en móvil: con menos, iOS hace zoom al enfocar. */}
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
          className="block w-full resize-y rounded-md border-[1.5px] border-stone-300 bg-white px-2.5 py-2 font-mono text-base outline-none placeholder:font-sans placeholder:text-stone-400 focus-visible:border-fiuba focus-visible:ring-3 focus-visible:ring-fiuba/20 md:text-xs/5 md:placeholder:text-sm"
        />
      </label>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div aria-live="polite" className="min-w-0 flex-1">
          {respuesta.estado === "error" ? (
            <Estado tono="error">{respuesta.mensaje}</Estado>
          ) : (
            <Estado tono={diagnostico.tono}>{diagnostico.mensaje}</Estado>
          )}
        </div>
        <button
          type="submit"
          disabled={!diagnostico.enviable || enviando}
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
