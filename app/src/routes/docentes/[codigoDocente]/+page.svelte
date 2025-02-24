<script lang="ts">
  import type { FormOptions } from "formsnap";
  import type { PageData } from "./$types";

  import { PUBLIC_TURNSTILE_SITE_KEY } from "$env/static/public";

  import { formCalificacionDocente as schema } from "$lib/zod/schema";
  import { mode } from "mode-watcher";
  import { toast } from "svelte-sonner";

  import { ChevronLeft, Loader2 } from "lucide-svelte";
  import { Turnstile } from "svelte-turnstile";
  import SuperDebug from "sveltekit-superforms/client/SuperDebug.svelte";

  import * as Form from "$lib/components/ui/form";
  import Link from "$lib/components/link.svelte";
  import { Toaster } from "$lib/components/ui/sonner";
  import SelectValoresCalificacion from "./select-valores-calificacion.svelte";

  export let data: PageData;

  const configFormulario: FormOptions<typeof schema> = {
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
    multipleSubmits: "prevent",
  };
</script>

<Toaster />

<main class="mx-auto max-w-screen-sm space-y-6 p-4 xs:space-y-8">
  <Link
    href={`/materias/${data.docente.codigo_materia}/${data.docente.codigo_catedra}`}
    class="flex items-center gap-1 underline">
    <ChevronLeft class="w-4" />
    Ir a cátedra del docente
  </Link>

  <h1 class="text-5xl font-bold tracking-tight">{data.docente.nombre}</h1>

  <Form.Root
    method="POST"
    form={data.form}
    {schema}
    options={configFormulario}
    let:config
    let:formValues
    let:submitting
    class="space-y-6 xs:space-y-8">
    <div class="space-y-4">
      <Form.Field {config} name="acepta-critica">
        <SelectValoresCalificacion id="acepta-critica" label="Acepta Crítica" />
      </Form.Field>
      <Form.Field {config} name="asistencia">
        <SelectValoresCalificacion id="asistencia" label="Asistencia" />
      </Form.Field>
      <Form.Field {config} name="buen-trato">
        <SelectValoresCalificacion id="buen-trato" label="Buen Trato" />
      </Form.Field>
      <Form.Field {config} name="claridad">
        <SelectValoresCalificacion id="claridad" label="Claridad" />
      </Form.Field>
      <Form.Field {config} name="clase-organizada">
        <SelectValoresCalificacion
          id="clase-organizada"
          label="Clase Organizada" />
      </Form.Field>
      <Form.Field {config} name="cumple-horario">
        <SelectValoresCalificacion id="cumple-horario" label="Cumple Horario" />
      </Form.Field>
      <Form.Field {config} name="fomenta-participacion">
        <SelectValoresCalificacion
          id="fomenta-participacion"
          label="Fomenta Participación" />
      </Form.Field>
      <Form.Field {config} name="panorama-amplio">
        <SelectValoresCalificacion
          id="panorama-amplio"
          label="Panorama Amplio" />
      </Form.Field>
      <Form.Field {config} name="responde-mails">
        <SelectValoresCalificacion id="responde-mails" label="Responde Mails" />
      </Form.Field>
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
          <Form.Select disabled={!formValues.comentario?.trim()}>
            <Form.SelectTrigger placeholder="Seleccionar" value={0} />
            <Form.SelectContent id="cuatrimestre">
              {#each data.cuatrimestres as cuatri}
                <Form.SelectItem value={cuatri.codigo}
                  >{cuatri.numero}C {cuatri.anio}</Form.SelectItem>
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
