# FIUBA Reviews

Aplicación web para leer y publicar opiniones de los docentes de FIUBA, agregadas por los mismos estudiantes, para que así tengas una mejor idea de que profesores te podrían gustar más y que cátedra elegir el cuatrimestre que viene. Reimplementación de [Dolly FIUBA](https://dollyfiuba.com) con adaptación de datos existentes.

Esta aplicación no pretende ser un reemplazo a la aplicación original, sino que más bien una propuesta de posibles cambios que se pueden hacer para, desde mi perspectiva, mejorar la experiencia de esta muy útil aplicación creada por la misma comunidad de estudiantes.

## Utilización

La aplicación se divide en dos partes: el cliente web y la base de datos. La primera depende de que la segunda esté corriendo para poder compilarse.

Si querés correr la aplicación compilada para producción, podes simplemente utilizar el [`compose.yaml`](https://github.com/regexPattern/fiuba-reviews/blob/main/compose.yaml) para levantar todo el proyecto. Para esto tenés que tener [Docker](https://www.docker.com/) instalado y correr el siguiente comando:

```bash
docker compose up
```

Los contenedores están configurados de tal manera que la aplicación web inicie a compilarse hasta que el servidor de la base de datos esté activo.

Tené en cuenta que durante el primer levantamiento del contenedor de la base de datos los tiempos de espera pueden prolongarse, ya que primero se tienen que insertar los datos de la aplicación original en la nueva base de datos, proceso detallado en la sección "[Adapación de los datos originales](https://github.com/regexPattern/fiuba-reviews/tree/main#adaptación-de-los-datos-originales)". Lo mismo aplica si se elimina el volumen del contenedor.

Por otra parte, la aplicación web se va a tener que recompilar cada vez que se inicie el contenedor, esto debido a que Docker Compose no tiene una forma de esperar a que uno de los servicios de red esté corriendo o pase su healthcheck para iniciar el build de la imagen de otro servicio.

## Desarrollo

La aplicación está escrita en [SvelteKit](https://kit.svelte.dev/). Entre otras especificaciones técnicas podría destacar:

- [TailwindCSS](https://tailwindcss.com/) para facilitar el estilizado.
- [DrizzleORM](https://orm.drizzle.team/) para hacer las queries a la base de datos sin tener un ORM tan abstracto.
- [shadcn-svelte](https://www.shadcn-svelte.com/) como libreria de componentes comunes para que quede más bonito.

Si querés desarrollar la aplicación en tu propia máquina en vez de iniciar el servidor de producción del Compose, tenés que configurar la variable de entorno `DATABASE_URL` antes, ya que es requerida para la compilación de la misma, incluídos la compilación en modo desarrollo.

Primero instalá las dependencias con el siguiente comando:

```bash
npm install
```

Podés iniciar únicamente el servicio de la base de datos del Compose, y configurar la variable de entorno que la aplicación web va a utilizar para que se conecte a ese servidor (el Compose está configurado para que haga binding de los puertos de tu máquina local). Finalmente, inicia servidor de la aplicación en modo de desarrollo. Para todo esto utilizá los siguientes comandos:

```bash
docker compose start base-de-datos
export DATABASE_URL=postgres://postgres:postgres@localhost:5432
npm run dev
```

## Utilidades

Desde el punto de vista técnico, de lo que más disfruté al construir esta aplicación fue de la automatización de la adaptación de los datos originales y la generación de las descripciones de los docentes. Para esto se desarrollaron dos utilidades escritas en Rust que automatizan estas tareas.

### Adaptación de los datos originales

Dolly ha sido una aplicación utilizada por muchos estudiantes de la facultad (me incluyo) durante muchos años. En este tiempo ha recopilado una grandísima cantidad de comentarios, de una igualmente grande cantidad de docentes, cátedras y materias. Por lo tanto, para aprovechar todo este trabajo se adaptaron los datos de la aplicación original en vez de iniciar desde cero.

Para esto se desarrolló la herramienta [`adaptador-datos`](https://github.com/regexPattern/fiuba-reviews/tree/main/crates/adaptador-datos), un programa que hace scraping de los datos de la aplicación original, y termina generando un archivo SQL que se carga a la base de datos de manera automática cuando se construye por primera vez la imagen de Docker de la misma.

### Generación de descripciones con inteligencia artificial

Mi principal caso de uso para esta aplicación es formarme una idea de cómo son los profesores de las diferentes cátedras de las materias que estoy a punto de inscribir en un cuatrimestre entrante. Para esto uno lee los comentarios de los demás estudiantes que cursaron anteriormente con el profesor en cuestión.

Con la ayuda del modelo de sumarización open source [BART](https://huggingface.co/facebook/bart-large-cnn), se generaron resúmenes de los comentarios de los docentes, para que quien use la aplicación pueda darse una idea general de qué opinan los demás estudiantes sobre un docente que no conoce.

Para facilitar la tarea de generar dichas descripciones se desarrolló una segunda utilidad, [`generador-descripciones`](https://github.com/regexPattern/fiuba-reviews/tree/main/crates/generador-descripciones), que utiliza el modelo de sumarización anteriormente mencionado a través de [Inference API](https://huggingface.co/inference-api).

Cuando se inicia una nueva base de datos utilizando el adaptor, ningún docente cuenta su descripción generada a partir del resumen de todos los comentarios asociados al mismo, ya que estos datos no están en la aplicación original de Dolly cuando se descargan los datos, ni se pueden generar automáticamente al momento de crear el script SQL con el que se inicia la base de datos ya que Inference API tiene un límite de solicitudes por hora, por lo que esta segunda utilidad tiene que ser corrida manualmente cada cierto tiempo para incrementalmente ir actualizando los registros de los docentes con las descripciones que se pudieron generar al momento de correr la utilidad.

Para correr la aplicación vas a necesitar [generar una llave para Inference API](https://huggingface.co/docs/api-inference/quicktour) y configurar las variables de entorno `DATABASE_URL` e `INFERENCE_API_KEY` al momento de ejecutar el programa. Para hacer todo eso podés correr los siguientes comandos (necesitas tener [Rust](https://www.rust-lang.org/) instalado):

```bash
export DATABASE_URL=...
export INFERENCE_API_KEY=...
cargo run --release
```
