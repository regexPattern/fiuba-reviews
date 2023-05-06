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
        .connect(&std::env::var("DATABASE_URL").unwrap())
        .await?;

    tracing::info!("establecida conexion con la base de datos");

    let app = Router::new()
        .route("/materias", get(materias::get_all))
        .route("/materias/:codigo", get(materias::by_codigo))
        .route("/catedras/:codigo_materia", get(catedras::by_materia))
        .route(
            "/catedras/:codigo_materia/docentes",
            get(docentes::by_catedra),
        )
        .route("/docentes/:codigo", get(docentes::by_codigo))
        .route("/comentarios/:codigo_docente", get(comentarios::by_docente))
        .with_state(pool);

    let addr: SocketAddr = "0.0.0.0:5000".parse().unwrap();

    tracing::info!("escuchando en '{addr}'");

    Ok(Server::bind(&addr).serve(app.into_make_service()).await?)
}
