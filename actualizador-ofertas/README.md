# Actualizador de ofertas

Planificación de la reescritura del actualizador de FIUBA Reviews. Este documento describe el problema y el flujo previsto; **no es una implementación**.

La aplicación recibe ofertas académicas parseadas, las compara con la base de datos y genera una **migración SQL revisable**. Generar la migración no debe modificar la base: su aplicación es un paso separado.

## Conceptos y reglas del dominio

- **Plan de estudios:** conjunto de materias de una carrera, con su año y equivalencias. No es lo mismo que una oferta cuatrimestral.
- **Oferta:** materias y comisiones publicadas para una carrera y un cuatrimestre. Es lo que ingresa desde el parser del SIU.
- **Materia:** se identifica por su código oficial, sin distinguir carrera. Un mismo código presente en varias carreras representa la misma materia.
- **Docente:** hoy la DB registra docentes **por materia**, con un UUID, un nombre de presentación y, cuando se conoce, `nombre_siu`. Sus comentarios y calificaciones deben conservarse al reconocer su identidad.
- **Cátedra/comisión:** para este actualizador, una comisión distinguible y su grupo de docentes. Dos comisiones A/B o 1/2 son diferentes aunque tengan exactamente los mismos docentes. Agregar o quitar un docente también da lugar a una cátedra diferente para fines prácticos.
- Las comisiones de **condicionales no participan** de la actualización de cátedras.
- Para una materia se utiliza el cuatrimestre más reciente disponible. Las ofertas de ese mismo cuatrimestre se unen entre carreras; las de cuatrimestres anteriores no se mezclan con ellas.
- Desactivar no significa borrar: las materias históricas, docentes, reseñas y cátedras anteriores deben seguir existiendo.

## Flujo general

```text
1. Ingress de ofertas parseadas + lectura del estado de la DB
                         ↓
2. Sanitización (filtrado y normalización de comisiones y docentes)
                         ↓
3. Resolución de materias y códigos pendientes
                         ↓
4. Consolidación de la oferta más reciente de cada materia
                         ↓
5. Resolución incremental de identidades docentes, con intervención humana
                         ↓
6. Resolución de comisiones y cálculo de altas/reactivaciones/bajas
                         ↓
7. Revisión del diff y generación de la migración SQL
                         ↓
            Aplicación separada y posterior relectura de la DB
```

### 1. Ingress

Recibir las ofertas producidas por `../parser-siu`. Su contrato actual es:

- Oferta: `carrera`, `anio`, `cuatrimestre`, `materias`.
- Materia: `codigo`, `nombre`, `comisiones`.
- Comisión: `nombre`, `docentes`, `clases`.
- Docente: `nombre`, `rol` (puede ser `null`).
- Clase: `tipo`, `dia`, `inicio`, `fin`, `aula`.

Conservar el contenido original y su procedencia: archivo, carrera y período. El parser sigue siendo genérico; las reglas particulares de FIUBA Reviews pertenecen al actualizador.

Leer también un snapshot lógico de la DB: materias, planes vigentes, equivalencias, docentes, reseñas relacionadas, cátedras, prioridades de roles y períodos de última actualización. Los cambios posteriores se calculan sobre un estado proyectado en memoria, no escribiendo durante la resolución.

La DB ya contiene `oferta_comisiones_raw` y `oferta_comisiones`, pero el JSON histórico de esta última usa `catedras` y códigos numéricos internos: **no tiene el mismo contrato ni toda la información del parser nuevo**.

TODO: Realmente las ofertas de comisiones ya guardadas en la base de datos no importan, solo son registro historico. Realmente lo importante para comparar si tengo que actualizar o no es el cuatrimestre de ultima actualizacion de la materia.

### 2. Sanitización

Estado actual: `sanitizar_ofertas_carreras` aplica `sanitizar_comisiones_materia` a cada materia y devuelve **todas** las materias, incluso las que queden sin comisiones. No se separa un catálogo observado de una oferta utilizable por ahora. Toda la información de exclusión se emite solo con logs (`tracing`), sin estructuras de motivos.

La validación de estructura, período y campos necesarios la hace `serde` al deserializar, por lo que no hay una fase aparte.

Implementado:

- Materia sin comisiones: log `warn!` («materia sin comisiones»).
- Exclusión de comisiones condicionales, clasificadas por su denominación (el nombre empieza con `CONDICIONAL`), sin eliminar las regulares. No se usa una búsqueda ingenua de `COND`: los fixtures contienen `Curso Conde Comisión 1`, que no significa condicionales.
- Exclusión de comisiones sin docentes.

Pendiente (hay tests declarados sin cuerpo en `src/oferta.rs`):

