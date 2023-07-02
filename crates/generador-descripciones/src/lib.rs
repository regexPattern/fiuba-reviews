#![allow(unused)]
mod hugging_face;

use std::sync::Arc;

use anyhow::Context;
use async_scoped::TokioScope;
use futures::{
    future::{self, Future},
    stream::{FuturesUnordered, StreamExt},
};
use reqwest::Client;
use sqlx::PgPool;
use tokio::{sync::Semaphore, task::JoinSet};

const MAX_INFERENCE_API_REQUESTS_CONCURRENTES: usize = 10;

// Proporcion entre la cantidad de comentarios actuales y la cantidad de comentarios que tenia un
// docente al momento de la ultima actualizacion de su descripcion. Por ejemplo, para un valor de
// 2, se va a habilitar la regeneracion de la descripcion del docente cuando el numero de
// comentarios actuales del misma sea mas del doble del que tuvo durante la ultima actualizacion.
//
const FACTOR_ACTUALIZACION_DESCRIPCION: usize = 2;

struct Docente {
    codigo: String,
    comentarios_ultima_descripcion: i32,
}

struct Comentario {
    contenido: String,
}

pub async fn actualizar(conexion: &PgPool, api_key: String) -> anyhow::Result<Option<String>> {
    let cliente_http = Client::new();

    let docentes = Docente::obtener_de_db(&conexion).await?;
    let semaphore = Arc::new(Semaphore::new(MAX_INFERENCE_API_REQUESTS_CONCURRENTES));

    let (_, tasks) = TokioScope::scope_and_block(|s| {
        for docente in docentes {
            let cliente_http = Client::clone(&cliente_http);
            let conexion = PgPool::clone(conexion);
            let permiso = Arc::clone(&semaphore);
            let api_key = &api_key;

            s.spawn(async move {
                let _permiso = permiso.acquire_owned().await.unwrap();
                let value = docente.query_sql(cliente_http, conexion, api_key).await;
                value
            });
        }
    });

    let mut updated_values = Vec::with_capacity(10);
    for task in tasks {
        if let Some(update_tuple) = task.unwrap()? {
            updated_values.push(update_tuple);
        }
    }

    if updated_values.is_empty() {
        return Ok(None);
    }

    let query_sql = format!(
        r#"
UPDATE Docente as d
SET descripcion = a.descripcion,
    comentarios_ultima_descripcion = a.comentarios_ultima_descripcion
FROM (VALUES
    {}
) as a(codigo, descripcion, comentarios_ultima_descripcion)
WHERE a.codigo = d.codigo;
"#,
        updated_values.join(",")
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
        let comentarios = self
            .obtener_comentarios_de_db(conexion)
            .await
            .map_err(|err| {
                tracing::error!(
                    "error obteniendo comentarios de la base de datos para '{}'",
                    self.codigo
                );
                err
            })?;

        if !self.require_nueva_descripcion(comentarios.len() as i32) {
            return Ok(None);
        }

        let descripcion =
            hugging_face::generar_descripcion(cliente_http, &self.codigo, &comentarios, &api_key)
                .await
                .map_err(|err| {
                    tracing::error!(
                        "error generando la descripcion para docente '{}'",
                        self.codigo,
                    );
                    err
                })
                .unwrap();

        tracing::info!("actualizada descripcion de '{}'", self.codigo);

        let query = format!(
            r#"('{}', '{}', {})"#,
            self.codigo,
            descripcion.replace("'", "''"),
            comentarios.len()
        );

        Ok(Some(query))
    }

    async fn obtener_de_db(conexion: &PgPool) -> anyhow::Result<Vec<Self>> {
        Ok(sqlx::query_as!(
            Docente,
            r#"
SELECT codigo, comentarios_ultima_descripcion
FROM Docente 
"#
        )
        .fetch_all(conexion)
        .await?)
    }

    async fn obtener_comentarios_de_db(&self, pool: PgPool) -> anyhow::Result<Vec<String>> {
        let comentarios = sqlx::query_as!(
            Comentario,
            r#"
SELECT contenido
FROM Comentario
WHERE codigo_docente = $1"#,
            self.codigo
        )
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
