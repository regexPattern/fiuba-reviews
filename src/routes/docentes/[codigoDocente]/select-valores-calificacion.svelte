<script lang="ts">
  import * as Form from "$lib/components/ui/form";

  export let id: string;
  export let label: string;
</script>

<div class="grid grid-cols-2">
  <Form.Label for={id} class="flex items-center">{label}</Form.Label>
  <Form.Select>
    <Form.SelectTrigger placeholder="Seleccionar" />
    <Form.SelectContent {id}>
      {#each Array.from({ length: 5 }, (_, i) => i + 1) as val}
        <Form.SelectItem value={val}>{val}</Form.SelectItem>
      {/each}
    </Form.SelectContent>
  </Form.Select>
</div>
