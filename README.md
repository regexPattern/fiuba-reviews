# FIUBA Reviews

**FIUBA Reviews** es una aplicación creada para que los alumnos de la Facultad de Ingeniería de la Universidad de Buenos Aires (FIUBA) puedan subir calificaciones y comentarios sobre los docentes. El objetivo es ayudar a los alumnos a tomar decisiones informadas sobre en qué cátedras inscribirse para cada materia.

## Motivación

El proyecto surgió como una propuesta para realizar mejoras en [Dolly FIUBA](https://github.com/lugfi/dolly). Como alumno de la facultad, solía utilizar Dolly cada cuatrimestre para obtener una idea general de qué cátedras eran más recomendadas, porque no siempre contaba con un amigo que hubiese cursado una materia específica y pudiera darme una recomendación. Finalmente se terminó convirtiendo en una aplicación autónoma, en vez de un simple rediseño del frontend de la aplicación original.

**FIUBA Reviews** se construye utilizando los datos que Dolly FIUBA recolectó durante casi una década, ordenándolos y dándoles una estructura más organizada. Ahora que Dolly dejó de funcionar, pretendo que este proyecto sea el nuevo espacio donde los alumnos puedan compartir y leer opiniones sobre los docentes libremente, un grupo al que también pertenezco como usuario de mi propia aplicación.

## Desarrollo

Por ahora la aplicación esta en un estado de cambio constante, pero en líneas generales la estructura es la siguiente:

- **`app`**: Lo principal es la aplicación web en sí, construída en Svelte con SvelteKit, donde se agrupan el frontend y el backend en un mismo proyecto. Utiliza TailwindCSS para los estilos.
- **`resumidor-comentarios`**: Es una utilidad que me permite generar o actualizar los resumenes de comentarios de los docentes on-demand y de manera incremental. Utiliza IA para generar estos resúmenes.

Además, la aplicación está hosteada en Vercel y la base de datos en Supabase.

## Contribución

En los próximos meses, mi intención es pulir el repositorio, limpiar el código, modularizarlo mejor y facilitar el inicio de un entorno de desarrollo tras clonar el repositorio. Esto va a permitir que aquellos que deseen contribuir con ideas de mejora puedan hacerlo más fácilmente. Por ahora, lo más recomendado es crear un issue con tus sugerencias, ya que planeo realizar varios cambios para hacer el proyecto más mantenible a largo plazo. Para más información, consulta [este issue](https://github.com/regexPattern/fiuba-reviews/issues/23).
