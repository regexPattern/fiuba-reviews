mod hugging_face;

use std::sync::Arc;

use async_scoped::TokioScope;
use reqwest::Client;
use sqlx::{types::Uuid, FromRow, PgPool};
use tokio::sync::Semaphore;

// Hugging Face tiene un limite bastante generoso de request simultaneas, pero en general he notado
// que alrededor de 20 es lo suficientemente bueno como para que no te retorne error el servidor
// por tantas requests por segundo.
const MAX_INFERENCE_API_REQUESTS_CONCURRENTES: usize = 20;

// Proporcion entre la cantidad de comentarios actuales y la cantidad de comentarios que tenia un
// docente al momento de la ultima actualizacion de su descripcion. Por ejemplo, para un valor de
// 2, se va a habilitar la regeneracion de la descripcion del docente cuando el numero de
// comentarios actuales del misma sea mas del doble del que tuvo durante la ultima actualizacion.
const FACTOR_ACTUALIZACION_DESCRIPCION: usize = 2;

// Interrumpir la generacion tras cierto numeros de intento consecutivos. Esto usualmente pasa si
// alcanzamos el maximo de request por hora de Hugging Face, con lo que no tiene sentido seguir
// gastando network haciendo mas requests.
const MAX_ERRORES_CONSECUTIVOS: usize = 10;

#[derive(FromRow)]
struct Docente {
    codigo: Uuid,
    comentarios_ultima_descripcion: i32,
}

#[derive(FromRow)]
struct Comentario {
    contenido: String,
}

pub async fn actualizar(conexion: &PgPool, api_key: String) -> anyhow::Result<Option<String>> {
    let cliente_http = Client::new();

    let docentes = Docente::obtener_de_db(&conexion).await?;
    let cantidad_docentes = docentes.len();

    let semaphore = Arc::new(Semaphore::new(MAX_INFERENCE_API_REQUESTS_CONCURRENTES));

    let (_, tasks) = TokioScope::scope_and_block(|s| {
        for docente in docentes {
            let cliente_http = Client::clone(&cliente_http);
            let conexion = PgPool::clone(conexion);
            let permiso = Arc::clone(&semaphore);
            let api_key = &api_key;

            s.spawn(async move {
                let _permiso = permiso.acquire_owned().await.unwrap();
                let codigo_docente = docente.codigo.clone();

                let resultado_descripcion =
                    docente.query_sql(cliente_http, conexion, api_key).await;

                tracing::debug!("terminado el procesamiento de docente '{codigo_docente}'");

                (codigo_docente, resultado_descripcion)
            });
        }
    });

    let mut tuplas_actualizaciones = Vec::with_capacity(cantidad_docentes);
    let mut numero_errores = 0;

    for task in tasks {
        let (codigo_docente, resultado_descripcion) = task.unwrap();
        match resultado_descripcion {
            Ok(tupla_actualizacion) => {
                if let Some(valores) = tupla_actualizacion {
                    tracing::info!("actualizada la descripcion de docente '{codigo_docente}'");
                    tuplas_actualizaciones.push(valores);
                }
                numero_errores = 0;
            }
            Err(err) => {
                tracing::error!("error generando descripcion para docente '{codigo_docente}'");
                tracing::error!("descripcion error: {err}");
                numero_errores += 1;
            }
        }

        if numero_errores == MAX_ERRORES_CONSECUTIVOS {
            tracing::info!("generacion interrumpida tras {numero_errores} errores consecutivos");
            break;
        }
    }

    if tuplas_actualizaciones.is_empty() {
        return Ok(None);
    }

    let query_sql = format!(
        r#"
UPDATE docente AS d
SET descripcion = v.descripcion,
    comentarios_ultima_descripcion = v.comentarios_ultima_descripcion
FROM (VALUES
    {}
) AS v(codigo, descripcion, comentarios_ultima_descripcion)
WHERE v.codigo = d.codigo;
"#,
        tuplas_actualizaciones.join(",")
    );

    Ok(Some(query_sql))
}

impl Docente {
    async fn query_sql(
        self,
        cliente_http: Client,
        conexion: PgPool,
        api_key: &str,
    ) -> anyhow::Result<Option<String>> {
        let comentarios = self.obtener_comentarios_de_db(conexion).await?;

        if !self.require_nueva_descripcion(comentarios.len() as i32) {
            return Ok(None);
        }

        let descripcion =
            hugging_face::generar_descripcion(cliente_http, &self.codigo, &comentarios, &api_key)
                .await?;

        let query = format!(
            r#"('{}', '{}', {})"#,
            self.codigo,
            descripcion.replace("'", "''"),
            comentarios.len()
        );

        Ok(Some(query))
    }

    async fn obtener_de_db(conexion: &PgPool) -> anyhow::Result<Vec<Self>> {
        Ok(sqlx::query_as::<_, Docente>(
            r#"
SELECT codigo, comentarios_ultima_descripcion
FROM docente
"#,
        )
        .fetch_all(conexion)
        .await?)
    }

    async fn obtener_comentarios_de_db(&self, pool: PgPool) -> anyhow::Result<Vec<String>> {
        let comentarios = sqlx::query_as::<_, Comentario>(
            r#"
SELECT contenido
FROM comentario
WHERE codigo_docente = $1"#,
        )
        .bind(&self.codigo)
        .fetch_all(&pool)
        .await?;

        let comentarios = comentarios.into_iter().map(|c| c.contenido);

        Ok(comentarios.collect())
    }

    fn require_nueva_descripcion(&self, comentarios_actuales: i32) -> bool {
        comentarios_actuales
            > self.comentarios_ultima_descripcion * FACTOR_ACTUALIZACION_DESCRIPCION as i32
    }
}
