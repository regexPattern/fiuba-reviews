import type { NextConfig } from "next";
import "./src/env";

const nextConfig: NextConfig = {
  cacheComponents: true,
  partialPrefetching: true,
  typedRoutes: true,
  allowedDevOrigins: ["*.local", "192.168.*.*"]
};

export default nextConfig;
