"use client";

import { Command as CommandPrimitive } from "cmdk";
import { SearchIcon } from "lucide-react";
import * as React from "react";
import { cn } from "@/lib/utils";

function Command({ className, ...props }: React.ComponentProps<typeof CommandPrimitive>) {
  return (
    <CommandPrimitive
      data-slot="command"
      className={cn(
        "flex w-full flex-col overflow-hidden rounded-md border bg-popover text-popover-foreground",
        className
      )}
      {...props}
    />
  );
}

function CommandInput({
  className,
  trailingAction,
  ...props
}: React.ComponentProps<typeof CommandPrimitive.Input> & { trailingAction?: React.ReactNode }) {
  return (
    <div className="flex items-center gap-2 border-b px-3">
      <SearchIcon aria-hidden="true" className="size-4 shrink-0 text-muted-foreground" />
      <div className="flex h-11 min-w-0 flex-1 items-center">
        {/* Keep 16px for iOS focus zoom, but render at 14px on mobile. */}
        <CommandPrimitive.Input
          data-slot="command-input"
          className={cn(
            // eslint-disable-next-line shadcn/no-arbitrary-values -- Exact 14/16 scale avoids iOS focus zoom; reciprocal width preserves layout.
            "h-full w-[calc(100%/0.875)] min-w-0 shrink-0 origin-left scale-[0.875] bg-transparent text-base outline-none placeholder:text-muted-foreground disabled:opacity-50 md:w-full md:scale-100 md:text-sm",
            className
          )}
          {...props}
        />
      </div>
      {trailingAction}
    </div>
  );
}

function CommandList({ className, ...props }: React.ComponentProps<typeof CommandPrimitive.List>) {
  return (
    <CommandPrimitive.List
      data-slot="command-list"
      className={cn("max-h-96 overflow-x-hidden overflow-y-auto p-1", className)}
      {...props}
    />
  );
}

function CommandEmpty(props: React.ComponentProps<typeof CommandPrimitive.Empty>) {
  return (
    <CommandPrimitive.Empty
      data-slot="command-empty"
      className="py-6 text-center text-sm"
      {...props}
    />
  );
}

function CommandItem({ className, ...props }: React.ComponentProps<typeof CommandPrimitive.Item>) {
  return (
    <CommandPrimitive.Item
      data-slot="command-item"
      className={cn(
        "flex cursor-pointer items-center gap-3 rounded-sm px-2 py-3 text-sm outline-none select-none data-[disabled=true]:pointer-events-none data-[disabled=true]:opacity-50 data-[selected=true]:bg-accent data-[selected=true]:text-accent-foreground",
        className
      )}
      {...props}
    />
  );
}

export { Command, CommandInput, CommandList, CommandEmpty, CommandItem };
