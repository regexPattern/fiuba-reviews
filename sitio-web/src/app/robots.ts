import type { MetadataRoute } from "next";
import { env } from "@/env";

export default function robots(): MetadataRoute.Robots {
  return {
    // Las cátedras no se bloquean acá: Google tiene que poder leerlas para ver su `noindex`.
    rules: { userAgent: "*", allow: "/", disallow: "/api/" },
    sitemap: new URL("/sitemap.xml", env.NEXT_PUBLIC_BASE_URL).toString()
  };
}
