import { TURNSTILE_SECRET_KEY } from "$env/static/private";

interface TokenValidateResponse {
  "error-codes": string[];
  success: boolean;
  action: string;
  cdata: string;
}

export async function validateToken(token: string) {
  const response = await fetch("https://challenges.cloudflare.com/turnstile/v0/siteverify", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ response: token, secret: TURNSTILE_SECRET_KEY })
  });

  const data: TokenValidateResponse = await response.json();

  return {
    success: data.success,
    error: data["error-codes"]?.length ? data["error-codes"][0] : null
  };
}
