# FIUBA Reviews

Reimplementación de https://dollyfiuba.com con adaptación de los datos existentes.

El proyecto consiste de dos partes principales: el cliente web y la base de datos (que incluye el adaptador de datos y el generador de descripciones).

Para correrlo, se debe correr primero la base de datos (para automatizar esto está el archivo [`compose.yaml`](https://github.com/regexPattern/fiuba-reviews/blob/main/compose.yaml)), para que luego la aplicación web pueda renderizarse con los datos de la misma. A continuación se detalla cómo iniciar ambos componentes, y más información de los mismos:

## Cliente web

Es una aplicación web escrita en SvelteKit. Entre otras especificaciones técnicas podría destacar:
- [TailwindCSS](https://tailwindcss.com/) para facilitar el estilizado.
- [DrizzleORM](https://orm.drizzle.team/) para hacer las queries a la base de datos sin tener un ORM tan abstracto.
- [shadcn-svelte](https://www.shadcn-svelte.com/) como libreria de componentes comunes para que quede más bonito.

La aplicación requiere que la variable de entorno `DATABASE_URL` esté configurada al correr el servidor.

## Base de datos

La aplicación utiliza Postgres como motor de base de datos, sin embargo lo más relevante desde el aspecto técnico es cómo se popula la base de datos con los datos de la aplicación de Dolly original, y cómo se logró hacer este proceso reproducible. Personalmente me resulta la parte más emocionante de la aplicación, tanto por su proceso de implementación como por el resultado.

### Adaptación de los datos originales

Este proyecto inició con la idea de escribir un nuevo cliente frontend para Dolly, sin un nuevo backend, sin una nueva base de datos con un esquema diferente ni nada por el estilo, por lo que el verdadero objetivo de esta aplicación siempre fue utilizar los datos de la aplicación original. He utilizado Dolly desde que ingresé a la facultad, y ya en ese entonces la aplicación tenía algunos años funcionando, habiendo recopilado una grandísima cantidad de comentarios, de una igualmente grande cantidad de docentes, cátedras y materias.

Sin embargo, cuando inicié a plantearme ideas de cómo crear un nuevo cliente, terminé concluyendo que lo más óptimo iba a ser reestructurar todo por dentro, para que finalmente el consumo de los datos desde el nuevo cliente web me resultara más conveniente. Para esto, tuve que encontrar una forma de utilizar los datos ya existentes, pero adaptarlos a una estructura más cómoda para mí. 

Acá es cuando nace [`adaptador-datos`](https://github.com/regexPattern/fiuba-reviews/tree/main/crates/adaptador-datos). Un programa escrito en Rust que hace scraping de los datos de la aplicación original, y termina generando un archivo SQL que se carga a la base de datos de manera automática cuando se construye por primera vez la imagen de Docker de la misma.

### Generación de descripciones con inteligencia artificial

Durante el desarrollo inicial del proyecto, surgió la idea de agregar una nueva feature a la aplicación, que aportara algo significativo a los alumnos cuando la usaran. Ya que el objetivo principal de la misma es que los alumnos lean los comentarios los docentes de las cátedras de las materias que quieren ver (probablemente las que quieran inscribir en el cuatrimestre entrante, o al menos es mi prinicipal caso de uso), la mejor forma que se me ocurrió para facilitar dicha tarea, era hacer un resumen de los comentarios de los docentes, algo que sirve para darse una idea de la opinión general de cada profesor, especialmente útil en el caso de los profesores que tienen muchos comentarios.

Es así que aprovechando la sustancial mejora de los modelos de sumarización en los últimos tiempos y la disponibilidad de un servicio como [Inference API](https://huggingface.co/inference-api) que provee una forma simple de acceder a dichos modelos de manera gratuita, surge la segunda utilidad de esta aplicación, [`generador-descripciones`](https://github.com/regexPattern/fiuba-reviews/tree/main/crates/generador-descripciones). Otro programa, también escrito en Rust, que se conecta al servidor de la base de datos ya corriendo, y uno por uno, identifica los docentes tienen los comentarios suficientes para que el modelo de sumarización seleccionado haga un trabajo decente, manda a generar dichas descripciones, y finalmente actualiza los docentes cuya descripción pudo generar en la base de datos conectada.

Hay algunos aspectos del funcionamiento de esta utilidad que vale la pena estacar:
- Se require que las variables de entorno `DATABASE_URL` e `INFERENCE_API_KEY` estén configuradas (es posible definirlas en un archivo `.env` si se desea). Respecto a la segunda variable, esta es la API key de Inference API, que permite acceder a los modelos de Hugging Face disponibles. Podés encontrar más información sobre como crear una llave personal [acá](https://huggingface.co/docs/api-inference/quicktour).
- Al momento de elaboración de este proyecto, el tier gratuito de Inference API tiene un límite de requests consecutivas, por lo que es posible que se tenga que correr la utilidad múltiples veces para terminar de actualizar la base de datos con los resúmenes de todos los docentes.

## Moraleja y detalles aún más técnicos

Si bien la gracia del proyecto es tener una aplicación funcional que pueda ser de utilidad a un estudiante de la facultad, o bien servir como sugerencia para cambios que me parece que podrían venir bien a la aplicación original de Dolly, la más significativo de este proyecto para mí ha sido el aprendizaje técnico que he obtenido durante la elaboración del mismo. Gracias a este proyecto aprendí a usar Docker, tecnología de la que sabía poco más que sobre su concepto, pero que ahora me resulta una de las herramientas más utiles con las que me puedo valer.

Los extensos tiempos de compilación de Rust me llevaron a aprender cómo usar cache para la construcción de las imágenes y así evitar recompilar los binarios cuando quisiera construir una nueva imagen de la base de datos tras la primera compilación de los mismos. El hecho de que con solo correr el comando `docker compose up` se compile la utilidad de adaptación de los datos iniciales, se corra la misma, y luego se genere un archivo SQL que automáticamente va a popular la base de datos antes de que se inicie la misma me sigue resultado fascinante aún después de haberlo logrado implementar.

Por otra parte, la utilidad de generación de resúmenes me permitió por primera vez darle un uso práctico a la inteligencia artificial en uno de mis proyectos. Este mismo proyecto también me hizo experimentar un poco más con programación asíncrona en Rust, pesadilla de la que aún tengo mucho por aprender, pero que me llevó a escribir el [código](https://github.com/regexPattern/fiuba-reviews/blob/c7f219e50ad843b30cf7e9c0fb06f5f5b3379321/crates/generador-descripciones/src/lib.rs#L65) más hermoso que considero que he escrito, aún 6 meses (al momento que escrito este README) después de haberlo escrito.

Finalmente, la aplicación web, que creí que iba a ser el único foco de la aplicación en un inicio, quizá terminó siendo lo menos interesante técnicamente, pero lo más estresante definitivamente, ya que a veces se me iba el perfeccionismo  de la mano, y terminaba reescribiendo todo el frontend desde cero solo porque no me gustaban el layout que le estaba dando. Esta parte del proceso me ayudó a reincorporarme un poco a las prácticas modernas del desarrollo web, rama de la programación que había dejado hace algún tiempo, y aprender las bases de Svelte, un framework del que ahora me puedo valer para mis futuros proyectos web.

A pesar de todos los detalles técnicos con los que podría extender mucho mas este README, lo que más me gustó del proceso de desarrollo de esta aplicación, fue que surgió de una necesidad que yo como estudiante sentí que podía solucionar. Es una aplicación que me va a servir por lo menos a mí, y espero que también a la comunidad también, directa o indirectamente.