- Si hay docentes con nombre vacío, la comisión se conserva tal cual pero se emite un `warn!` indicando que está incompleta; no se la convierte silenciosamente en un equipo reducido.
- Normalizar variantes de roles (`JTP`, `Jefe/a Trabajos Practicos`, etc.). Mantener roles desconocidos para revisión; docentes sin rol van al final al ordenar.
- Una misma comisión en múltiples turnos se considera como comisiones diferentes: no se fusionan ni se pierde información de turnos. No se espera que haya otras comisiones repetidas, por lo que no hay deduplicación.

### 3. Resolver materias y códigos

Al introducir los planes nuevos se crearon materias con códigos provisionales `COD###`. La resolución debe seguir siendo incremental mientras queden pendientes.

1. Buscar primero por código oficial.
2. Si no existe, buscar candidatos provisionales por nombre normalizado y contexto del plan/carrera; consultar sus equivalencias con planes anteriores.
3. Resolver automáticamente solo los casos inequívocos según las reglas acordadas. Diferencias de nombre, múltiples candidatos o código destino ocupado requieren revisión.
4. Registrar el mapeo confirmado `código provisional → código oficial`, con su fuente. Una materia realmente nueva se distingue de un cambio de código y necesita una decisión de alta y pertenencia a planes.
5. Reescribir las referencias del estado proyectado antes de resolver docentes o cátedras.

**Actualizar el código y recuperar el historial son responsabilidades distintas.** La versión anterior también copiaba docentes, comentarios y calificaciones Dolly desde materias equivalentes del plan anterior hacia la nueva. Si se completan todos los códigos por una migración inicial, no debe perderse ese trabajo de recuperación de historial.

La copia debe ser trazable y ejecutarse una sola vez. Cuando varias equivalencias aportan al mismo docente, hay que resolver la identidad y evitar duplicados; no fusionar personas solamente porque comparten nombre. Conservar los registros del plan anterior.

El SQL de `main` actualizaba códigos por igualdad de nombre sin acentos/mayúsculas; **las equivalencias se usaban después para copiar historial**, no como condición del match inicial. Además, elegía por nombre el docente con más calificaciones Dolly entre las equivalencias: es un comportamiento anterior a revisar, no una regla obligatoria para la reescritura.

### 4. Consolidar por materia y período

Para cada materia ya identificada:

1. Comparar períodos por `(año, cuatrimestre)`, no por el ID autoincremental de la tabla.
2. Seleccionar su período más reciente entre las ofertas ingresadas.
3. Unir las comisiones de todas las carreras que ofrecen esa materia **en ese período**.
4. Deduplicar la misma comisión repetida entre carreras conservando las fuentes.
5. Si dos fuentes discrepan sobre la misma comisión, mostrar el conflicto: no concatenar docentes a ciegas ni elegir arbitrariamente el último archivo.

Las ofertas viejas pueden servir como evidencia para códigos o alias, pero no deben hacer retroceder el estado activo de una materia ya actualizada a un período posterior. Si la oferta más reciente es incompleta, no sustituirla silenciosamente por una vieja como si fuera vigente.

Una segunda ejecución puede traer otra carrera o una corrección **del mismo cuatrimestre**. Para recalcular la unión hacen falta las ofertas ya conocidas de ese período, no solamente el archivo nuevo. El período de última actualización por sí solo no basta para decidir que no hay trabajo pendiente: también importan las fuentes y su contenido.

### 5. Resolver docentes de forma incremental

La resolución se hace antes de comparar definitivamente los grupos docentes:

1. Reutilizar vínculos SIU ya confirmados en el ámbito de la materia.
2. Para nombres SIU desconocidos, proponer candidatos entre los docentes históricos de esa materia, incluyendo los recuperados por equivalencias.
3. El usuario elige **vincular con uno existente**, **crear uno nuevo** o **dejar pendiente**. La similitud de nombres es una ayuda, no autorización para fusionar.
4. Al vincular, conservar el UUID existente y sus comentarios/calificaciones. El nombre de presentación no tiene por qué reemplazarse por el texto del SIU.
5. Guardar decisiones y alias para reutilizarlos en ejecuciones futuras. No exigir resolver todos los docentes históricos: muchos todavía no aparecen en las ofertas disponibles.

Las decisiones guardadas antes de aplicar SQL son un borrador, no prueba de que la DB ya fue actualizada. Al retomar, hay que reconciliarlas con el snapshot actual. El formato de persistencia de ese borrador queda por definir.

Si hay docentes pendientes, no dar por resuelta una comisión usando solo los reconocidos. Como criterio conservador, bloquear las bajas y el cierre de actualización de esa materia hasta resolver las ambigüedades; otras materias independientes pueden continuar.

### 6. Resolver cátedras/comisiones

Comparar el conjunto completo de comisiones entrantes con las existentes, activas e inactivas:

