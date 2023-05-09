mod catedra;
mod comentario;
mod materia;

use std::{
    net::{Ipv4Addr, SocketAddr},
    num::ParseIntError,
    time::Duration,
};

use axum::{routing::get, Router, Server};
use sqlx::postgres::PgPoolOptions;
use thiserror::Error;

#[derive(Error, Debug)]
enum ErrorPuerto {
    #[error("variable de entorno `BACKEND_PORT` no encontrada")]
    VariableNoEncontrada,

    #[error("variable de entorno `BACKEND_PORT` invalida: `{0}`")]
    ValorInvalido(ParseIntError),
}

pub async fn escuchar() -> anyhow::Result<()> {
    let puerto: u16 = std::env::var("BACKEND_PORT")
        .map_err(|_| ErrorPuerto::VariableNoEncontrada)?
        .parse()
        .map_err(|err| ErrorPuerto::ValorInvalido(err))?;

    let pool = PgPoolOptions::new()
        .acquire_timeout(Duration::from_secs(5))
        .connect(&std::env::var("DATABASE_URL")?)
        .await?;

    tracing::info!("conexion establecida con la base de datos");

    let app = Router::new()
        .route("/materia", get(materia::index))
        .route("/materia/:codigo_materia/catedras", get(materia::catedras))
        .route(
            "/catedra/:codigo_catedra/docentes",
            get(catedra::docentes_con_comentarios),
        )
        .with_state(pool);

    let addr = SocketAddr::from(("0.0.0.0".parse::<Ipv4Addr>().unwrap(), puerto));

    tracing::info!("escuchando en `{addr}`");

    Ok(Server::bind(&addr).serve(app.into_make_service()).await?)
}
