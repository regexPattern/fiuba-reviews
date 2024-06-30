# Che FIUBA

Aplicación web para leer y publicar opiniones de los docentes de FIUBA, agregadas por los mismos estudiantes, para que así tengas una mejor idea de que profesores te podrían gustar más y que cátedra elegir el cuatrimestre que viene. Reimplementación de [Dolly FIUBA](https://dollyfiuba.com) con adaptación de datos existentes.

Esta aplicación no pretende ser un reemplazo a la aplicación original, sino que más bien una propuesta de posibles cambios que se pueden hacer para, desde mi perspectiva, mejorar la experiencia de usuario de la aplicación original.

## Utilización

Podés acceder a la versión live de la aplicación desde https://che-fiuba.com.

Si por el contrario querés correr un build local, tomá en cuenta que al momento de compilar la página, debe existir la variable de entorno `DATABASE_URL`. Podés usar el archivo [`compose.yaml`](https://github.com/regexPattern/che-fiuba/blob/main/compose.yaml) para arrancar la base de datos usando Docker Compose:

```bash
git clone https://github.com/regexPattern/che-fiuba

cd che-fiuba
npm install

docker compose up -d

export DATABASE_URL=postgres://postgres:postgres@localhost:5432

npm run build
npm run preview
```

Tené en cuenta que si usas el Docker Compose, durante el primer levantamiento del contenedor de la base de datos los tiempos de espera pueden prolongarse, ya que primero se tienen que insertar los datos de la aplicación original en la nueva base de datos, proceso detallado en la sección "[Adapación de los datos originales](#adaptación-de-los-datos-originales)". Lo mismo aplica si se elimina el volumen del contenedor.

## Desarrollo

La aplicación está escrita en [SvelteKit](https://kit.svelte.dev/). Entre otras especificaciones técnicas podría destacar:

- [TailwindCSS](https://tailwindcss.com/) para facilitar el estilizado.
- [DrizzleORM](https://orm.drizzle.team/) para hacer las queries a la base de datos sin tener un ORM tan abstracto.
- [shadcn-svelte](https://www.shadcn-svelte.com/) como libreria de componentes comunes para que quede más bonito.

Si querés desarrollar la aplicación en tu propia máquina corré el siguiente comando:

```bash
docker compose up -d

export DATABASE_URL=postgres://postgres:postgres@localhost:5432

npm run dev
```

Además de la aplicación web central, se desarrollaron dos herramientas que facilitan la automatización de la adaptación de los datos originales y la generación de las descripciones de los docentes.

### Adaptación de los datos originales

Para aprovechar todos los datos que Dolly ha recopilado durante años, se adaptaron los datos de la aplicación original en vez de iniciar desde cero.

Para esto la herramienta [`adaptador-datos`](https://github.com/regexPattern/che-fiuba/tree/main/crates/adaptador-datos) hace scraping de los datos de la aplicación original, y genera un archivo SQL que se carga a la base de datos de manera automática cuando se construye por primera vez la imagen de Docker de la misma.

También cuenta con funcionalidad de actualización, para poder sincronizar la base de datos existente e incorporar los nuevos datos que se fueron agregaron a Dolly desde que esta se inicializó.

### Resumen de comentarios con inteligencia artificial

Se generaron resúmenes de los comentarios de los docentes con la ayuda modelos de inteligencia artificial, para que quien use la aplicación pueda darse una idea general de qué opinan los demás estudiantes sobre un docente que no conoce. La versión de la aplicación que está activa en línea utiliza [GPT-3.5 Turbo](https://platform.openai.com/docs/models/gpt-3-5-turbo) de OpenAI.

La utilidad [`resumidor-comentarios`](https://github.com/regexPattern/che-fiuba/tree/main/crates/resumidor-comentarios) facilita la generación de los resúmenes y actualización de la base de datos. Esta herramienta debe ser utilizada de manera manual, y puede ser adaptada para soportar varios modelos de inteligencia artificial.
