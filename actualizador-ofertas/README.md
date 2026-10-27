# Actualizador de ofertas

Estado y próximos pasos de la reescritura del actualizador de FIUBA Reviews. Ya están conectadas la sanitización, la reconciliación automática de códigos provisionales con recuperación de historial y la consolidación por materia y período. **El próximo bloque de desarrollo es completar el paso 5: confirmar y persistir la resolución de identidades docentes**, con algunos pendientes de las etapas anteriores detallados abajo.

El objetivo final es recibir ofertas académicas parseadas, compararlas con la base de datos y generar una **migración SQL revisable**, cuya aplicación sea un paso separado. **Todavía no se genera esa migración:** la reconciliación actual ejecuta SQL sobre la conexión recibida y el llamador debe administrar la transacción. El ejecutable actual es un prototipo gráfico y todavía no ejecuta este flujo.

## Migraciones

Este proyecto administra exclusivamente el esquema compartido mediante **SQLx CLI** y los SQL de `migrations/`. Ver [instalación, comandos y transición desde Drizzle](../README.md#migraciones-de-la-base-de-datos). Ejecutar `sqlx migrate add`, `sqlx migrate info` y `sqlx migrate run` desde este directorio. Los tests `#[sqlx::test]` aplican automáticamente ese mismo historial a bases aisladas; la interfaz gráfica no aplica migraciones al iniciarse.

## Estado actual de la implementación

- `src/main.rs` abre el prototipo gráfico de `src/gui.rs` con GPUI Kit; todavía muestra una pantalla de prueba sin conexión a la DB ni carga de ofertas.
- `src/lib.rs`: `proponer_actualizaciones` sanitiza, reconcilia materias en la DB, consolida las ofertas y devuelve propuestas de docentes resueltos y pendientes, con sus candidatos. El llamador debe administrar la conexión y la transacción.
- `src/propuesta.rs` y `src/propuesta/docente.rs`: búsqueda de vínculos SIU existentes y candidatos por similitud dentro de cada materia. Las propuestas se ordenan por cantidad de pendientes y código de materia.
- `src/siu.rs`: tipos compartidos del SIU y reexportaciones. `src/siu/sanitizar.rs`, `src/siu/normalizar.rs` y `src/siu/consolidar.rs` contienen cada implementación con sus tests unitarios que pasan con `cargo test`.
- `src/catalogo.rs` y `src/catalogo/reconciliar_materias.sql`: extracción de pares únicos `(codigo, nombre)`, reemplazo de códigos provisionales y copia de docentes, calificaciones Dolly y comentarios desde equivalencias.
- Pendiente inmediato: conectar las propuestas a la interfaz y confirmar/persistir las decisiones de identidad docente. La comparación con el período de última actualización de la DB tampoco está implementada.

**Límite del ensayo actual:** `reconciliar_materias` no administra transacciones; si se la llama sin una transacción externa, sus escrituras pueden persistir. Incluso con una transacción externa, el `rollback` no revierte los avances de secuencias PostgreSQL usados al copiar calificaciones/comentarios. No equivale a una planificación de solo lectura.

## Conceptos y reglas del dominio

- **Plan de estudios:** conjunto de materias de una carrera, con su año y equivalencias. No es lo mismo que una oferta cuatrimestral.
- **Oferta:** materias y comisiones publicadas para una carrera y un cuatrimestre. Es lo que ingresa desde el parser del SIU.
- **Materia:** se identifica por su código oficial, sin distinguir carrera. Un mismo código presente en varias carreras representa la misma materia.
- **Docente:** hoy la DB registra docentes **por materia**, con un UUID, un nombre de presentación y, cuando se conoce, `nombre_siu`. Sus comentarios y calificaciones deben conservarse al reconocer su identidad.
- **Cátedra/comisión:** para este actualizador, una comisión distinguible y su grupo de docentes. Dos comisiones A/B o 1/2 son diferentes aunque tengan exactamente los mismos docentes. Agregar o quitar un docente también da lugar a una cátedra diferente para fines prácticos.
- Las comisiones de **condicionales no participan** de la actualización de cátedras.
- Para una materia se utiliza el cuatrimestre más reciente disponible con comisiones después de sanitizar. Las ofertas de ese mismo cuatrimestre se unen entre carreras; las de cuatrimestres anteriores no se mezclan con ellas.
- Desactivar no significa borrar: las materias históricas, docentes, reseñas y cátedras anteriores deben seguir existiendo.

## Flujo general

```text
1. Ingress de ofertas parseadas + lectura del estado de la DB
                         ↓
2. Sanitización (filtrado y normalización de comisiones y docentes)
                         ↓
3. Resolución de materias y códigos pendientes
                         ↓
4. Consolidación de la oferta más reciente con comisiones (se ignoran las vacías)
                         ↓
5. Resolución incremental de identidades docentes, con intervención humana
                         ↓
6. Resolución de comisiones y cálculo de altas/reactivaciones/bajas
                         ↓
7. Revisión del diff y generación de la migración SQL
                         ↓
            Aplicación separada y posterior relectura de la DB
```

### 1. Ingress — carga de fixtures implementada; snapshot pendiente

La biblioteca deserializa ofertas desde JSON, pero el prototipo gráfico todavía no permite seleccionar entradas ni conectar a la DB. Tampoco hay una lectura completa del estado de la DB ni se conserva el archivo de procedencia en el modelo.

Recibir las ofertas producidas por `../parser-siu`. Su contrato actual es:

- Oferta: `carrera`, `anio`, `cuatrimestre`, `materias`.
- Materia: `codigo`, `nombre`, `comisiones`.
- Comisión: `nombre`, `docentes`, `clases`.
- Docente: `nombre`, `rol` (puede ser `null`).
- Clase: `tipo`, `dia`, `inicio`, `fin`, `aula`. El modelo actual del actualizador no conserva `aula`.

Conservar el contenido original y su procedencia: archivo, carrera y período. El parser sigue siendo genérico; las reglas particulares de FIUBA Reviews pertenecen al actualizador.

Leer también un snapshot lógico de la DB: materias, planes vigentes, equivalencias, docentes, reseñas relacionadas, cátedras, prioridades de roles y períodos de última actualización. El diseño objetivo es calcular los cambios posteriores sobre un estado proyectado, sin escrituras persistentes durante la resolución; hoy la reconciliación del catálogo escribe sobre la conexión recibida y requiere una transacción externa.

La DB ya contiene `oferta_comisiones_raw` y `oferta_comisiones`, pero el JSON histórico de esta última usa `catedras` y códigos numéricos internos: **no tiene el mismo contrato ni toda la información del parser nuevo**.

Las ofertas de comisiones guardadas en la DB son un registro histórico, no la entrada del flujo nuevo. Para contrastar la vigencia de una oferta importa el cuatrimestre de última actualización de cada materia; el reprocesamiento de correcciones del mismo período requiere además definir la cobertura de las entradas (ver paso 4).

### 2. Sanitización — implementada

`sanitizar_ofertas_carreras` modifica las ofertas en el lugar: excluye comisiones condicionales o sin docentes, normaliza variantes de roles y conserva los turnos distintos. Mantiene docentes con nombre vacío con una advertencia y **todas las materias**, incluso las que quedan sin comisiones. Los diagnósticos se emiten mediante `tracing`.

La deserialización con `serde` comprueba la estructura y los tipos de la entrada.

### 3. Resolver materias y códigos — reconciliación automática implementada

`catalogo::extraer_materias` usa **todas las materias sanitizadas**, incluidas las que no tienen comisiones y las de períodos anteriores, y extrae los pares únicos `(codigo, nombre)`. Esta etapa no construye `OfertaMateria`: la consolidación permanece en `siu`.

`catalogo::reconciliar_materias` ejecuta el SQL de `src/catalogo/reconciliar_materias.sql`:

- Busca en `materia_vigente` códigos provisionales `XX###` (regex `^XX[0-9]{3}$`) con `usa_equivalencias`, por igualdad de nombre ignorando mayúsculas y acentos. Asume que esa vista materializada está al día.
- Reemplaza el código provisional por el del SIU y se apoya en las FKs con `ON UPDATE CASCADE` para propagarlo.
- Copia docentes desde `codigos_materias_fuente`, con UUID nuevos, junto con sus resúmenes, calificaciones Dolly y comentarios; remapea el vínculo entre comentario y calificación y conserva los registros de origen.
- Si varias equivalencias tienen docentes de igual nombre, elige una sola fuente: más calificaciones Dolly, luego más comentarios y finalmente menor código de docente. No une el historial de las fuentes descartadas ni copia `nombre_siu`.
- Devuelve códigos anterior/nuevo, equivalencias y cantidades copiadas por materia, emite diagnósticos y refresca las vistas materializadas cuando hubo reemplazos.

La condición `^XX[0-9]{3}$` evita volver a copiar el historial de una materia cuyo código ya fue reemplazado por uno oficial. No recupera por sí sola historial de materias que ya tenían código oficial antes de ejecutar esta etapa.

**Todavía pendiente:** detectar y presentar matches ambiguos, códigos destino ocupados, materias desconocidas y altas; registrar decisiones y fuentes; probar la reconciliación contra PostgreSQL. La selección por nombre del SQL actual no equivale a una resolución humana de identidad docente. Las reglas siguientes describen el alcance objetivo, no funcionalidades ya completas.

Al introducir los planes nuevos se crearon materias con códigos provisionales `COD###`, renombrados a `XX###` por la migración `migrations/0005_rename_placeholder_codes.sql` (desde este proyecto) conservando los tres dígitos. La resolución debe seguir siendo incremental mientras queden pendientes. Las referencias `COD###` del diagnóstico histórico de abajo corresponden al snapshot anterior a esa migración.

1. Buscar primero por código oficial.
2. Si no existe, buscar candidatos provisionales por nombre normalizado y contexto del plan/carrera; consultar sus equivalencias con planes anteriores.
3. Resolver automáticamente solo los casos inequívocos según las reglas acordadas. Diferencias de nombre, múltiples candidatos o código destino ocupado requieren revisión.
4. Registrar el mapeo confirmado `código provisional → código oficial`, con su fuente. Una materia realmente nueva se distingue de un cambio de código y necesita una decisión de alta y pertenencia a planes.
5. Reescribir las referencias del estado proyectado antes de resolver docentes o cátedras.

**Actualizar el código y recuperar el historial son responsabilidades distintas.** La versión anterior también copiaba docentes, comentarios y calificaciones Dolly desde materias equivalentes del plan anterior hacia la nueva. Si se completan todos los códigos por una migración inicial, no debe perderse ese trabajo de recuperación de historial.

La copia debe ser trazable y ejecutarse una sola vez. Cuando varias equivalencias aportan al mismo docente, hay que resolver la identidad y evitar duplicados; no fusionar personas solamente porque comparten nombre. Conservar los registros del plan anterior.

El SQL de `main` actualizaba códigos por igualdad de nombre sin acentos/mayúsculas; **las equivalencias se usaban después para copiar historial**, no como condición del match inicial. Además, elegía por nombre el docente con más calificaciones Dolly entre las equivalencias: es un comportamiento anterior a revisar, no una regla obligatoria para la reescritura.

### 4. Consolidar por materia y período — implementado

Implementado: `consolidar_ofertas_materias` agrupa por código, selecciona el período más reciente con comisiones comparando `(año, cuatrimestre)` y une carreras y comisiones de ese período, deduplicando comisiones exactamente iguales.

**Implementado:** dentro de la consolidación se ignoran las materias sin comisiones **antes de comparar períodos**. Si una materia tiene comisiones en 2025 pero su oferta de 2026 queda vacía tras sanitizar, se conserva la oferta de 2025 con su período original. Esto funciona independientemente del orden de entrada. Las carreras que no aportan comisiones no se incluyen en la oferta consolidada. Los metadatos de todas las materias, incluidas las vacías, ya habrán participado de la reconciliación del catálogo. Ignorar una oferta vacía no implica desactivar cátedras existentes.

También queda pendiente conservar la procedencia de cada comisión y mostrar conflictos entre fuentes, sin concatenar docentes a ciegas. Al contrastar con la DB, una oferta vieja no debe hacer retroceder una materia ya actualizada a un período posterior.

Una segunda ejecución puede traer otra carrera o una corrección **del mismo cuatrimestre**. Para recalcular la unión hacen falta las ofertas ya conocidas de ese período, no solamente el archivo nuevo. El período de última actualización por sí solo no basta para decidir que no hay trabajo pendiente: también importan las fuentes y su contenido.

### 5. Resolver docentes de forma incremental — próximo paso

La búsqueda de propuestas ya está conectada a `proponer_actualizaciones`: reconoce vínculos por `nombre_siu` normalizado y busca hasta cinco candidatos por similitud dentro de cada materia. Separa los docentes resueltos de los pendientes, sin confirmar ni escribir nuevos vínculos. Hay tests unitarios de agrupación y tests de integración con PostgreSQL que usan las migraciones de este proyecto.

El punto de partida es una oferta consolidada por materia y un catálogo reconciliado sobre la misma conexión, donde pueden leerse también los docentes recién copiados por equivalencias. La consolidación ya excluye del trabajo de cátedras las materias sin comisiones.

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
- `sincronizacion.go` y su SQL cambian códigos y copian historial en la DB durante el flujo. La reescritura ya implementa una reconciliación SQL con escrituras, ensayada dentro de una transacción que se revierte. Sigue pendiente convertir esos cambios en una migración revisable, separada de su aplicación.
- `patches.go` detecta principalmente docentes/comisiones nuevas: no basta como modelo para bajas puras o para volver a procesar correcciones del mismo período.

## Casos pendientes para validar el flujo

La sanitización y la consolidación básica tienen, junto con la normalización de roles, tests unitarios en `src/siu/sanitizar.rs`, `src/siu/normalizar.rs` y `src/siu/consolidar.rs`. Pasan con `cargo test`, pero no validan el SQL ni la integración con PostgreSQL. Quedan por cubrir la integración y las etapas pendientes:

- Reconciliación de códigos con copia de historial, conservación de origen y remapeo de comentarios/calificaciones.
- Reejecución tras confirmar una reconciliación, matches ambiguos, código destino ocupado y restricciones de claves foráneas.

- Materia sin comisiones que aporta un código oficial o historial por equivalencias, pero no una actualización de cátedras.
- Cubierto por tests unitarios: una oferta reciente vacía conserva la anterior con comisiones y su período original (en ambos órdenes de entrada).
- Nueva fuente del mismo período en otra ejecución: recalcular sin perder comisiones ya conocidas.
- Vinculación de un docente histórico que conserva UUID y reseñas, y reutilización del vínculo en otra ejecución.
- A/B con idénticos docentes siguen siendo dos cátedras; cambio de integrantes crea otra.
- Desaparición de una comisión sin altas: se detecta la baja.
- Oferta incompleta o docentes pendientes: no produce bajas destructivas ni cierre de actualización.
- Repetición de entradas y decisiones: no duplica historial ni genera cambios adicionales injustificados.
