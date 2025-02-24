<script lang="ts">
  export let icon: ConstructorOfATypedSvelteComponent;
  export let title: string;
  export let desc: string;
</script>

<div class="flex items-start gap-4">
  <div class="rounded-lg border-2 border-fiuba bg-background p-1.5 shadow">
    <svelte:component this={icon} class="stroke-fiuba" />
  </div>
  <div class="space-y-3">
    <span class="font-medium">{title}</span>
    <p class="text-muted-foreground">{desc}</p>
  </div>
</div>
