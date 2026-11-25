// Extiende assets/fiuba.jpg para que entre toda la fachada, con un modelo de edición de imágenes
// vía OpenRouter. Uso: OPENROUTER_API_KEY=... node scripts/extend-fiuba-photo.mjs [--dry-run]
// Después: node scripts/extend-fiuba-photo.mjs --compose <generada.png> pega encima la foto real y
// escribe assets/fiuba-extendida.jpg.
import { existsSync } from "node:fs";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";
import sharp from "sharp";

const photo = fileURLToPath(new URL("../assets/fiuba.jpg", import.meta.url));
const workDir = process.env.FIUBA_WORK_DIR ?? "/tmp/fiuba-outpaint";
const model = process.env.FIUBA_IMAGE_MODEL ?? "google/gemini-3-pro-image";

// Lienzo 4:3 con la foto abajo a la izquierda: la fachada sigue hacia la derecha y hacia arriba.
export const canvas = { width: 1792, height: 1344, left: 64, top: 576 };
const blank = "#ff00ff";

// Sobre la foto (1024x768): el árbol de la esquina superior izquierda y las escaleras en sombra.
const erased = `<svg xmlns="http://www.w3.org/2000/svg" width="1024" height="768">
  <polygon points="0,0 215,0 215,140 105,250 40,298 0,316" fill="${blank}"/>
  <rect x="0" y="612" width="1024" height="156" fill="${blank}"/>
</svg>`;

export const prompt = `The first image is a photo of the Facultad de Ingeniería (University of Buenos Aires, Paseo Colón building) placed on a larger canvas. The flat magenta (#ff00ff) areas are blank. Outpaint the photo: replace every magenta area with a natural continuation of the same photograph, same camera position, lens, perspective and sunlight.

Requirements:
- Keep the canvas layout: the existing photo stays exactly where it is and at the same size, occupying only the lower-left part of the output. Do not zoom in, crop or re-frame; the output must show the whole canvas area.
- Do not change, move, restyle or re-render the non-magenta part of the photo: every column, the frieze with the inscription "FACULTAD DE INGENIERIA", the lattice windows, doors and shadows on the facade must stay pixel-identical and at the same position and scale.
- Complete the building toward the right and top: the colonnade continues with the same giant fluted columns, then ends in the solid stone corner block with a tall lattice window, exactly like the left end of the facade and like the reference photos. The stepped cornice and roofline continue along the same perspective line until the top-right corner of the building is fully visible, with clear blue sky above and to the right of it. The building's top-right roof corner must be well inside the frame, with a wide margin of empty sky between it and the top and right edges; nothing of the building may touch the top or right edge.
- The whole facade must be inside the frame, from its left corner to its right corner and from the steps to the roof, with sky around it.
- The magenta area in the top-left (where a tree was) becomes plain blue sky with a few soft clouds. No trees or branches anywhere.
- The magenta band at the bottom becomes the wide stone entrance staircase in even daylight, continuing across the full width under the extended building. No hard dark shadow on the steps, no cars, no people, no scaffolding.
- Photorealistic, no text other than the existing inscription, no borders, no magenta left anywhere.

The other images are reference photos of the same building from other angles: use them only to know how the ends of the facade, the cornice and the corner blocks look.`;

async function buildCanvas() {
  const cleaned = await sharp(photo)
    .autoOrient()
    .composite([{ input: Buffer.from(erased) }])
    .png()
    .toBuffer();
  return sharp({
    create: { width: canvas.width, height: canvas.height, channels: 3, background: blank }
  })
    .composite([{ input: cleaned, left: canvas.left, top: canvas.top }])
    .png()
    .toBuffer();
}

// La parte de la foto que se conserva tal cual: la fachada, sin el cielo, el árbol, las escaleras
// en sombra, los autos ni la columna cortada por el borde derecho.
const kept = `<svg xmlns="http://www.w3.org/2000/svg" width="1024" height="768">
  <polygon points="880,-12 880,598 365,598 365,528 12,528 12,335 55,310 118,256 230,191 330,132 430,73 520,21 556,-12" fill="#fff"/>
</svg>`;

