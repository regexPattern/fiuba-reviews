# Idea General

Estoy construyendo un actualizador de catedras (comisiones) para mi aplicacion de reviews de los docentes de mi facultad. Para esto, estoy obteniendo los listados de las ofertas de comisiones de todas las materias de todas las carreras de la facultad (solo las ofertas que esten publicadas para cada cuatrimestre del que cuento con los planes, hay materias que a veces no ofrecen catedra en un cuatrimestre, por lo tanto, podria no tener informacion de todas).

Una vez se tienen cargadas todas las ofertas de comisiones de todas las materias disponibles, me quedo solamente con la version mas reciente de cada materia, es decir, si para una materia hay ofertas del primer cuatrimestre de 2025 y tambien del segundo de 2024, voy a preferir quedarme con la mas reciente. Esto se hace porque las ofertas son retornadas en orden cronologico descendente, dado el codigo de su cuatrimestre, que tambien sigue esta convencion en la base de datos.

# Sincronizacion de codigos y migracion desde equivalencias

Una vez hecho esto se tienen que sincronizar las materias de la base de datos que todavia no estan asociadas con su equivalente del siu. Esto sucede porque cuando se populo la base de datos con las materias de los nuevos planes, se desconocian los codigos de las mismas. Cabe destacar que solo tengo interes en trabajar con las materias de los nuevos planes de ahora en adelante. Una materia del nuevo plan que esta desincronizada con una del siu tiene un codigo con un placeholder COD, por lo que codigo_siu != codigo_db, y tampoco tiene ni docentes, ni comentarios, ni calificaciones, ni catedras asociadas. Sin embargo, si puedo asegurar que los nombres normalizados coinciden para dos materias que son iguales, aunque no esten sincronizadas aun, es por esto que la query de sincronizacion que utilizo hace este match por nombre normalizado.

Cuando se sincroniza una materia, ademas de corregir el codigo y poner el codigo correcto (el traido de las ofertas de comisiones del siu), se copian todos los docentes, comentarios y calificaciones de las equivalencias de esa materia en los planes que no estan vigentes. Se puede asumir que la combinacion (nombre_docente, codigo_materia) es unica para las materias del plan anterior, lo que significa que, si cuando se migran estos datos de una materia que tiene multiples equivalencias a una sola, por ejemplo, si la materia CB001 del nueva plan tiene como equivalencia a la materia M1, M2, M3, y estas tres materias, si bien son tecnicamente diferentes en la base de datos, realmente contienen los mismos datos duplicados, cuando se haga la migracion esta condicion de unicidad se rompe, por lo que, para evitar conflictos, nos terminamos quedando con la copia del docente que mas informacion (calificaciones) tenga asociada.

# Preparacion para actualizacion

Luego de la sincronizacion viene el paso de preparacion para planificar las actualizacion que se van a hacer en la base de datos. Recordemos que el objetivo final del actualizador es justamente actualizar las comisiones disponibles para una materia (catedras). Actualmente estoy en el proceso de diseño de esta funcionalidad, quiero que me ayudes a implementarla. 

Las materias a actualizar que traemos de los planes del siu incluyen todas las materias que no han sido actualizadas al ultimo cuatrimestre. Pero realmente podrian ya estar totalmente al dia en cuanto a su oferta de comisiones y planilla de docentes, que son las cosas que a mi me interesa actualizar en mi app.

La oferta de comisiones de cada materia proveniente del siu posee tanto un listado de comisiones/catedras como de los docentes que las componen. Actualizar la oferta de comisiones en la base de datos implica tener armadas las catedras que estan presentes en esta ultima version de la oferta de comisiones que se tiene para cada materia disponible. Como una catedra no es mucho mas que una agrupacion de docentes, esto implica que deben estar registrados los docentes listados en el siu en la base de datos. Contamos con el listado de docentes que necesitamos encontrar en la base de datos para armar las catedras de una materia, son justamente los que traemos del siu.

Es decir, una materia ya esta actualizada si todos los docentes encontrados en el siu para la misma ya estan registrados en la base datos, y si todas las catedras que tiene en el siu ya estan registradas tambien.

Como digo, Una catedra no es mas que una agrupacion de docentes. Si bien en la base de datos estan identificadas con un codigo, lo verdaderamente representativo de las catedras es su nombre, que se puede obtener concatenando los nombres de los docentes que la componen. Asi, por ejemplo, la catedra de los docentes "Carlos Castillo", "Irene Cardona" y "Lionel Messi", seria la catedra "Cardona-Castillo-Messi", ordenada en orden alfabetico.

IMPORTANTE: La resolucion de las nuevas catedras encontradas en las ofertas del siu se va a hacer manualmente en la app web que voy a construir despues. La idea es simplemente detectar cuales son las catedras que se deben crear.

## Matcheo de docentes

Para cada docente de esta lista se pueden presentar 3 situaciones:
1. Hay un docente en la base de datos cuyo `nombre_siu` coincide con el nombre del docente encontrado en el siu. En este caso podemos considerar al docente como resuelto.
2. No hay un docente en la base de datos cuyo `nombre_siu` coincida, pero hay varios docentes cuyo campo `nombre` coincide con el nombre del docente encontrado en el siu. Estos otros docentes no deben tener un `nombre_siu` asignado, ya que esto significa que ya estan "asociados" con otro docente del siu. Para esto se utilizara fuzzy matching.
3. No hay ningun docente que matchee en nombre con el docente del siu. En este caso significa que estamos ante un docente totalmente nuevo que hay que registrar.

IMPORTANTE: La resolucion de los docentes que tienen multiples matches, al igual que la resolucion de catedras, se va a hacer de manera manual en la app web que voy a construir despues. La idea es simplemente saber cuales docentes ya tienen un match perfecto, y dejar los nombres de los que no lo tienen todavia.

## Persistencia de catedras

De un cuatrimestre a otro, las catedras pueden continuar existiendo. Es decir, si la catedra "Cardona-Castillo-Messi" ya existe, no hay necesidad de eliminar este registro y crear uno nuevo con los mismos docentes. Si se diera el caso en el que se agregara un docente adicional, para simplificar, ahi si deberiamos crear una nueva catedra. Lo mismo para si se retira un docente.

A medida vayan pasando los cuatrimestres de actualizacion, van a ir quedando catedras registradas en la base de datos que no existan. Para esto tenemos dos opciones: borramos siempre todas las catedras que no se mantienen en cada cuatrimestre, o usamos algun tipo de flag para saber si estan activas o son de cuatrimestres anteriores. Para elegir alguna de estas dos opciones, hay que tener en consideracion que pasa cuando una materia no tiene actualizaciones para este cuatrimestre.
