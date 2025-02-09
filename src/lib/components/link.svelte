<script lang="ts">
  import { cn } from "$lib/utils";

  export let href: string;

  let className = "";
  export { className as class };

  const isExternal = !href.startsWith("/");
</script>

{#if isExternal}
  <a
    {href}
    target="_blank"
    rel="noopener noreferrer"
    class={cn("underline", className)}
    {...$$restProps}><slot /></a
  >
{:else}
  <a {href} {...$$restProps} class={className}><slot /></a>
{/if}
