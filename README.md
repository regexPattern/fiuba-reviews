# FIUBA Reviews

FIUBA Reviews es una aplicación creada para que los alumnos de la Facultad de Ingeniería de la Universidad de Buenos Aires (FIUBA) puedan subir calificaciones y comentarios sobre los docentes. El objetivo es ayudar a los alumnos a tomar decisiones informadas sobre en qué cátedras inscribirse para cada materia.

![Document](https://github.com/user-attachments/assets/f4e574de-9919-48fe-8d7d-fb8fbf0a5512)

El proyecto continúa el legado de [Dolly FIUBA](https://github.com/lugfi/dolly), la aplicación de reviews original que estuvo activa por más de 10 años, también construída por estudiantes de la facultad. Al ya no estar más activa, FIUBA Reviews mantiene viva toda la información recopilada por Dolly durante todos esos años que estuvo al servicio a los estudiantes, pero dándole un enfoque más moderno.

## Migraciones de la base de datos

El actualizador en Rust es el único responsable de las migraciones. El historial SQL vive en `actualizador-ofertas/migrations/` y se administra con **SQLx CLI**. Drizzle se usa solamente como ORM en la web: `web/src/server/db/schema.ts` debe mantenerse alineado con el SQL, pero no genera ni aplica migraciones.

Desde la raíz:

```sh
cargo install sqlx-cli --version 0.9.0 --no-default-features --features postgres,rustls
cd actualizador-ofertas
export DATABASE_URL='postgres://usuario:clave@localhost:5432/fiuba_reviews'

sqlx database create      # Solo si la base todavía no existe
sqlx migrate add descripcion_del_cambio
# Editar y revisar el nuevo archivo SQL en migrations/
sqlx migrate info
sqlx migrate run
cargo test
```

Las migraciones nuevas se escriben a mano y se versionan en este directorio. No modificar las ya aplicadas: SQLx verifica sus checksums. Los tests `#[sqlx::test]` usan automáticamente el mismo directorio y requieren una `DATABASE_URL` con permisos para crear bases de prueba e instalar las extensiones del esquema.

**Bases existentes con Drizzle:** los seis archivos históricos se conservan sin cambios, pero SQLx registra su estado en `_sqlx_migrations`, no en el historial de Drizzle. No ejecutar `sqlx migrate run` directamente sobre una base previamente migrada con Drizzle: intentaría crear objetos existentes. Antes de la transición, hacer un respaldo, verificar qué migraciones están realmente aplicadas y establecer el historial equivalente de SQLx con sus versiones y checksums en una copia de prueba. No borrar tablas ni volver a ejecutar el SQL histórico para adoptar el nuevo migrador. A partir de la transición, usar únicamente SQLx.

---

> [!WARNING]
> 🚧 REAME en construcción.
