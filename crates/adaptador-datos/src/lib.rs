mod catedra;
mod comentario;
mod docente;
mod materia;

use http_cache::{CACacheManager, CacheMode, HttpCache};
use http_cache_reqwest::Cache;
use reqwest::Client;
use reqwest_middleware::{ClientBuilder, ClientWithMiddleware};
use std::{
    collections::{HashMap, HashSet},
    sync::Arc,
};
use tokio::task::JoinHandle;
use uuid::Uuid;

use materia::Materia;

#[derive(Debug)]
struct TaskMateria {
    sql: SqlElementosMateria,
    codigos_docentes: HashMap<(i16, String), Uuid>,
}

#[derive(Default, Debug)]
struct SqlElementosMateria {
    catedras: Vec<String>,
    docentes: Vec<String>,
    calificaciones: Vec<String>,
    rel_catedras_docentes: Vec<String>,
}

pub async fn init_query() -> anyhow::Result<String> {
    let cliente_http = Arc::new(crear_cliente_http());

    let materias = Materia::descargar_todas(&cliente_http).await?;
    let comentarios = comentario::descargar_todos(&cliente_http).await?;

    let mut sql_materias = Vec::with_capacity(materias.len());

    let mut handles: Vec<JoinHandle<anyhow::Result<TaskMateria>>> =
        Vec::with_capacity(materias.len());

    for mat in materias {
        sql_materias.push(mat.sql());

        let cliente_http = Arc::clone(&cliente_http);

        handles.push(tokio::spawn(async move {
            let catedras = mat
                .descargar_catedras(&cliente_http)
                .await
                .inspect_err(|err| {
                    tracing::error!("error descargando catedras de materia {}", mat.codigo);
                    tracing::debug!("descripcion error: {err}");
                })?;

            let mut sql_elementos_materia = SqlElementosMateria {
                catedras: Vec::with_capacity(catedras.len()),
                ..Default::default()
            };

            let mut codigos_docentes = HashMap::new();

            for cat in catedras {
                sql_elementos_materia.catedras.push(cat.sql(mat.codigo));

                for (nombre, calificacion) in cat.docentes {
                    let codigo = codigos_docentes.entry(nombre).or_insert_with_key(|nombre| {
                        let codigo = Uuid::new_v4();
                        sql_elementos_materia
                            .docentes
                            .push(docente::sql(&codigo, &nombre, mat.codigo));

                        codigo
                    });

                    sql_elementos_materia
                        .rel_catedras_docentes
                        .push(docente::sql_rel_catedra(&codigo, &cat.codigo));

                    sql_elementos_materia
                        .calificaciones
                        .push(calificacion.sql(&codigo));
                }
            }

            let codigos_docentes = codigos_docentes
                .into_iter()
                .map(|(n, c)| ((mat.codigo, n), c))
                .collect();

            let output = TaskMateria {
                sql: sql_elementos_materia,
                codigos_docentes,
            };

            return Ok(output);
        }));
    }

    let mut sql_elementos_materias = SqlElementosMateria::default();
    let mut codigos_docentes = HashMap::new();

    for task in handles {
        if let Ok(materia) = task.await.unwrap() {
            sql_elementos_materias.catedras.extend(materia.sql.catedras);
            sql_elementos_materias.docentes.extend(materia.sql.docentes);
            sql_elementos_materias
                .rel_catedras_docentes
                .extend(materia.sql.rel_catedras_docentes);
            sql_elementos_materias
                .calificaciones
                .extend(materia.sql.calificaciones);

            codigos_docentes.extend(materia.codigos_docentes);
        }
    }

    let cuatrimestres: HashSet<_> = comentarios
        .keys()
        .map(|c| c.nombre_cuatrimestre.as_str())
        .collect();

    let sql_cuatrimestres: Vec<_> = cuatrimestres
        .into_iter()
        .map(comentario::sql_cuatrimestre)
        .collect();

    let mut sql_comentarios = Vec::with_capacity(comentarios.len());

    for (md, coms) in comentarios {
        if let Some(codigo_docente) = codigos_docentes.get(&(md.codigo_materia, md.nombre_docente))
        {
            sql_comentarios.extend(
                coms.iter().map(|c| {
                    comentario::sql_comentario(c, codigo_docente, &md.nombre_cuatrimestre)
                }),
            );
        }
    }

    Ok(vec![
        String::from_utf8_lossy(include_bytes!("../sql/schema.sql")).to_string(),
        "BEGIN;\n".to_string(),
        materia::sql_bulk_insert(&sql_materias),
        catedra::sql_bulk_insert(&sql_elementos_materias.catedras),
        docente::sql_bulk_insert_docentes(&sql_elementos_materias.docentes),
        docente::sql_bulk_insert_rel_catedras_docentes(
            &sql_elementos_materias.rel_catedras_docentes,
        ),
        docente::sql_bulk_insert_calificaciones(&sql_elementos_materias.calificaciones),
        comentario::sql_bulk_insert_cuatrimestres(&sql_cuatrimestres),
        // comentario::sql_bulk_insert_comentarios(&sql_comentarios),
        "COMMIT;\n".to_string(),
    ]
    .join("\n"))
}

fn crear_cliente_http() -> ClientWithMiddleware {
    let cache_mode = if cfg!(debug_assertions) {
        tracing::debug!("Forzando cache en las requests a Dolly");
        CacheMode::ForceCache
    } else {
        CacheMode::Default
    };

    ClientBuilder::new(Client::new())
        .with(Cache(HttpCache {
            mode: cache_mode,
            manager: CACacheManager::default(),
            options: None,
        }))
        .build()
}
