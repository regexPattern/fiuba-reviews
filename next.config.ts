import type { NextConfig } from "next";
import "./src/env";

const nextConfig: NextConfig = {
  cacheComponents: true,
  partialPrefetching: true,
  typedRoutes: true,
  async redirects() {
    // URLs de la versión anterior del sitio.
    return [
      {
        source: "/materia/:codigo",
        has: [{ type: "query", key: "catedra", value: "(?<catedra>.+)" }],
        destination: "/materias/:codigo/catedras/:catedra",
        permanent: true
      },
      { source: "/materia/:codigo", destination: "/materias/:codigo", permanent: true }
    ];
  }
};

export default nextConfig;
