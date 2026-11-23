import type { NextConfig } from "next";
import "./src/env";

const nextConfig: NextConfig = {
  cacheComponents: true,
  typedRoutes: true
};

export default nextConfig;
