import Image from "next/image";

export default function Home() {
  return (
    <main className="mx-auto max-w-6xl p-6">
      <h1 className="text-3xl font-semibold tracking-tight">FIUBA Reviews</h1>
      <Image
        src="/fiuba-scene.svg"
        alt="Ilustración de la fachada de la Facultad de Ingeniería de la UBA en puntos de luz"
        width={1720}
        height={1200}
        className="mt-6 h-auto w-full max-w-4xl bg-slate-950"
        preload
        unoptimized
      />
    </main>
  );
}
