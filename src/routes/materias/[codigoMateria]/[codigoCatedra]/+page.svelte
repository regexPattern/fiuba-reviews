<script lang="ts">
	import Link from "$lib/components/link.svelte";
	import Promedios from "$lib/components/listado-promedios-docente.svelte";
	import PlaceholderCatedra from "$lib/components/placeholder-catedra.svelte";
	import { Button } from "$lib/components/ui/button";
	import { Popover, PopoverContent, PopoverTrigger } from "$lib/components/ui/popover";
	import { ChevronDown, PlusCircle, Star } from "lucide-svelte";

	import type { PageData } from "./$types";

	export let data: PageData;
</script>

{#await data.streamed.docentes}
	<PlaceholderCatedra />
{:then docentes}
	{#each docentes as doc (doc.codigo)}
		<!-- TODO: Tengo que hacer que los ids de los docentes tengan alguna utilidad para ser usados como links. -->
		{@const urlSafeName = encodeURIComponent(doc.nombre.toLowerCase())}
		<section id={urlSafeName} class="space-y-4">
			<h2 class="text-4xl font-bold tracking-tight">{doc.nombre}</h2>

			{#if doc.resumen_comentarios}
				<div
					class="divide-y divide-border rounded-lg border bg-secondary dark:divide-slate-700 dark:border-slate-700 [&>*]:p-3"
				>
					<p class={`text-secondary-foreground before:content-['"'] after:content-['"']`}>
						{doc.resumen_comentarios}
					</p>
					<div class="text-sm text-slate-500">
						Resumen generado por IA.
						<Link
							href="https://github.com/regexPattern/che-fiuba#resumen-de-comentarios-con-inteligencia-artificial"
							class="after:content-link"
						>
							Más información.
						</Link>
					</div>
				</div>
			{/if}

			<div class="flex flex-col gap-2 xs:flex-row xs:items-center">
				{#if doc.promedio}
					<Popover>
						<PopoverTrigger asChild let:builder>
							<Button builders={[builder]} variant="outline" class="items-center gap-1.5">
								<Star class="h-4 w-4 fill-current text-yellow-500" />
								<span>Promedio: {doc.promedio.toFixed(1)}</span>
								<ChevronDown class="h-[1.2rem] w-[1.2rem]" />
							</Button>
						</PopoverTrigger>
						<PopoverContent class="w-max">
							<Promedios cantidadCalificaciones={doc.cantidadCalificaciones} {...doc.promedios} />
						</PopoverContent>
					</Popover>
				{:else}
					<Button variant="outline" class="items-center gap-1.5">
						<Star class="h-4 w-4 fill-none text-yellow-500" />
						<span>Sin calificaciones</span>
					</Button>
				{/if}

				<Button class="items-center gap-1.5" href={`/docentes/${doc.codigo}`}>
					Calificar <PlusCircle class="h-[1.2rem] w-[1.2rem]" />
				</Button>
			</div>

			<div class="flex flex-col gap-2 divide-y">
				{#each doc.comentarios as com (com.codigo)}
					<div class="pt-2 [&:first-child]:pt-0">
						<p class={`inline before:content-['"'] after:content-['"']`}>
							{com.contenido}
						</p>
						&dash;
						<span class="text-sm text-muted-foreground">{com.cuatrimestre}</span>
					</div>
				{/each}
			</div>
		</section>
	{/each}
{/await}
