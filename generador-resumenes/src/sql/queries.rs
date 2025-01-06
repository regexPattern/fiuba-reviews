use std::collections::HashMap;

use sqlx::{types::Uuid, PgPool};

use super::Sql;

const MIN_COMENTARIOS_ACTUALIZACION: usize = 3;
const PROPORCION_COMENTARIOS_ACTUALIZACION: usize = 2;

pub async fn comentarios_docentes(
    db: &PgPool,
) -> anyhow::Result<HashMap<Uuid, (String, Vec<String>)>> {
    let comentarios: Vec<(Uuid, String, String)> = sqlx::query_as(&format!(
        r"
SELECT d.codigo, d.nombre, c.contenido
FROM comentario c
INNER JOIN docente d
ON c.codigo_docente = d.codigo
WHERE c.codigo_docente IN (
  SELECT d.codigo
  FROM docente d
  INNER JOIN comentario c
  ON c.codigo_docente = d.codigo
  GROUP BY d.codigo
  HAVING COUNT(c) > (d.comentarios_ultimo_resumen * {})
  AND COUNT(c) > {}
);
",
        PROPORCION_COMENTARIOS_ACTUALIZACION, MIN_COMENTARIOS_ACTUALIZACION
    ))
    .fetch_all(db)
    .await?;

    tracing::info!("comentarios obtenidos de la base de datos");

    let mut comentarios_de_docente = HashMap::new();

    for (codigo_docente, nombre_docente, contenido) in comentarios {
        let comentarios_de_docente = comentarios_de_docente
            .entry(codigo_docente)
            .or_insert((nombre_docente, Vec::new()));

        comentarios_de_docente.1.push(contenido);
    }

    let cantidad_docentes = comentarios_de_docente.len();

    tracing::info!("encontrados {cantidad_docentes} docentes que requiren actualización");

    Ok(comentarios_de_docente.into_iter().take(5).collect())
}

pub async fn actualizar_resumen_docentes(
    resumenes_docentes: Vec<(Uuid, String, usize)>,
) -> Option<String> {
    let resumenes_docentes: Vec<_> = resumenes_docentes.into_iter().map(|r| r.sql()).collect();

    if !resumenes_docentes.is_empty() {
        Some(format!(
            r"
UPDATE docente AS d
SET resumen_comentarios = val.resumen_comentarios,
    comentarios_ultimo_resumen = val.comentarios_ultimo_resumen
FROM (
    VALUES
        {}
)
AS val(codigo_docente, resumen_comentarios, comentarios_ultimo_resumen)
WHERE d.codigo::text = val.codigo_docente;
",
            resumenes_docentes.join(",\n        ")
        ))
    } else {
        None
    }
}
