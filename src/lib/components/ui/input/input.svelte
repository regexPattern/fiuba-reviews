<script lang="ts">
  import { cn } from "$lib/utils";
  import type { HTMLInputAttributes } from "svelte/elements";

  import type { InputEvents } from ".";

  type $$Props = HTMLInputAttributes;
  type $$Events = InputEvents;

  let className: $$Props["class"] = undefined;
  export let value: $$Props["value"] = undefined;
  export { className as class };
</script>

<input
  class={cn(
    "flex h-10 w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm ring-offset-background file:border-0  file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
    className,
  )}
  bind:value
  on:blur
  on:change
  on:click
  on:focus
  on:keydown
  on:keypress
  on:keyup
  on:mouseover
  on:mouseenter
  on:mouseleave
  on:paste
  on:input
  {...$$restProps} />
