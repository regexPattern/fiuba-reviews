"use client";

import type { OfertaCarrera } from "@/server/db/queries/ofertas";
import type { Icon } from "@phosphor-icons/react";
import type { ReactNode } from "react";
import { Select } from "@base-ui/react/select";
import { AndroidLogo, AppleLogo, Desktop } from "@phosphor-icons/react/ssr";
import { ArrowUpRight, Check, ChevronDown, ChevronRight } from "@withicons/react/line";
import Image from "next/image";
import { useState } from "react";
import siuLogo from "@/assets/siu-logo.png";
import { FormularioOferta } from "@/components/planes/formulario-oferta";
import { type Plataforma, usePlataforma } from "@/lib/plataforma";

const URL_OFERTA_SIU = "https://guaraniautogestion.fi.uba.ar/g3w/oferta_comisiones";

/** Iconos de Phosphor: withicons no tiene logos de marcas. */
const PLATAFORMAS: { valor: Plataforma; etiqueta: string; icono: Icon }[] = [
  { valor: "windows", etiqueta: "Windows / Linux", icono: Desktop },
  { valor: "mac", etiqueta: "macOS", icono: AppleLogo },
  { valor: "android", etiqueta: "Android", icono: AndroidLogo },
  { valor: "ios", etiqueta: "iPhone / iPad", icono: AppleLogo }
];

/** Cómo seleccionar y copiar la página en cada plataforma (pasos 2 y 3). */
const COPIAR: Record<Plataforma, { seleccionar: ReactNode; copiar: ReactNode }> = {
  windows: copiarConTeclado("Ctrl"),
  mac: copiarConTeclado("⌘"),
  android: {
    seleccionar: (
      <Gesto>
        Mantené presionada cualquier palabra y tocá <Opcion>Seleccionar todo</Opcion> en el menú que
        aparece.
      </Gesto>
    ),
    copiar: (
      <Gesto>
        Tocá <Opcion>Copiar</Opcion> en el mismo menú.
      </Gesto>
    )
  },
  ios: {
    seleccionar: (
      <Gesto>
        Mantené presionada la primera palabra de la página y arrastrá el final de la selección hasta
        abajo de todo, para que entre desde tu carrera hasta la última comisión.
      </Gesto>
    ),
    copiar: (
      <Gesto>
        Tocá <Opcion>Copiar</Opcion> en el menú que aparece sobre la selección.
      </Gesto>
    )
  }
};

/** En la compu es lo mismo con Ctrl o con ⌘. */
function copiarConTeclado(modificador: string) {
  return {
    seleccionar: (
      <Gesto>
        Apretá <Atajo teclas={[modificador, "A"]} /> para seleccionar toda la página.
      </Gesto>
    ),
    copiar: (
      <Gesto>
        Apretá <Atajo teclas={[modificador, "C"]} /> para copiar todo lo seleccionado.
      </Gesto>
    )
  };
}

export function PasosOferta({ ofertas }: { ofertas: OfertaCarrera[] }) {
  const detectada = usePlataforma();
  const [elegida, setElegida] = useState<Plataforma | null>(null);
  const plataforma = elegida ?? detectada;

  return (
    <div className="space-y-5">
      <SelectorPlataforma valor={plataforma} onChange={setElegida} />
      <ol className="space-y-5">
        <Paso numero={1} titulo="Abrí la oferta de comisiones en el SIU">
          <a
            href={URL_OFERTA_SIU}
            target="_blank"
            rel="noopener noreferrer"
            className="boton-secundario w-fit"
          >
            <Image
              src={siuLogo}
              alt="SIU Guaraní"
              width={16}
              height={16}
              draggable={false}
              className="size-4 shrink-0 rounded-sm"
            />
            <ChevronRight size={14} aria-hidden="true" className="shrink-0 text-stone-400" />
            Reportes
            <ChevronRight size={14} aria-hidden="true" className="shrink-0 text-stone-400" />
            Oferta de comisiones
            <ArrowUpRight size={16} aria-hidden="true" className="shrink-0 text-stone-500" />
          </a>
        </Paso>
        <Paso numero={2} titulo="Seleccioná toda la página">
          {COPIAR[plataforma].seleccionar}
        </Paso>
        <Paso numero={3} titulo="Copiala">
          {COPIAR[plataforma].copiar}
        </Paso>
        <Paso numero={4} titulo="Pegala acá abajo y enviala">
          <FormularioOferta ofertas={ofertas} />
        </Paso>
      </ol>
    </div>
  );
}

