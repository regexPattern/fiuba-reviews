<script lang="ts">
  export let icon: ConstructorOfATypedSvelteComponent;
  export let title: string;
  export let desc: string;
</script>

<div class="flex items-start gap-4">
  <div class="rounded-lg border-2 border-purple-500 bg-background p-1.5 shadow">
    <svelte:component this={icon} class="stroke-purple-500" />
  </div>
  <div class="space-y-3">
    <span class="font-medium">{title}</span>
    <p class="text-muted-foreground">{desc}</p>
  </div>
</div>
