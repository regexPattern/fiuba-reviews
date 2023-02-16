mod entidad;
mod remoto;

use std::{collections::HashMap, sync::Arc};

use anyhow::Result;
use entidad::{EntidadComentario, EntidadDocente};
use http_cache_reqwest::{CACacheManager, Cache, CacheMode, HttpCache};
use remoto::ComentarioRemoto;
use reqwest::Client;
use reqwest_middleware::{ClientBuilder, ClientWithMiddleware};
use reqwest_tracing::{SpanBackendWithUrl, TracingMiddleware};

pub async fn query_sql() -> Result<String> {
    let cliente_http = crear_client_http();

    let materias = remoto::descargar_materias(&cliente_http).await?;

    let rutinas = materias.into_iter().map(|materia| {
        let cliente_http = Arc::clone(&cliente_http);
        tokio::spawn(async move { materia.generar_adaptador(cliente_http.as_ref()).await })
    });

    let mut materias = vec![];
    let mut catedras = vec![];
    let mut mapa_materias_docentes = HashMap::new();

    for rutina in rutinas {
        let adaptador = match rutina.await.unwrap() {
            Ok(adaptador) => adaptador,
            Err(_) => continue,
        };

        let codigo_materia = adaptador.entidad_materia.codigo;

        materias.push(adaptador.entidad_materia);
        catedras.extend(adaptador.entidades_catedras);

        for docente in adaptador.entidades_docentes {
            mapa_materias_docentes.insert(
                (codigo_materia, docente.calificacion.nombre.clone()),
                docente,
            );
        }
    }

    let comentarios = remoto::descargar_comentarios(&cliente_http).await?;
    let comentarios = generar_entidades_comentarios(comentarios, &mapa_materias_docentes);

    Ok(entidad::exportar_query_sql(
        materias,
        catedras,
        mapa_materias_docentes.into_values().collect::<Vec<_>>(),
        comentarios,
    ))
}

fn crear_client_http() -> Arc<ClientWithMiddleware> {
    Arc::new(
        ClientBuilder::new(Client::new())
            .with(Cache(HttpCache {
                mode: CacheMode::ForceCache,
                manager: CACacheManager::default(),
                options: None,
            }))
            .with(TracingMiddleware::<SpanBackendWithUrl>::new())
            .build(),
    )
}

fn generar_entidades_comentarios(
    comentarios: Vec<ComentarioRemoto>,
    mapa_materias_docentes: &HashMap<(u32, String), EntidadDocente>,
) -> Vec<EntidadComentario> {
    let mut entidades = vec![];

    for comentario in comentarios {
        let codigo_docente = match mapa_materias_docentes
            .get(&(comentario.codigo_materia, comentario.nombre_docente))
        {
            Some(docente) => docente.codigo,
            None => continue,
        };

        for contenido in comentario.contenido_comentarios {
            entidades.push(EntidadComentario {
                cuatrimestre: comentario.cuatrimestre.clone(),
                contenido,
                codigo_docente,
            });
        }
    }

    entidades
}
