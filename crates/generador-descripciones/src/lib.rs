use reqwest::Client;
use serde::Serialize;
use sqlx::PgPool;

const HF_INTERENCE_API: &str =
    "https://api-inference.huggingface.co/models/facebook/bart-large-cnn";

struct Docente {
    codigo: String,
    comentarios_ultima_descripcion: i32,
}

struct Comentario {
    contenido: String,
}

pub async fn query_actualizacion_db(pool: &PgPool) -> anyhow::Result<String> {
    let http = Client::new();

    let docentes = sqlx::query_as!(
        Docente,
        r#"
SELECT codigo, comentarios_ultima_descripcion
FROM Docente 
"#
    )
    .fetch_all(pool)
    .await?;

    let mut tasks = Vec::with_capacity(docentes.len());
    let mut value_queries = Vec::with_capacity(docentes.len());

    for docente in docentes {
        let pool = PgPool::clone(pool);
        let http = Client::clone(&http);

        tasks.push(tokio::spawn((move || async move {
            docente.query_actualizacion_values(http, pool).await
        })()));
    }

    for task in tasks {
        if let Ok(query) = task.await? {
            if let Some(query) = query {
                value_queries.push(query);
            }
        } else {
            tracing::error!("error ejecutando task asincrona");
        }
    }

    let values = value_queries.join(",");

    Ok(format!(
        r#"
UPDATE Docente as d
SET descripcion = a.descripcion,
    comentarios_ultima_descripcion = a.comentarios_ultima_descripcion
FROM (VALUES
    {values}
) as a(codigo, descripcion, comentarios_ultima_descripcion)
WHERE a.codigo = d.codigo;
"#
    ))
}

impl Docente {
    async fn query_actualizacion_values(
        &self,
        http: Client,
        pool: PgPool,
    ) -> anyhow::Result<Option<String>> {
        let comentarios = self.obtener_comentarios(pool).await?;

        if comentarios.is_empty()
            || (comentarios.len() as i32) < self.comentarios_ultima_descripcion * 2
        {
            return Ok(None);
        }

        let descripcion = Self::generar_descripcion(http, &comentarios).await?;

        tracing::info!("actualizada descripcion de {}", self.codigo);

        let query = format!(
            r#"
('{}', '{}', {})
"#,
            self.codigo,
            descripcion.replace("'", "''"),
            comentarios.len()
        );

        Ok(Some(query))
    }

    async fn obtener_comentarios(&self, pool: PgPool) -> anyhow::Result<Vec<String>> {
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

    async fn generar_descripcion(http: Client, comentarios: &[String]) -> anyhow::Result<String> {
        #[derive(Serialize)]
        struct Payload {
            inputs: String,
        }

        // http.post(HF_INTERENCE_API)
        //     .bearer_auth("")
        //     .json(&Payload {
        //         inputs: comentarios.join("."),
        //     })
        //     .send()
        //     .await?;

        Ok("ESTA ES MI NUEVA DESCRIPCION".into())
    }
}
