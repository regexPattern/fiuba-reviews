<script lang="ts">
  import { Input, type InputEvents } from "$lib/components/ui/input";
  import { getFormField } from "formsnap";
  import type { HTMLInputAttributes } from "svelte/elements";

  type $$Props = HTMLInputAttributes;
  type $$Events = InputEvents;

  const { attrStore, value } = getFormField();
</script>

<Input
  {...$attrStore}
  bind:value={$value}
  {...$$restProps}
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
  on:input />
