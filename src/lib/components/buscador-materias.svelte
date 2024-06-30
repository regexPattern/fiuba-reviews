<script lang="ts">
	import { goto } from "$app/navigation";
	import CommandStore from "$lib/command";
	import { Button } from "$lib/components/ui/button";
	import {
		CommandDialog,
		CommandInput,
		CommandItem,
		CommandList
	} from "$lib/components/ui/command";
	import { cn } from "$lib/utils";
	import { Search } from "lucide-svelte";

	export let label: string;

	type Materia = {
		nombre: string;
		codigo: string;
		codigoEquivalencia: string | null;
	};

	export let materias: Materia[];

	const filterQueryMinLength = 2;
	let filtered: Materia[];

	let className = "";
	export { className as class };

	let search = "";
	let debounceTimeout: number | undefined;

	$: if (search.length < filterQueryMinLength || !$CommandStore) {
		filtered = [];
	}

	async function debounceSearch(e: Event) {
		clearTimeout(debounceTimeout);

		debounceTimeout = setTimeout(() => {
			if (e.target instanceof HTMLInputElement) {
				search = e.target.value;

				if (search.length < filterQueryMinLength) {
					return;
				}

				filtered = materias.filter(
					(m) =>
						m.codigo.includes(search) ||
						m.codigoEquivalencia?.includes(search) ||
						m.nombre.includes(search.toUpperCase())
				);
			}
		}, 300);
	}
</script>

<Button
	class={cn("flex justify-between gap-1 px-3 py-2", className)}
	on:click={() => ($CommandStore = !$CommandStore)}
	{...$$restProps}
>
	<span>{label}</span>
	<Search class="h-4 w-4" />
</Button>

<CommandDialog bind:open={$CommandStore} shouldFilter={false}>
	<CommandInput placeholder="Código o nombre de una materia" on:input={debounceSearch} />
	<CommandList>
		{#each filtered as mat (mat.codigo)}
			{@const slug = mat.codigoEquivalencia || mat.codigo}
			<CommandItem
				value={mat.codigo}
				onSelect={async () => {
					await goto(`/materias/${slug}`);
					$CommandStore = false;
				}}
				class="flex cursor-pointer items-start space-x-1.5"
			>
				<span class="font-mono font-semibold">{mat.codigo}</span>
				<span class="font-bold">&bullet;</span>
				<span>
					{mat.nombre}
					{#if mat.codigoEquivalencia}
						(Equivalente a {mat.codigoEquivalencia})
					{/if}
				</span>
			</CommandItem>
		{/each}
	</CommandList>
</CommandDialog>
