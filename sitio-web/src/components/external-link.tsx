import { cn } from "@/lib/utils";

export function ExternalLink({ className, ...props }: React.ComponentProps<"a">) {
  return (
    <a
      target="_blank"
      rel="external noopener noreferrer"
      className={cn("underline underline-offset-2", className)}
      {...props}
    />
  );
}
