<script lang="ts">
	import { AutoScroll } from "@splidejs/splide-extension-auto-scroll";
	import { Splide, SplideSlide } from "@splidejs/svelte-splide";

	import Link from "./link.svelte";

	export let materias: {
		codigo: number;
		nombre: string;
		cantidadCatedras: number;
		cantidadComentarios: number;
	}[];
</script>

<Splide
	options={{
		autoWidth: true,
		arrows: false,
		pagination: false,
		gap: "1.5rem",
		type: "loop",
		padding: { top: "1rem", bottom: "1rem" },
		autoScroll: {
			speed: 0.5
		}
	}}
	extensions={{ AutoScroll }}
>
	{#each materias as mat (mat.codigo)}
		<SplideSlide
			class="m-0.5 rounded-lg border bg-slate-50 text-center text-sm text-muted-foreground dark:bg-slate-900"
		>
			<Link href={`/materias/${mat.codigo}`} class="flex flex-col p-3">
				<span class="font-medium text-foreground">{mat.codigo} &bull; {mat.nombre}</span>
				<span>{mat.cantidadCatedras} cátedras</span>
				<span>{mat.cantidadComentarios} comentarios</span>
			</Link>
		</SplideSlide>
	{/each}
</Splide>
