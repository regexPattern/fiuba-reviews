import { piece, spread } from "ascii.rest/kit";
import { svg } from "ascii.rest/svg";
import { writeFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";
import sharp from "sharp";

const photo = fileURLToPath(new URL("../assets/fiuba-extendida.jpg", import.meta.url));
const output = new URL("../public/fiuba-scene.svg", import.meta.url);
const cols = 172;
const rows = 120;
const ground = "#090c18";

// La foto extendida es 4:3; se recorta cielo de arriba para quedar en la proporción de la grilla.
const { width, height } = await sharp(photo).metadata();
const cropHeight = Math.round((width * rows) / cols);
const { data } = await sharp(photo)
  .extract({ left: 0, top: height - cropHeight, width, height: cropHeight })
  .resize(cols, rows, { fit: "fill" })
  .toColourspace("srgb")
  .removeAlpha()
  .raw()
  .toBuffer({ resolveWithObject: true });

const dots = " ·•●";
const bayer = [0, 8, 2, 10, 12, 4, 14, 6, 3, 11, 1, 9, 15, 7, 13, 5];
const shades = 16;
// Escena de día: cielo celeste y piedra clara, con las sombras en un gris cálido.
const palette = [
  ...spread(["#4f9fe3", "#74b9ee", "#9fd2f6", "#cdeafc"], shades),
  ...spread(["#2e2a3a", "#5d5560", "#948a88", "#cdbfae", "#eedfc9", "#fff6e6"], shades)
];
const clamp = (v) => Math.min(1, Math.max(0, v));

const scene = piece(
  { name: "fiuba", category: "scenes", cols, rows, cell: 1, fps: 0, ground, palette },
  (_, surface) => {
    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {
        const i = (y * cols + x) * 3;
        const r = data[i] / 255;
        const g = data[i + 1] / 255;
        const b = data[i + 2] / 255;
        const luma = 0.2126 * r + 0.7152 * g + 0.0722 * b;
        // Cielo: azul, o nube clara sin el tinte cálido de la piedra.
        const sky = (b - r > 0.08 && b > g) || (b - r > 0.0 && luma > 0.6);

        // El cielo no sale de la foto: es un degradé parejo, más claro hacia el horizonte, para que
        // las nubes y el empalme con el cielo generado no compitan con la fachada.
        // Sobre fondo oscuro un punto chico se lee casi negro: la piedra en sombra se aclara de más para
        // que se vea como en la foto.
        const tint = sky ? (y / rows) ** 1.3 : clamp((luma - 0.1) / 0.78) ** 0.45;
        const light = sky ? 0.86 + 0.14 * tint : 0.5 + 0.5 * tint;

        // El umbral ordenado rompe las zonas planas sin desdibujar la grilla de puntos.
        const threshold = (bayer[(y % 4) * 4 + (x % 4)] + 0.5) / 16 - 0.5;
        const size = Math.max(1, Math.min(3, Math.round(light * 3 + threshold)));

        const tone = Math.round(tint * (shades - 1));
        surface.set(x, y, dots[size], (sky ? 0 : shades) + tone);
      }
    }
  }
);

await writeFile(output, svg(scene, { label: "Fachada de la FIUBA, ilustración en puntos" }));
