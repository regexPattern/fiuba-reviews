import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  cacheComponents: true,
  partialPrefetching: true,
  typedRoutes: true,
  async redirects() {
    return [
      { source: "/materias", destination: "/", permanent: true },
      { source: "/materias/:codigo", destination: "/materia/:codigo", permanent: true }
    ];
  }
};

export default nextConfig;
