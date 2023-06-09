<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;

	let search = "";

	$: filtered_materias = data.materias.filter((m) => {
		return m.search_terms.some((term) => term?.toLowerCase().includes(search));
	});
</script>

<input bind:value={search} placeholder="Buscar materias" />
<ul>
	{#each filtered_materias as m}
		<li>
			<a href={`/materias/${m.codigo_equivalencia ?? m.codigo}/catedras`}
				>{m.codigo} - {m.nombre}
				{#if m.codigo_equivalencia}
					-> {m.codigo_equivalencia}
				{/if}
			</a>
		</li>
	{/each}
</ul>
