mod catedras;
mod comentarios;
mod docentes;
mod materias;

use std::{net::SocketAddr, time::Duration};

use axum::{routing::get, Router, Server};
use sqlx::postgres::PgPoolOptions;

pub async fn run() -> anyhow::Result<()> {
    let pool = PgPoolOptions::new()
        .acquire_timeout(Duration::from_secs(5))
        .connect(&std::env::var("DATABASE_URL")?)
        .await?;

    tracing::info!("establecida conexion con la base de datos");

    let app = Router::new()
        .route("/materias", get(materias::get_all))
        .route("/materias/:codigo", get(materias::by_codigo))
        .route("/materias/:codigo/catedras", get(catedras::by_materia))
        .route(
            "/catedras/:codigo_catedra/docentes",
            get(docentes::by_catedra),
        )
        .route("/docentes/:codigo", get(docentes::by_codigo))
        .route("/docentes/:codigo/catedras", get(catedras::by_docente))
        .route("/comentarios/:codigo_docente", get(comentarios::by_docente))
        .with_state(pool);

    let addr: SocketAddr = "0.0.0.0:5000".parse().unwrap();

    tracing::info!("escuchando en '{addr}'");

    Ok(Server::bind(&addr).serve(app.into_make_service()).await?)
}
