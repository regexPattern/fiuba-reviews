<script lang="ts">
	import { Label } from "$lib/components/ui/label";
	import { cn } from "$lib/utils";
	import type { Label as LabelPrimitive } from "bits-ui";
	import { getFormField } from "formsnap";

	type $$Props = LabelPrimitive.Props;

	let className: $$Props["class"] = undefined;
	export { className as class };

	const { errors, ids } = getFormField();
</script>

<Label
	for={$ids.input}
	class={cn($errors && "text-destructive underline decoration-wavy underline-offset-4", className)}
	{...$$restProps}
>
	<slot />
</Label>
