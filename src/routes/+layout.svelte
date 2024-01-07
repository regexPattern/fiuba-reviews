<script lang="ts">
	import Link from "$lib/components/link.svelte";
	import ThemePicker from "$lib/components/theme-picker.svelte";
	import Button from "$lib/components/ui/button/button.svelte";
	import { Search } from "lucide-svelte";
	import { ModeWatcher } from "mode-watcher";

	import "../app.css";
</script>

<svelte:head>
	<title>FIUBA Reviews</title>
</svelte:head>

<ModeWatcher />

<header class="sticky top-0 z-[40] border-b bg-background/75 backdrop-blur-lg">
	<div class="container flex h-16 items-center justify-start gap-2 p-3 sm:justify-between">
		<Link href="/" class="mr-2 flex-1 font-semibold sm:mr-auto sm:p-0">
			<span class="hidden sm:inline">FIUBA<span class="font-medium">&nbsp;Reviews</span></span>
			<span class="sm:hidden">F<span class="font-medium">R</span></span>
		</Link>

		<nav>
			<ul>
				<li class="flex items-center">
					<Button href="/materias" variant="outline" class="space-x-1">
						<span>Materias</span>
						<Search class="h-4 w-4" />
					</Button>
				</li>
			</ul>
		</nav>

		<ThemePicker class="h-10 w-10 shrink-0" />

		<Link
			href="https://github.com/regexPattern/fiuba-reviews"
			class="flex h-10 w-10 shrink-0 items-center justify-center"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				x="0px"
				y="0px"
				viewBox="0 0 24 24"
				class="h-[26px] w-[26px] fill-black dark:fill-white"
			>
				<path
					d="M10.9,2.1c-4.6,0.5-8.3,4.2-8.8,8.7c-0.5,4.7,2.2,8.9,6.3,10.5C8.7,21.4,9,21.2,9,20.8v-1.6c0,0-0.4,0.1-0.9,0.1 c-1.4,0-2-1.2-2.1-1.9c-0.1-0.4-0.3-0.7-0.6-1C5.1,16.3,5,16.3,5,16.2C5,16,5.3,16,5.4,16c0.6,0,1.1,0.7,1.3,1c0.5,0.8,1.1,1,1.4,1 c0.4,0,0.7-0.1,0.9-0.2c0.1-0.7,0.4-1.4,1-1.8c-2.3-0.5-4-1.8-4-4c0-1.1,0.5-2.2,1.2-3C7.1,8.8,7,8.3,7,7.6C7,7.2,7,6.6,7.3,6 c0,0,1.4,0,2.8,1.3C10.6,7.1,11.3,7,12,7s1.4,0.1,2,0.3C15.3,6,16.8,6,16.8,6C17,6.6,17,7.2,17,7.6c0,0.8-0.1,1.2-0.2,1.4 c0.7,0.8,1.2,1.8,1.2,3c0,2.2-1.7,3.5-4,4c0.6,0.5,1,1.4,1,2.3v2.6c0,0.3,0.3,0.6,0.7,0.5c3.7-1.5,6.3-5.1,6.3-9.3 C22,6.1,16.9,1.4,10.9,2.1z"
				/>
			</svg>
		</Link>
	</div>
</header>

<slot />
