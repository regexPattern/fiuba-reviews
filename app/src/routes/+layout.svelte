<script lang="ts">
	import "../app.css";
</script>

<header class="border-b">
	<nav>
		<ul class="flex">
			<li class="px-8 py-4"><a href="/">Inicio</a></li>
			<li class="px-8 py-4"><a href="/materias">Materias</a></li>
			<li class="px-8 py-4"><a href="/docentes">Docentes</a></li>
		</ul>
	</nav>
</header>

<div class="container mx-auto mt-8">
	<slot />
</div>

<footer />
