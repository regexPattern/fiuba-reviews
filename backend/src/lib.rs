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
        .route("/materias", get(materias::listar))
        .route("/materias/:codigo_materia", get(materias::informacion))
        .route(
            "/materias/:codigo_materia/catedras",
            get(catedras::por_materia),
        )
        .route("/catedras/:codigo_catedra", get(catedras::informacion))
        .route(
            "/catedras/:codigo_catedra/docentes",
            get(docentes::por_catedra),
        )
        .route("/docentes/:codigo_docente", get(docentes::informacion))
        .route(
            "/docentes/:codigo_docente/catedras",
            get(catedras::por_docente),
        )
        .route(
            "/comentarios/:codigo_docente",
            get(comentarios::por_docente),
        )
        .with_state(pool);

    let addr: SocketAddr = "0.0.0.0:5000".parse().unwrap();

    tracing::info!("escuchando en '{addr}'");

    Ok(Server::bind(&addr).serve(app.into_make_service()).await?)
}
