<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLAnchorAttributes } from "svelte/elements";

  interface Props extends HTMLAnchorAttributes {
    href: string;
    children: Snippet;
  }

  let { href, children, ...others }: Props = $props();
</script>

<a {href} target="_blank" rel="external noopener noreferrer" {...others}>
  {@render children?.()}
</a>