- Misma comisión y mismo equipo: conservar su código; reactivar si corresponde.
- Comisión nueva, o comisión cuyo equipo ganó/perdió integrantes: crear una cátedra nueva.
- Cátedra anterior ausente de una oferta completa y resuelta: desactivarla, no borrarla.
- Dos comisiones con igual equipo pero distinto turno/identificador: mantenerlas separadas.

El diff debe detectar **bajas aunque no haya altas**. No alcanza con preguntar si todas las comisiones entrantes ya existen.

#### Identidad y posible hash (decisión pendiente)

Un hash solamente de nombres docentes no alcanza: colisionaría entre comisiones A/B con el mismo equipo y cambiaría al corregir o canonizar un nombre.

Recomendación: separar el **código público estable** de una **firma de comparación**, formada conceptualmente por:

```text
materia + identificador normalizado de comisión + conjunto de identidades docentes resueltas
```

El identificador debe preservar diferencias como `01A`/`01B`. Si el nombre no permite distinguir turnos, usar los datos de clases para desambiguar y solicitar revisión cuando sea necesario. Queda por decidir si un cambio de nombre u horario por sí solo constituye otra cátedra; no asumir que aula u horario son identificadores estables.

Para mostrar docentes: prioridad de rol, desempate alfabético y sin rol al final. Conviene que ese orden de presentación no determine la identidad: un ascenso de rol no cambia necesariamente el equipo. Si se elige el hash propuesto de nombres ordenados por rol, habrá que aceptar y documentar esos cambios de código, incluir materia y comisión, definir serialización inequívoca y comprobar colisiones.

La DB actual usa UUID para `catedra.codigo` y no guarda nombre de comisión, turnos ni firma. Distinguir comisiones en ejecuciones sucesivas requiere persistir esa información (cambio de esquema o estado auxiliar). Las cátedras históricas no tienen datos suficientes para reconstruir automáticamente A/B: los casos ambiguos necesitan revisión y una política de conservación de enlaces antiguos.

#### Ausencia, cobertura y oferta activa

Una ausencia solo permite una baja cuando la entrada es completa para el alcance que se está actualizando. No desactivar globalmente materias o cátedras porque no llegaron archivos de una carrera, porque falló el parser o porque se descartaron docentes vacíos.

Hay que distinguir:

- No se recibió información suficiente: conservar el último estado conocido y advertir su antigüedad.
- Se confirmó una oferta completa sin esa materia/comisión: puede marcarse como no ofertada/desactivada en ese alcance.
- La materia existe pero solo tiene condicionales: no genera cátedras regulares; las bajas requieren igualmente cobertura suficiente.

El esquema actual tiene `catedra.activa`, pero no `materia.activa` ni historial de oferta por período. Si se quiere representar explícitamente «materia no ofertada este cuatrimestre», hay que definir cómo persistirlo; no eliminarla de `plan_materia` por una ausencia cuatrimestral.

### 7. Revisar y generar SQL

Antes del SQL, mostrar un resumen con:

- Códigos corregidos y materias pendientes/desconocidas.
- Historial recuperado y posibles duplicados.
- Docentes vinculados, nuevos y pendientes.
- Comisiones mantenidas, reactivadas, creadas y desactivadas.
- Fuentes, período, exclusiones y conflictos por materia.

Generar una migración transaccional, con orden compatible con las claves foráneas, y verificaciones del estado base para no aplicarla sobre datos que cambiaron desde el snapshot. Mantener determinismo para los mismos datos y decisiones. La reaplicación debe ser inocua o rechazarse explícitamente; nunca duplicar docentes, reseñas o comisiones.

Actualizar `codigo_cuatrimestre_ultima_actualizacion` al período efectivamente reconciliado, incluso si no hubo cambios, pero no marcar como completas las materias bloqueadas o incompletas. No usar el máximo cuatrimestre global de la DB como sustituto del período de cada oferta.

**Atención al esquema real:** en la DB inspeccionada se usa `plan.vigente`, no `esta_vigente`. Tampoco todas las referencias a `materia.codigo` tienen `ON UPDATE CASCADE` (por ejemplo, `docente.codigo_materia` no lo tiene). La migración debe considerar estas restricciones, no copiar sin más el SQL de `main`.

## Investigación de códigos pendientes

En el snapshot consultado de `fiuba-reviews` hay **450 materias vinculadas a planes vigentes**, de las cuales **195 tienen códigos `COD###`**. Son cifras de diagnóstico, no constantes de la aplicación.

El cruce por nombre normalizado encontró:

- **21** candidatos en los siete fixtures de `tests/fixtures/ofertas/`.
- **24 en total** al sumar las ofertas históricas de `oferta_comisiones`.
- **171** sin coincidencia exacta en esas fuentes locales.

