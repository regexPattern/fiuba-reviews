<script lang="ts">
  import { Switch } from "$lib/components/ui/switch";
  import type { Switch as SwitchPrimitive } from "bits-ui";
  import { getFormField } from "formsnap";

  type $$Props = SwitchPrimitive.Props;
  type $$Events = SwitchPrimitive.Events;

  export let onCheckedChange: $$Props["onCheckedChange"] = undefined;

  const { name, setValue, attrStore, value } = getFormField();
</script>

<Switch
  {...$attrStore}
  checked={typeof $value === "boolean" ? $value : false}
  onCheckedChange={(v) => {
    onCheckedChange?.(v);
    setValue(v);
  }}
  {...$$restProps}
  on:click
  on:keydown />
<input hidden {name} value={$value} />
