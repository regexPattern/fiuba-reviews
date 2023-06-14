<script lang="ts">
	import "../app.css";
</script>

<header
	class="sticky top-0 z-20 h-16 border-b bg-white supports-[backdrop-filter]:bg-white/50 supports-[backdrop-filter]:backdrop-blur"
>
	<div class="flex h-full items-center justify-between px-4 lg:container lg:mx-auto">
		<div>Reviews</div>
		<nav>
			<ul class="flex space-x-4">
				<li><a href="/">Inicio</a></li>
				<li><a href="/materias">Materias</a></li>
				<li><a href="/docentes">Docentes</a></li>
			</ul>
		</nav>
	</div>
</header>

<div class="h-full grow pb-6 lg:container lg:mx-auto lg:px-4 lg:py-8">
	<slot />
</div>
