mod catedra;
mod docente;
mod materia;

use http_cache::{CACacheManager, CacheMode, HttpCache};
use http_cache_reqwest::Cache;
use reqwest::Client;
use reqwest_middleware::ClientBuilder;
use std::sync::Arc;

use materia::Materia;

pub async fn init_query() -> anyhow::Result<String> {
    let cliente_http = Arc::new(
        ClientBuilder::new(Client::new())
            .with(Cache(HttpCache {
                mode: CacheMode::ForceCache,
                manager: CACacheManager::default(),
                options: None,
            }))
            .build(),
    );

    let materias = Materia::descargar_todas(&cliente_http).await?;
    let mut materias_inserts = Vec::with_capacity(materias.len());

    let mut handles = Vec::with_capacity(materias.len());

    for materia in materias.into_iter().take(10) {
        materias_inserts.push(materia.sql());

        let cliente_http = Arc::clone(&cliente_http);

        handles.push(tokio::spawn(async move {
            let catedras = materia.descargar_catedras(&cliente_http).await.unwrap();

            for catedra in catedras {
            }
            return materia.codigo;
        }));
    }

    for handle in handles {
        let _codigo_materia = handle.await.unwrap();
    }

    dbg!(materias_inserts);

    Ok("".into())
}
