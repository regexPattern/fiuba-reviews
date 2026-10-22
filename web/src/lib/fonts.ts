import { IBM_Plex_Mono, Inter } from "next/font/google";

export const ibmPlexMono = IBM_Plex_Mono({
  subsets: ["latin"],
  weight: ["400", "500", "600"],
  display: "swap",
  variable: "--font-ibm-plex-mono"
});

export const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-inter"
});
