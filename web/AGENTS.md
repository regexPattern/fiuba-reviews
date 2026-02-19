# Runtime de JavaScript

Utilizá Bun como runtime de JavaScript.

# Consideraciones de diseño

- La página src/routes/materia/[codigo_materia]/[codigo_catedra]/+page.svelte realmente no tiene contenido que mostrar, sino que simplemente redirige hacia la primera cátedra de cada materia.

# Convenciones de código

- Mantené el código lo más simple posible. Podés detallar consideraciones a futuro para futuros cambios, pero dejalos fuera de la implementación a menos que se te solicited explícitamente.
- Preferí utilizar comillas dobles al momento de escribir atributos/props de tipo texto, pero para aquellos que necesiten evaluar alguna expresión de JavaScript, utilizá llaves ({}). En el caso en el que el texto dobles dentro, por ejemplo, cuando se agregan dichas comillas con before y after en CSS usando Tailwind, utiliza llaves y backticks como expresión JavaScript para una string. Por ejemplo:

Prefiri:
```
<div atributo={`before:content-['"']`}></div>
```

Por sobre esto:
```
<div atributo="before:content-['\"']"}></div>
```

O esto:
```
<div atributo="before:content-['&quot;']"></div>
```
