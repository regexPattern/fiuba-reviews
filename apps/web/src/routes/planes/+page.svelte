<script lang="ts">
  import type { PageData } from "./$types";

  import { PUBLIC_TURNSTILE_SITE_KEY } from "$env/static/public";

  import { formPlanSiu as schema } from "$lib/zod/schema";
  import { Loader2 } from "lucide-svelte";
  import { mode } from "mode-watcher";
  import { toast } from "svelte-sonner";
  import { type FormOptions } from "formsnap";

  import { Turnstile } from "svelte-turnstile";
  import SuperDebug from "sveltekit-superforms/client/SuperDebug.svelte";

  import * as Alert from "$lib/components/ui/alert";
  import * as Form from "$lib/components/ui/form";
  import Link from "$lib/components/link.svelte";
  import { Toaster } from "$lib/components/ui/sonner";

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
  };
</script>

<Toaster />

<main class="mx-auto max-w-screen-sm space-y-6 p-4 xs:space-y-8">
  <Alert.Root>
    <Alert.Title class="text-xl"
      >Actualización de oferta de cátedras y docentes</Alert.Title>
    <Alert.Description class="text-md mt-6">
      La <Link href="https://ofertahoraria.fi.uba.ar/"
        >página de la oferta horaria</Link> dejó de actualizarse luego del primer
      cuatrimestre del 2024, por eso necesito de tu ayuda para poder actualizar los
      listados de cátedras y docentes de las materias de la facultad. Si querés colaborar,
      podés copiar el contenido que te aparece en el SIU siguiendo las instrucciones
      desde una computadora (igual que a como se hace en <Link
        href="https://fede.dm/FIUBA-Plan/">FIUBA Plan</Link
      >):
      <br />
      <br />
      <ol class="ml-2 list-inside list-decimal">
        <li>
          En el SIU, andá a <Link
            href="https://guaraniautogestion.fi.uba.ar/g3w/oferta_comisiones"
            >Reportes > Oferta de comisiones.</Link>
        </li>
        <li>
          Seleccioná todo el contenido de la página <kbd
            class="hidden rounded border px-1 py-0.5 font-mono text-sm tracking-widest sm:inline">
            (CTRL/CMD + A)</kbd
          >.
        </li>
        <li>
          Copia todo <kbd
            class="hidden rounded border px-1 py-0.5 font-mono text-sm tracking-widest sm:inline">
            (CTRL/CMD + C)</kbd
          >.
        </li>
        <li>
          Pegalo en el cuadro de texto de abajo <kbd
            class="hidden rounded border px-1 py-0.5 font-mono text-sm tracking-widest sm:inline">
            (CTRL/CMD + V)</kbd
          >.
        </li>
      </ol>
      <br />
      En la lista desplegable de selección de carrera de abajo podés ver cuáles carreras
      hacen falta y cuáles ya fueron actualizadas por otros alumnos. Si estudias
      alguna de las que no se han enviado, se agredece si podés hacer el aporte para
      terminar de actualizar los listados.
      <br />
      <br />
      La actualización no va a ser inmediata para este cuatrimestre pero apunto a
      que esté lista a partir del próximo cuatrimestre.
      <br />
      <br />
      Muchas gracias.
    </Alert.Description>
  </Alert.Root>

  <Form.Root
    method="POST"
    form={data.form}
    {schema}
    options={configFormulario}
    let:config
    let:formValues
    let:submitting
    let:errors
    class="space-y-6 xs:space-y-8">
    <Form.Field {config} name="carrera">
      <div class="grid gap-4 sm:grid-cols-2">
        <Form.Label for="carrera" class="flex items-center"
          >Seleccioná tu carrera</Form.Label>
        <Form.Select>
          <Form.SelectTrigger placeholder="Seleccionar" />
          <Form.SelectContent id="carrera">
            {#each data.carrerasFaltantes as carrera}
              <Form.SelectItem value={carrera}>{carrera}</Form.SelectItem>
            {/each}
            {#each data.planesRegistrados as plan}
              <Form.SelectItem value={plan.carrera} disabled
                ><span
                  ><span class="line-through">{plan.carrera}</span> (Ya enviada)</span
                ></Form.SelectItem>
            {/each}
          </Form.SelectContent>
          {#if errors.carrera}
            {errors.carrera}
          {/if}
        </Form.Select>
      </div>
    </Form.Field>

    <Form.Field {config} name="contenido-siu">
      <Form.Item>
        <Form.Label for="contenido-siu">Contenido copiado del SIU</Form.Label>
        <Form.Textarea id="contenido-siu" />
      </Form.Item>
    </Form.Field>

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
