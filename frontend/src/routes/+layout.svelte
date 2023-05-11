<script lang="ts">
	import "../app.css";
</script>

<header class="sticky top-0 z-20 border-b p-6">
	<nav>
		<ul class="flex space-x-4">
			<li><a href="/">Materias</a></li>
			<li><a href="/catedras">Catedras</a></li>
			<li><a href="/docentes">Docentes</a></li>
		</ul>
	</nav>
</header>

<main class="container mx-auto">
	<slot />
</main>
