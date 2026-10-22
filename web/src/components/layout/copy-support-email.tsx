"use client";

import { Mail, MailCheck } from "@withicons/react/chrome";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function CopySupportEmail({ email }: { email: string }) {
  const [status, setStatus] = useState<"idle" | "copied" | "error">("idle");

  useEffect(() => {
    if (status === "idle") return;
    const timeout = window.setTimeout(() => setStatus("idle"), 2000);
    return () => window.clearTimeout(timeout);
  }, [status]);

  async function copyEmail() {
    try {
      await navigator.clipboard.writeText(email);
      setStatus("copied");
    } catch {
      setStatus("error");
    }
  }

  const Icon = status === "copied" ? MailCheck : Mail;

  return (
    <>
      <Button
        variant="ghost"
        className="flex h-auto w-full justify-start p-2 font-normal"
        onClick={copyEmail}
        aria-label="Copiar mail de contacto"
      >
        <span
          key={status}
          className="-mx-0.5 inline-flex shrink-0 animate-in duration-200 zoom-in motion-reduce:animate-none"
        >
          <Icon size={20} className="size-5" aria-hidden="true" />
        </span>
        <span className="relative text-left whitespace-nowrap">
          <span className={cn("text-primary", status !== "idle" && "invisible")}>
            Mandame un mail
          </span>
          {status !== "idle" && (
            <span
              key={status}
              className="absolute inset-0 animate-in duration-200 fade-in motion-reduce:animate-none"
            >
              {status === "copied" ? "¡Copiado!" : "No se pudo copiar"}
            </span>
          )}
        </span>
      </Button>
      <output className="sr-only">{status === "copied" && "Mail copiado"}</output>
    </>
  );
}
