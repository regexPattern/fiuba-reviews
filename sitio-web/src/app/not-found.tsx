import type { Metadata } from "next";
import { ErrorView } from "@/components/error-view";

export const metadata: Metadata = {
  title: "No encontrado | FIUBA Reviews",
  robots: { index: false, follow: false }
};

export default function NotFound() {
  return <ErrorView codigo={404} mensaje="La página que buscás no existe." />;
}