Estos son candidatos a verificar, no cambios ya aplicados:

| Provisional | Materia | Código observado |
| --- | --- | --- |
| COD016 | Análisis Sísmico de Estructuras | TB048 |
| COD027 | Aprovechamientos Hidráulicos | TA097 |
| COD041 | Centrales Hidráulicas | TA091 |
| COD055 | Comunicaciones Aeroespaciales A | TB156 |
| COD064 | Construcciones Hidráulicas | TA096 |
| COD084 | Dinámica de las Estructuras | TB047 |
| COD089 | Diseño Estructural | TA106 |
| COD092 | Diseño Vial Rural y Urbano | TA082 |
| COD096 | Diseño de Producto Digital | TA028 |
| COD140 | Fluidodinámica de los Hidrocarburos | TB059 |
| COD144 | Fundamentos Matemáticos de la Criptografía | CB101 |
| COD168 | Gestión Socioambiental de las Obras Civiles | TA081 |
| COD175 | Hidrología Aplicada | TB045 |
| COD176 | Hidráulica Aplicada | TA067 |
| COD177 | Hidráulica Fluvial | TA093 |
| COD178 | Hidráulica General | TB038 |
| COD203 | Ingeniería Sanitaria I | TA074 |
| COD213 | Inspección y Ejecución de Estructuras de Hormigón | TA107 |
| COD241 | Introducción a la Mecánica Computacional en Sólidos | TB134 |
| COD256 | Lenguajes y Compiladores I | TB027 |
| COD289 | Modelos Hidráulicos | TA092 |
| COD295 | Método de los Elementos Finitos | TB049 |
| COD415 | Tecnologías Emergentes en la Transición Energética | TA042 |
| COD440 | Transporte y Distribución de Hidrocarburos | TA124 |

La búsqueda web encontró fuentes oficiales parciales, **no un catálogo completo verificado que permita eliminar la fase incremental**:

- [Planes de estudio 2023–2024 de FIUBA](https://www.fi.uba.ar/institucional/plan2020/planes-de-estudio-2023): índice de resoluciones para contrastar nombres y pertenencia a planes; no es por sí mismo un mapeo completo a códigos SIU.
- [Campus de Agrimensura](https://campusgrado.fi.uba.ar/course/index.php?categoryid=66): muestra `TA195 Geodesía Física` (candidata para `COD154`) y `TC037 Teledetección Avanzada` (para `COD419`). También `TA180 Geodesia Astronómica y Mediciones`, cuyo nombre abreviado requiere contrastar con `COD153`.
- [Curso de Tecnologías de Información Geoespacial](https://campusgrado.fi.uba.ar/course/view.php?id=689&skipcheck=1): muestra `TA192`, candidato para `COD416`.

No extrapolar un código antiguo al nuevo plan sin evidencia. Los resultados aproximados pueden confundir materias distintas: «Análisis de Datos» con «Análisis de Casos», o «Lenguajes y Compiladores II» con «I». Los hallazgos parciales pueden alimentar un catálogo confirmado con fuentes, mientras el resto sigue pendiente.

## Diferencias respecto de `main`

Referencia: `actualizador-siu/servidor/` en la rama `main`.

- `ofertas.go` conserva la primera oferta por código al recorrer períodos descendentes: **no hace la unión entre carreras del mismo período** descrita como requisito de esta reescritura.
- Deduplica cátedras exclusivamente por nombres docentes, fusionando turnos que ahora deben distinguirse.
- `sincronizacion.go` y su SQL cambian códigos y copian historial en la DB durante el flujo. La nueva aplicación debe producir esos cambios como parte de una migración, no aplicarlos al planificar.
- `patches.go` detecta principalmente docentes/comisiones nuevas: no basta como modelo para bajas puras o para volver a procesar correcciones del mismo período.

## Casos mínimos para validar el flujo cuando se implemente

- Materia sin docentes que aporta un código oficial, pero no una actualización de cátedras.
- Comisión condicional excluida sin excluir las regulares; «Curso Conde» conservado.
- Mismo código y período en dos carreras: unión; períodos distintos: solo el más reciente.
- Nueva fuente del mismo período en otra ejecución: recalcular sin perder comisiones ya conocidas.
- Vinculación de un docente histórico que conserva UUID y reseñas, y reutilización del vínculo en otra ejecución.
- A/B con idénticos docentes siguen siendo dos cátedras; cambio de integrantes crea otra.
- Desaparición de una comisión sin altas: se detecta la baja.
- Oferta incompleta o docentes pendientes: no produce bajas destructivas ni cierre de actualización.
- Repetición de entradas y decisiones: no duplica historial ni genera cambios adicionales injustificados.
