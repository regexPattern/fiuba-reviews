# Guía para agentes de código

Este documento describe cómo deben interactuar los agentes que escriben código en este proyecto, con qué herramientas cuentan y qué estilo deben respetar.

## 1. Stack tecnológico a tener en cuenta

Los agentes deben entender que el proyecto está construido sobre el siguiente stack:

- **JavaScript/TypeScript** ejecutado en **Bun**
- **Svelte 5 + SvelteKit 2** como framework frontend
- **bits-ui** como librería de componentes UI
- **TailwindCSS** para estilos
- **Postgres** como base de datos, utilizando **drizzle** como ORM

## 2. Estilo y convenciones de código

Para garantizar coherencia y mantenibilidad, seguir estas reglas:

### 2.1 Simplicidad y legibilidad

- Optimizar para la claridad antes que para la “cleverness”.
- Evitar pasos innecesarios o código demasiado indirecto.
- Mantener bajo el número de variables intermedias salvo que aporten claridad.
- Para condicionales o bucles de una sola línea, siempre incluí las llaves.

### 2.2 Nombres de variables

- Usar **nombres en español** cuando sea natural.
- Usar inglés sólo cuando tenga más sentido (ej: `buffer`, `tmp`, `payload`, `slug`).
- Evitar abreviaturas innecesarias.

### 2.3 Formato del código

Este proyecto usa `prettier`. Formatear antes de finalizar:

```sh
bunx prettier --write ./path/to/file
```

## 3. Herramientas MCP disponibles

Los agentes tienen acceso a distintas herramientas que pueden usar durante la generación de código. Se describen a continuación:

### 3.1 Herramientas para Svelte

Los agentes cuentan con un MCP especializado para Svelte 5 + SvelteKit. Las herramientas funcionan así:

#### list-sections

- Usar siempre primero
- Muestra la lista de secciones de documentación disponibles
- Ayuda a ubicar material relevante según contexto

#### get-documentation

- Se usa después de list-sections
- Recupera el contenido completo de las secciones pertinentes
- Puede recibir una o varias secciones a la vez
- Elegir las secciones basándose en use_cases

#### svelte-autofixer

- Debe ejecutarse cada vez que un agente escriba código Svelte
- Repetir hasta que no queden sugerencias ni errores
- Sólo entonces el código se considera entregable
- Nota: Los agentes deben usar esta pipeline al recibir cualquier instrucción relacionada a Svelte o SvelteKit.

### 3.2 Herramientas para bits-ui

La documentación de bits-ui está disponible vía webfetch en https://bits-ui.com/llms.txt. Ese endpoint incluye:

- Índice de componentes
- Descripción general de la librería
- Rutas para acceder a documentación específica

bits-ui se usa como librería UI principal en combinación con Tailwind.

Cuando utilices los ejemplos de la documentación, mantene solo los estilos necesarios para tener un componentes básico. En la documentación hay bastantes ejemplos con clases con propiedades o variables custom de tamaño o de color, etc. Ignoralas, mantene los estilos mínimos.

## 3.3 Herramientas para Postgres

Los agentes cuentan con acceso a Postgres en modo read-only a través del MCP. Usos recomendados:

- Consultar el schema
- Entender relaciones y tablas
- Ver datos relevantes
- Analizar cómo drizzle modela la base

No se puede escribir ni mutar datos desde esta herramienta.

# 4. Principios finales para agentes

Antes de finalizar cualquier respuesta:

- Validar coherencia con el stack indicado
- Respetar estilo y convenciones de este documento
- Minimizar complejidad innecesaria
- Entregar código formateado con prettier

En caso de trabajar con Svelte:

- list-sections
- get-documentation
- svelte-autofixer
