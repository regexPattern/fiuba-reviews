import { ImageResponse } from "next/og";
import { readFile } from "node:fs/promises";
import { join } from "node:path";

export const OG_SIZE = { width: 1200, height: 630 } as const;

const FUENTES_DIR = join(process.cwd(), "src/assets/fonts");

async function fuentes() {
  const [regular, semibold] = await Promise.all([
    readFile(join(FUENTES_DIR, "source-serif-4-latin-400-normal.woff")),
    readFile(join(FUENTES_DIR, "source-serif-4-latin-600-normal.woff"))
  ]);

  return [
    { name: "Source Serif 4", data: regular, weight: 400 as const, style: "normal" as const },
    { name: "Source Serif 4", data: semibold, weight: 600 as const, style: "normal" as const }
  ];
}

function Marca({ size, gap }: { size: number; gap: number }) {
  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        fontSize: size,
        letterSpacing: "-0.025em",
        gap,
        fontWeight: 600
      }}
    >
      <span style={{ color: "#4facd4" }}>FIUBA</span>
      <span style={{ color: "#ececec" }}>Reviews</span>
    </div>
  );
}

export async function ogInicio() {
  return new ImageResponse(
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        width: "100%",
        height: "100%",
        backgroundColor: "#342F37",
        padding: 60,
        fontFamily: "Source Serif 4"
      }}
    >
      <Marca size={80} gap={12} />
    </div>,
    { ...OG_SIZE, fonts: await fuentes() }
  );
}

export async function ogMateria(materia: { codigo: string; nombre: string }) {
  return new ImageResponse(
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        width: "100%",
        height: "100%",
        backgroundColor: "#342F37",
        padding: 60,
        fontFamily: "Source Serif 4"
      }}
    >
      <div style={{ display: "flex", color: "#342F37", fontSize: 32 }}>placeholder</div>

      <div style={{ display: "flex", justifyContent: "center", alignItems: "center" }}>
        <div
          style={{
            display: "flex",
            color: "#ececec",
            fontSize: 64,
            letterSpacing: "-0.025em",
            maxHeight: 240,
            lineHeight: 1.25,
            overflow: "hidden"
          }}
        >
          {materia.codigo} • {materia.nombre}
        </div>
      </div>

      <Marca size={32} gap={8} />
    </div>,
    { ...OG_SIZE, fonts: await fuentes() }
  );
}
