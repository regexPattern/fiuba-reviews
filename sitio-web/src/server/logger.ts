import "server-only";
import pino from "pino";
import { env } from "@/env";

export const logger = pino({
  level: env.LOG_LEVEL,
  transport: process.env.NODE_ENV === "development" ? { target: "pino-pretty" } : undefined
});