/** Dropdown con la plataforma: viene detectada, así que es secundario y alcanza con una línea. */
function SelectorPlataforma({
  valor,
  onChange
}: {
  valor: Plataforma;
  onChange: (plataforma: Plataforma) => void;
}) {
  const { icono: IconoElegido } = PLATAFORMAS.find((p) => p.valor === valor) ?? PLATAFORMAS[0];

  return (
    <Select.Root
      items={PLATAFORMAS.map(({ valor, etiqueta }) => ({ value: valor, label: etiqueta }))}
      value={valor}
      onValueChange={(nueva) => nueva && onChange(nueva)}
    >
      <div className="flex flex-wrap items-center gap-x-2 gap-y-1.5">
        <Select.Label className="text-sm text-stone-500">Instrucciones para</Select.Label>
        <Select.Trigger className="group boton-secundario">
          <IconoElegido size={16} aria-hidden="true" className="shrink-0 text-stone-500" />
          <Select.Value />
          <ChevronDown
            size={14}
            aria-hidden="true"
            className="text-stone-400 transition-transform duration-150 group-data-popup-open:rotate-180 motion-reduce:transition-none"
          />
        </Select.Trigger>
      </div>
      <Select.Portal>
        <Select.Positioner
          side="bottom"
          align="start"
          sideOffset={6}
          alignItemWithTrigger={false}
          className="z-50"
        >
          <Select.Popup className="min-w-(--anchor-width) origin-(--transform-origin) relieve bg-white p-1 text-sm transition-[transform,opacity] duration-150 outline-none data-ending-style:scale-95 data-ending-style:opacity-0 data-starting-style:scale-95 data-starting-style:opacity-0 motion-reduce:transition-none">
            <Select.List>
              {PLATAFORMAS.map(({ valor: opcion, etiqueta, icono: Icono }) => (
                <Select.Item
                  key={opcion}
                  value={opcion}
                  className="flex cursor-pointer items-center gap-2 rounded-md py-1.5 pr-2 pl-2 text-stone-700 outline-none select-none data-highlighted:bg-stone-100 data-highlighted:text-stone-900 data-selected:font-medium data-selected:text-stone-900"
                >
                  <Icono size={16} aria-hidden="true" className="shrink-0 text-stone-500" />
                  <Select.ItemText className="flex-1">{etiqueta}</Select.ItemText>
                  <Select.ItemIndicator className="ml-3 text-fiuba">
                    <Check size={14} aria-hidden="true" />
                  </Select.ItemIndicator>
                </Select.Item>
              ))}
            </Select.List>
          </Select.Popup>
        </Select.Positioner>
      </Select.Portal>
    </Select.Root>
  );
}

/** Paso numerado: el número en un círculo y una línea que lo une con el siguiente. */
function Paso({
  numero,
  titulo,
  children
}: {
  numero: number;
  titulo: string;
  children: ReactNode;
}) {
  return (
    <li className="relative flex gap-3 pb-1 not-last:after:absolute not-last:after:top-8 not-last:after:-bottom-4 not-last:after:left-3 not-last:after:w-px not-last:after:-translate-x-1/2 not-last:after:bg-stone-200">
      <span
        aria-hidden="true"
        className="flex size-6 shrink-0 items-center justify-center rounded-full border-[1.5px] border-stone-300 bg-white font-mono text-xs text-stone-600"
      >
        {numero}
      </span>
      <div className="min-w-0 flex-1 space-y-2.5">
        <h3 className="text-base/6 font-medium text-stone-900">{titulo}</h3>
        {children}
      </div>
    </li>
  );
}

/** Combinación de teclas, cada una en relieve, dentro de una oración. */
function Atajo({ teclas }: { teclas: string[] }) {
  return (
    <span className="mx-0.5 inline-flex items-center gap-1 align-middle">
      {teclas.map((tecla) => (
        <Tecla key={tecla}>{tecla}</Tecla>
      ))}
    </span>
  );
}

/** Tecla en relieve, como los botones. */
function Tecla({ children }: { children: ReactNode }) {
  return (
    <kbd className="inline-flex h-6 min-w-6 items-center justify-center relieve bg-white px-1.5 font-mono text-xs text-stone-700 [--relieve-alto:2px]">
      {children}
    </kbd>
  );
}

/** Instrucción de un paso. Interlineado de 28px para que las teclas en relieve no separen las líneas. */
function Gesto({ children }: { children: ReactNode }) {
  return <p className="text-sm/7 text-stone-600">{children}</p>;
}

/** Nombre de una opción de menú del sistema. */
function Opcion({ children }: { children: ReactNode }) {
  return <span className="font-medium text-stone-900">{children}</span>;
}
