# Migración de `sitio-web` a Next.js (Vercel)

Reescritura desde cero en Next.js 16 con **Cache Components** (`cacheComponents`, `"use cache"`,
`cacheTag`, `<Suspense>`). No se mantiene compatibilidad hacia atrás.

## Hecho (resumen)

Fases 0 a 6 implementadas en la rama `next`; `next build` OK (`/materia/[codigo]` sale como
Partial Prerender) y el servidor de producción fue probado a mano.

- **Stack**: Next 16 + React 19, Tailwind v4, shadcn/ui, `next-themes`, `next/font`, zod,
  `@t3-oss/env-nextjs` (`src/env.ts`), Drizzle + `postgres`. El SvelteKit viejo está en
  `../sitio-web-legacy/`.
- **Caché**: queries en `src/server/queries/` (`getMateriasVigentes`, `getMateria`,
  `getCatedrasMateria`, `getCatedra`, `getInicio`), todas con `"use cache"`, `cacheLife("max")`
  (inicio: `"hours"`) y tags `materias`, `materia:<c>`, `catedras:<c>`, `catedra:<uuid>`, `inicio`.
- **Rutas**: inicio, `/materia/[codigo]?catedra=<uuid>`, `/calificar`, `not-found`/`error`, OG
  images (inicio y materia), `sitemap`, `robots`, PostHog (`instrumentation-client.ts`).
- **Revalidación**: `POST /api/revalidate` con `Authorization: Bearer <REVALIDATE_SECRET>` y body
  `{ materiasCambiaron?, materias?, catedras?, tags?, inicio? }` (usa `revalidateTag(tag, "max")`).
- **Calificar**: Server Action con transacción + Turnstile + `updateTag` de las cátedras del
  docente, `catedras:<materia del docente>` e `inicio`. El link desde la cátedra incluye `?materia=`
  para volver a la materia vigente.

## Pendiente

### Fase 7 — Verificación y corte

- [ ] Probar con `next build && next start` el flujo completo de calificar: el detalle y el
      promedio del listado cambian y no se regeneran otras cátedras.
- [ ] Probar que invalidar una materia vieja actualiza el listado de la materia vigente que la
      muestra por equivalencia.
- [ ] Preview deploy en Vercel: revisar `x-vercel-cache` y TTFB en frío y en caliente.
- [ ] Lighthouse y Core Web Vitals en mobile.
- [ ] Variables de entorno en Vercel (ver `.env.development`); sacar `ISR_BYPASS_TOKEN`.
- [ ] Merge a `main`, borrar `sitio-web-legacy/` y actualizar `README.md` y `CONTRIBUTING.md`.

## Extras (al final de todo)

- [ ] **Dump local desactualizado**: `fiuba-reviews.sql` usa `catedra.esta_activa` pero el schema
      (y producción) usa `activa`. Regenerar el dump o aplicar el rename. Mientras tanto
      `.env.development` apunta a `fiuba-reviews` y falla; se probó contra una copia `fiuba_next`.
- [ ] **Región de Vercel**: configurar las funciones en la misma región que la DB y revisar el
      límite de conexiones del proveedor.
- [ ] **`actualizador-siu`**: que llame a `/api/revalidate` al terminar una corrida, con los
      códigos de materia (y/o cátedras) que cambiaron.
- [ ] **404 en materias inexistentes**: la primera visita a un código que no existe devuelve 200
      (con `noindex` y la página de 404); las siguientes devuelven 404. Evaluar `dynamicParams`.
- [ ] Revisar el warning de `next/font` por falta de métricas de `Google Sans Code`.

## Referencias

- Cache Components: https://nextjs.org/docs/app/api-reference/config/next-config-js/cacheComponents
- `use cache` / `cacheTag` / `cacheLife`: https://nextjs.org/docs/app/api-reference/directives/use-cache
- ISR con Cache Components: https://nextjs.org/docs/app/guides/incremental-static-regeneration
- shadcn/ui: https://ui.shadcn.com/docs/installation/next
