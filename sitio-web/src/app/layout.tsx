import { Google_Sans_Code, Inter, Source_Serif_4 } from "next/font/google";
import type { Metadata, Viewport } from "next";
import { BuscadorDialog } from "@/components/buscador/buscador-dialog";
import { BuscadorProvider } from "@/components/buscador/buscador-provider";
import { Header } from "@/components/layout/header";
import { ThemeProvider } from "@/components/theme-provider";
import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { env } from "@/env";
import { getMateriasVigentes } from "@/server/queries/materias-vigentes";
import "./globals.css";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });
const sourceSerif = Source_Serif_4({ subsets: ["latin"], variable: "--font-source-serif" });
const googleSansCode = Google_Sans_Code({
  subsets: ["latin"],
  variable: "--font-google-sans-code"
});

export const metadata: Metadata = {
  metadataBase: new URL(env.NEXT_PUBLIC_BASE_URL),
  title: { default: "FIUBA Reviews", template: "%s" },
  authors: [{ name: "Carlos Eduardo Castillo Pereira" }],
  icons: { icon: [{ url: "/favicon.svg", type: "image/svg+xml" }, { url: "/favicon.ico" }] },
  openGraph: { type: "website", siteName: "FIUBA Reviews", locale: "es_AR" }
};

export const viewport: Viewport = { viewportFit: "cover" };

export default async function RootLayout({ children }: { children: React.ReactNode }) {
  const materias = await getMateriasVigentes();

  return (
    <html
      lang="es"
      suppressHydrationWarning
      className={`${inter.variable} ${sourceSerif.variable} ${googleSansCode.variable}`}
    >
      <body className="font-sans antialiased">
        <ThemeProvider>
          <TooltipProvider>
            <BuscadorProvider materias={materias}>
              <Header />
              <BuscadorDialog />
              <div className="pt-[calc(56px+env(safe-area-inset-top))]">{children}</div>
            </BuscadorProvider>
          </TooltipProvider>
          <Toaster />
        </ThemeProvider>
      </body>
    </html>
  );
}
