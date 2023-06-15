<script lang="ts">
	import "../app.css";
</script>

<header class="border-color sticky-bg-color sticky top-0 z-20 h-16 border-b">
	<div class="flex h-full items-center px-4 font-semibold lg:container lg:mx-auto">
		<div class="grow text-sm leading-6 text-slate-700 dark:text-slate-50">Reviews</div>
		<nav class="text-sm leading-6 text-slate-700 dark:text-slate-50">
			<ul class="flex space-x-4">
				<li><a href="/">Inicio</a></li>
				<li><a href="/materias">Materias</a></li>
				<li><a href="/docentes">Docentes</a></li>
			</ul>
		</nav>
	</div>
</header>

<div class="mb-4 h-full grow lg:container lg:mx-auto lg:my-8">
	<slot />
</div>

<style lang="postcss">
	nav > ul > li {
		@apply transition-colors hover:text-fiuba;
	}
</style>