// El modelo respeta el encuadre del lienzo pero vuelve a dibujar todo: se pegan los píxeles reales.
async function compose(generated) {
  const base = await sharp(generated)
    .resize(canvas.width, canvas.height, { fit: "fill" })
    .png()
    .toBuffer();
  const mask = await sharp(Buffer.from(kept)).blur(5).png().toBuffer();
  const original = await sharp(photo)
    .autoOrient()
    .ensureAlpha()
    .composite([{ input: mask, blend: "dest-in" }])
    .png()
    .toBuffer();
  await sharp(base)
    .composite([{ input: original, left: canvas.left, top: canvas.top }])
    .jpeg({ quality: 92 })
    .toFile(fileURLToPath(new URL("../assets/fiuba-extendida.jpg", import.meta.url)));
}

const dataUrl = (bytes, mime = "image/png") => `data:${mime};base64,${bytes.toString("base64")}`;

async function main() {
  const composeAt = process.argv.indexOf("--compose");
  if (composeAt !== -1) return compose(process.argv[composeAt + 1]);

  await mkdir(workDir, { recursive: true });
  // FIUBA_INPUT y FIUBA_PROMPT_FILE permiten una segunda pasada de corrección sobre un resultado.
  const input = process.env.FIUBA_INPUT
    ? await readFile(process.env.FIUBA_INPUT)
    : await buildCanvas();
  const text = process.env.FIUBA_PROMPT_FILE
    ? await readFile(process.env.FIUBA_PROMPT_FILE, "utf8")
    : prompt;
  if (!process.env.FIUBA_INPUT) await writeFile(`${workDir}/canvas.png`, input);
  if (process.argv.includes("--dry-run")) return;

  const key = process.env.OPENROUTER_API_KEY;
  if (!key) throw new Error("Falta OPENROUTER_API_KEY en el entorno");

  const images = [dataUrl(input)];
  for (const name of ["ref1.png", "ref2.png", "ref3.png"]) {
    const path = `${workDir}/${name}`;
    if (!existsSync(path)) continue;
    const small = await sharp(await readFile(path))
      .resize(1024, 1024, { fit: "inside" })
      .jpeg({ quality: 88 })
      .toBuffer();
    images.push(dataUrl(small, "image/jpeg"));
  }

  const response = await fetch("https://openrouter.ai/api/v1/chat/completions", {
    method: "POST",
    headers: { Authorization: `Bearer ${key}`, "Content-Type": "application/json" },
    body: JSON.stringify({
      model,
      modalities: ["image", "text"],
      image_config: { aspect_ratio: "4:3", image_size: "2K" },
      messages: [
        {
          role: "user",
          content: [
            { type: "text", text },
            ...images.map((url) => ({ type: "image_url", image_url: { url } }))
          ]
        }
      ]
    })
  });
  const body = await response.json();
  if (!response.ok || body.error) {
    throw new Error(`OpenRouter ${response.status}: ${JSON.stringify(body.error ?? body)}`);
  }

  const message = body.choices?.[0]?.message;
  const outputs = message?.images ?? [];
  if (outputs.length === 0) {
    throw new Error(`El modelo no devolvió imágenes: ${String(message?.content).slice(0, 500)}`);
  }
  for (const [i, { image_url }] of outputs.entries()) {
    const bytes = Buffer.from(image_url.url.slice(image_url.url.indexOf(",") + 1), "base64");
    const path = `${workDir}/extended-${process.env.FIUBA_VARIANT ?? "a"}-${i}.png`;
    await sharp(bytes).png().toFile(path);
    console.log(path);
  }
}

if (process.argv[1] === fileURLToPath(import.meta.url)) await main();
