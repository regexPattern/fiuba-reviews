# Convenciones

## Comentarios

- Siempre con `/* */`, incluso los de una línea. En JSX, `{/* */}`.
- Los docstrings de funciones arrancan con un verbo: "Devuelve…", "Formatea…", "Indica si…".
- No comentar lo que el nombre o el código ya explican, sobre todo en código interno. Dentro de una función, comentar el porqué y no el qué.
- Cortos e impersonales ("Se invalida…", no "Invalidamos…").

## Nombres

- Identificadores de dominio en castellano. Los nombres estándar de React y de UI van en inglés (`label`, `rootProps`, `formAction`, `Shortcut`). Lo mismo para convensiones informáticas que se suelen nombrar en inglés.
- Callbacks en inglés: `onX` para las props, `handleX` para las funciones que las manejan.
- Booleanos (funciones y variables) con `es…`/`esta…` (`esPosterior`, `estaAutorizado`), salvo que ya se lean como pregunta (`usaCmd`).
- Constantes en `SCREAMING_CASE`, con prefijo para los límites: `MAX_PUNTAJE`, `MIN_PUNTAJE`.
- Schemas de zod en PascalCase con prefijo `Schema` (`SchemaPuntaje`). Las funciones que los crean arrancan con un verbo (`crearSchemaCalificacion`).

## Organización

- Solo se testea la lógica de dominio, no la UI ni los helpers de formato.
- Los logs van con `logger` (pino) de `@/server/logger`: `log.info({ ...campos }, "Mensaje.")`.

## Verificación

- `pnpm typecheck` y `pnpm test`.
