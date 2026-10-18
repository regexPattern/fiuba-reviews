import { OG_SIZE, ogInicio } from "@/lib/og";

export const size = OG_SIZE;
export const contentType = "image/png";
export const alt = "FIUBA Reviews";

export default function Image() {
  return ogInicio();
}
