import "server-only";
import { env } from "@/env";

type TurnstileResponse = { success: boolean; "error-codes"?: string[] };

export async function validarTokenTurnstile(token: string) {
  const response = await fetch("https://challenges.cloudflare.com/turnstile/v0/siteverify", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ response: token, secret: env.TURNSTILE_SECRET_KEY })
  });

  const data: TurnstileResponse = await response.json();

  return { success: data.success, error: data["error-codes"]?.[0] ?? null };
}
