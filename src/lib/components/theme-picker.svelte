<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import {
		DropdownMenu,
		DropdownMenuContent,
		DropdownMenuItem,
		DropdownMenuTrigger
	} from "$lib/components/ui/dropdown-menu";
	import { Moon, Sun } from "lucide-svelte";
	import { setMode } from "mode-watcher";

	let className = "";
	export { className as class };
</script>

<DropdownMenu positioning={{ placement: "bottom-end" }}>
	<DropdownMenuTrigger asChild let:builder>
		<Button builders={[builder]} variant="outline" size="icon" class={className}>
			<Sun
				class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0"
			/>
			<Moon
				class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100"
			/>
		</Button>
	</DropdownMenuTrigger>
	<DropdownMenuContent class="z-[50]">
		<DropdownMenuItem on:click={() => setMode("light")}>Claro</DropdownMenuItem>
		<DropdownMenuItem on:click={() => setMode("dark")}>Oscuro</DropdownMenuItem>
		<DropdownMenuItem on:click={() => setMode("system")}>Dispositivo</DropdownMenuItem>
	</DropdownMenuContent>
</DropdownMenu>
