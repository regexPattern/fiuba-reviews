import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  cacheComponents: true,
  partialPrefetching: true,
  typedRoutes: true,
  experimental: {
    // Work around instant-validation client manifest errors in the dev worker.
    devValidationWorker: false
  },
  async redirects() {
    return [
      { source: "/materias", destination: "/", permanent: true },
      { source: "/materias/:codigo", destination: "/materia/:codigo", permanent: true }
    ];
  }
};

export default nextConfig;
