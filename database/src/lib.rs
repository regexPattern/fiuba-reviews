#![allow(unused, dead_code)]

mod dolly;

use dolly::{Catedra, ComentariosCuatrimestre, Docente, Materia};
use http_cache_reqwest::{CACacheManager, Cache, CacheMode, HttpCache};
use reqwest::Client;
use reqwest_middleware::ClientBuilder;

pub async fn run() -> anyhow::Result<String> {
    let client = ClientBuilder::new(Client::new())
        .with(Cache(HttpCache {
            mode: CacheMode::ForceCache,
            manager: CACacheManager::default(),
            options: None,
        }))
        .build();

    let comentarios = ComentariosCuatrimestre::fetch_all(&client).await?;
    let materias = Materia::fetch_all(&client).await?;

    let mut init_sql_buffer = vec![
        Materia::table_query(),
        Catedra::table_query(),
        Docente::table_query(),
    ];

    for materia in materias {
        init_sql_buffer.push(materia.insert_query());
        for catedra in materia.catedras {
            init_sql_buffer.push(catedra.insert_query(&materia.codigo));
        }
    }

    Ok(init_sql_buffer.join("\n\n"))
}
