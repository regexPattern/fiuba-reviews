<script lang="ts">
  import type { PageData } from "./$types";
  import type { FormOptions } from "formsnap";

  import * as Form from "$lib/components/ui/form";
  import * as Select from "$lib/components/ui/select";

  import { PUBLIC_TURNSTILE_SITE_KEY } from "$env/static/public";
  import InputCalificacion from "$lib/components/input-calificacion.svelte";
  import Link from "$lib/components/link.svelte";
  import { Toaster } from "$lib/components/ui/sonner";
  import { codigoDocente as schema } from "$lib/zod/schema";
  import { ChevronLeft, Loader2 } from "lucide-svelte";
  import { mode } from "mode-watcher";
  import { toast } from "svelte-sonner";
  import { Turnstile } from "svelte-turnstile";
  import SuperDebug from "sveltekit-superforms/client/SuperDebug.svelte";

  export let data: PageData;

  const formOptions: FormOptions<typeof schema> = {
    resetForm: true,
    onUpdated: ({ form }) => {
      if (form.valid) {
        toast.success(form.message);
      } else {
        for (const e of form.errors?._errors || []) {
          toast.error(e);
        }
      }
    },
    onError: "apply",
  };
</script>

<Toaster />

<main class="mx-auto max-w-screen-sm space-y-6 p-4 xs:space-y-8">
  <Link
    href={`/materias/${data.docente.codigo_materia}/${data.docente.codigo_catedra}`}
    class="flex items-center gap-1 underline"
  >
    <ChevronLeft class="w-4" />
    Ir a cátedra del docente
  </Link>

  <h1 class="text-5xl font-bold tracking-tight">{data.docente.nombre}</h1>

  <Form.Root
    method="POST"
    form={data.form}
    {schema}
    options={formOptions}
    let:config
    let:formValues
    let:submitting
    class="space-y-6 xs:space-y-8"
  >
    <div class="space-y-4">
      <InputCalificacion
        name="acepta-critica"
        label="Acepta Crítica"
        {config}
      />
      <InputCalificacion name="asistencia" label="Asistencia" {config} />
      <InputCalificacion name="buen-trato" label="Buen Trato" {config} />
      <InputCalificacion name="claridad" label="Claridad" {config} />
      <InputCalificacion
        name="clase-organizada"
        label="Clase Organizada"
        {config}
      />
      <InputCalificacion
        name="cumple-horario"
        label="Cumple Horario"
        {config}
      />
      <InputCalificacion
        name="fomenta-participacion"
        label="Fomenta Participación"
        {config}
      />
      <InputCalificacion
        name="panorama-amplio"
        label="Panorama Amplio"
        {config}
      />
      <InputCalificacion
        name="responde-mails"
        label="Responde Mails"
        {config}
      />
    </div>

    <div class="space-y-4">
      <Form.Field {config} name="comentario">
        <Form.Item>
          <Form.Label for="comentario">Comentario (Opcional)</Form.Label>
          <Form.Textarea id="comentario" />
        </Form.Item>
      </Form.Field>

      <Form.Field {config} name="cuatrimestre">
        <Form.Item>
          <Form.Label for="cuatrimestre">Cuatrimestre</Form.Label>
          <Form.Select disabled={!formValues.comentario}>
            <Form.SelectTrigger>
              <Select.Value placeholder="Seleccionar" />
            </Form.SelectTrigger>
            <Form.SelectContent id="cuatrimestre">
              {#each data.cuatrimestres as cuatri}
                <Form.SelectItem value={`${cuatri.anio}-${cuatri.numero}`}
                  >{cuatri.numero}C {cuatri.anio}</Form.SelectItem
                >
              {/each}
            </Form.SelectContent>
          </Form.Select>
        </Form.Item>
      </Form.Field>
    </div>

    <Turnstile siteKey={PUBLIC_TURNSTILE_SITE_KEY} theme={$mode} />

    <Form.Button type="submit" class="items-center gap-1" disabled={submitting}>
      {#if submitting}
        <span>Enviando</span>
        <Loader2 class="h-4 w-4 animate-spin" />
      {:else}
        Enviar
      {/if}
    </Form.Button>

    {#if import.meta.env.DEV}
      <SuperDebug data={formValues} />
    {/if}
  </Form.Root>
</main>
