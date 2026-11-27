import type { NextConfig } from "next";
import "./src/env";

const nextConfig: NextConfig = {
  transpilePackages: ["@fiuba-reviews/db", "@fiuba-reviews/parser-siu"],
  cacheComponents: true,
  partialPrefetching: true,
  typedRoutes: true,
  allowedDevOrigins: ["*.local", "192.168.*.*"]
};

export default nextConfig;
