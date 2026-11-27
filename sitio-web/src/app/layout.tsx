import type { Metadata, Viewport } from "next";
import "@fontsource-variable/geist";
import "@fontsource-variable/geist-mono";
import { Header } from "@/components/layout/header";
import { env } from "@/env";
import "./globals.css";

const DESCRIPCION =
  "Encontrá calificaciones y comentarios de los docentes de la facultad, subidos por otros estudiantes. Basado en el legendario Dolly FIUBA.";

export const metadata: Metadata = {
  metadataBase: new URL(env.NEXT_PUBLIC_BASE_URL),
  title: { default: "FIUBA Reviews", template: "%s | FIUBA Reviews" },
  description: DESCRIPCION,
  authors: [{ name: "Carlos Eduardo Castillo Pereira" }],
  robots: { index: true, follow: true, "max-snippet": -1, "max-image-preview": "large" },
  openGraph: {
    type: "website",
    siteName: "FIUBA Reviews",
    locale: "es_AR",
    url: "/",
    title: "FIUBA Reviews",
    description: DESCRIPCION
  },
  twitter: {
    card: "summary",
    title: "FIUBA Reviews",
    description: DESCRIPCION
  }
};

// `viewportFit: "cover"` hace que la página ocupe toda la pantalla, incluso debajo del notch y
// de la barra flotante de Safari (iOS 26+). Por eso cada elemento pegado a un borde (header,
// drawers, footers) suma `env(safe-area-inset-*)` a su padding para no quedar tapado.
export const viewport: Viewport = { viewportFit: "cover" };

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="es-AR">
      <body className="relative bg-white font-sans antialiased">
        <Header />
        <div className="pt-(--header-height) pb-[env(safe-area-inset-bottom)] md:pb-0">
          {children}
        </div>
      </body>
    </html>
  );
}
