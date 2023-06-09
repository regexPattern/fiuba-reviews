<script>
	import "../app.css";
</script>

<header
	class="sticky top-0 z-20 border-b border-light-border bg-light/80 text-dark backdrop-blur-md dark:border-dark-border dark:bg-dark/80 dark:text-light"
>
	<nav>
		<ul class="flex font-medium">
			<li class="p-4"><a href="/materias">Materias</a></li>
			<li class="p-4"><a href="/docentes">Docentes</a></li>
			<li class="p-4"><a href="/comentar">Comentar</a></li>
		</ul>
	</nav>
</header>
<slot />
