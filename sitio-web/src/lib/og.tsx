import { ImageResponse } from "next/og";

export const OG_SIZE = { width: 1200, height: 630 } as const;

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
      <span style={{ color: "#171717" }}>FIUBA</span>
      <span style={{ color: "#171717" }}>Reviews</span>
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
        backgroundColor: "#ffffff",
        padding: 60,
        fontFamily: "sans-serif"
      }}
    >
      <Marca size={80} gap={12} />
    </div>,
    { ...OG_SIZE }
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
        backgroundColor: "#ffffff",
        padding: 60,
        fontFamily: "sans-serif"
      }}
    >
      <div style={{ display: "flex", color: "#ffffff", fontSize: 32 }}>placeholder</div>

      <div style={{ display: "flex", justifyContent: "center", alignItems: "center" }}>
        <div
          style={{
            display: "flex",
            color: "#171717",
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
    { ...OG_SIZE }
  );
}
