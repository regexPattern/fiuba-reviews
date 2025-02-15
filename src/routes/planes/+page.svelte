<script lang="ts">
  import { PUBLIC_TURNSTILE_SITE_KEY } from "$env/static/public";
  import Link from "$lib/components/link.svelte";
  import {
    Alert,
    AlertDescription,
    AlertTitle,
  } from "$lib/components/ui/alert";
  import {
    Form,
    FormButton,
    FormField,
    FormItem,
    Label,
    Textarea,
  } from "$lib/components/ui/form";
  import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
  } from "$lib/components/ui/select";
  import { Toaster } from "$lib/components/ui/sonner";
  import { contenidoSiu as schema } from "$lib/zod/schema";
  import type { PageData } from "./$types";
  import { type FormOptions } from "formsnap";
  import { Loader2 } from "lucide-svelte";
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
  <Alert>
    <AlertTitle class="text-lg"
      >Actualización de oferta de cátedras y docentes</AlertTitle
    >
    <AlertDescription>
      <br />
      La <Link href="https://ofertahoraria.fi.uba.ar/"
        >página de la oferta horaria</Link
      > dejó de actualizarse luego del primer cuatrimestre del 2024, por eso necesito
      de tu ayuda para poder actualizar los listados de cátedras y docentes de las
      materias de la facultad. Si querés colaborar, podés copiar el contenido que
      te aparece en el SIU siguiendo las instrucciones desde una computadora (igual
      que a como se hace en <Link href="https://fede.dm/FIUBA-Plan/"
        >FIUBA Plan</Link
      >):
      <br />
      <br />
      <ol class="ml-2 list-inside list-decimal">
        <li>
          En el SIU, andá a <Link
            href="https://guaraniautogestion.fi.uba.ar/g3w/oferta_comisiones"
            >Reportes > Oferta de comisiones.</Link
          >
        </li>
        <li>
          Seleccioná todo el contenido de la página <kbd>(CTRL/CMD + A)</kbd>.
        </li>
        <li>Copia todo <kbd>(CTRL/CMD + C)</kbd>.</li>
        <li>
          Pegalo en el cuadro de texto de abajo <kbd>(CTRL/CMD + V)</kbd>.
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
    </AlertDescription>
  </Alert>

  <Form
    method="POST"
    form={data.form}
    {schema}
    options={formOptions}
    let:config
    let:formValues
    let:submitting
    let:errors
    class="space-y-6 xs:space-y-8"
  >
    <FormField {config} name="carrera">
      <div class="grid sm:grid-cols-2 gap-4">
        <Label for="carrera" class="flex items-center"
          >Seleccioná tu carrera</Label
        >
        <Select>
          <SelectTrigger>
            <SelectValue placeholder="Seleccionar" />
          </SelectTrigger>
          <SelectContent id="carrera">
            {#each data.carrerasFaltantes as carrera}
              <SelectItem value={carrera}>{carrera}</SelectItem>
            {/each}
            {#each data.planesRegistrados as plan}
              <SelectItem value={plan.carrera} disabled
                ><span
                  ><span class="line-through">{plan.carrera}</span> (Ya actualizada)</span
                ></SelectItem
              >
            {/each}
          </SelectContent>
          {#if errors.carrera}
            {errors.carrera}
          {/if}
        </Select>
      </div>
    </FormField>

    <FormField {config} name="contenido-siu">
      <FormItem>
        <Label for="contenido-siu">Contenido copiado del SIU</Label>
        <Textarea id="contenido-siu" />
      </FormItem>
    </FormField>

    <Turnstile siteKey={PUBLIC_TURNSTILE_SITE_KEY} theme={$mode} />

    <FormButton type="submit" class="items-center gap-1" disabled={submitting}>
      {#if submitting}
        <span>Enviando</span>
        <Loader2 class="h-4 w-4 animate-spin" />
      {:else}
        Enviar
      {/if}
    </FormButton>

    <!-- <SuperDebug data={formValues} /> -->
  </Form>
</main>
