<script lang="ts">
	import "../app.css";
</script>

<header
	class="sticky top-0 z-20 h-16 border-b border-slate-300 text-slate-800 supports-[backdrop-filter]:bg-white/50 supports-[backdrop-filter]:backdrop-blur dark:border-slate-700 dark:text-slate-50 dark:supports-[backdrop-filter]:bg-slate-900/50"
>
	<div class="flex h-full items-center px-4 font-semibold lg:container lg:mx-auto">
		<div class="grow">Reviews</div>
		<nav>
			<ul class="flex space-x-4">
				<li class="transition-colors hover:text-fiuba"><a href="/">Inicio</a></li>
				<li class="transition-colors hover:text-fiuba"><a href="/materias">Materias</a></li>
				<li class="transition-colors hover:text-fiuba"><a href="/docentes">Docentes</a></li>
			</ul>
		</nav>
	</div>
</header>

<div class="h-full grow pb-6 lg:container lg:mx-auto lg:px-4 lg:py-8">
	<slot />
</div>
