use std::{fs::File, io::Write};

use http_cache_reqwest::{CACacheManager, Cache, CacheMode, HttpCache};
use reqwest::Client;
use reqwest_middleware::ClientBuilder;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let cliente_http = ClientBuilder::new(Client::new())
        .with(Cache(HttpCache {
            mode: CacheMode::Default,
            manager: CACacheManager::default(),
            options: None,
        }))
        .build();

    let init_sql = migrador_datos::descargar(&cliente_http).await?;
    let mut archivo = File::create("init.sql")?;
    archivo.write_all(init_sql.as_bytes())?;

    Ok(())
}
