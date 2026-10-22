import type { Metadata, Viewport } from "next";
import { Header } from "@/components/layout/header";
import { env } from "@/env";
import { ibmPlexMono, inter } from "@/lib/fonts";
import "./globals.css";

export const metadata: Metadata = {
  metadataBase: new URL(env.NEXT_PUBLIC_BASE_URL),
  title: { default: "FIUBA Reviews", template: "%s" },
  authors: [{ name: "Carlos Eduardo Castillo Pereira" }],
  icons: { icon: [{ url: "/favicon.svg", type: "image/svg+xml" }, { url: "/favicon.ico" }] },
  openGraph: { type: "website", siteName: "FIUBA Reviews", locale: "es_AR" }
};

export const viewport: Viewport = { viewportFit: "cover" };

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="es" className={`${inter.variable} ${ibmPlexMono.variable}`}>
      <body className="flex min-h-dvh flex-col font-sans md:h-dvh antialiased">
        <Header />
        <div className="flex-1 pt-(--header-height) pb-[env(safe-area-inset-bottom)] md:min-h-0 md:overflow-y-auto md:pb-0">{children}</div>
      </body>
    </html>
  );
}
