use std::net::SocketAddr;

use axum::{routing::get, Router, Server};
use sea_orm::Database;
use tower_http::cors::{Any, CorsLayer};

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/", get(root))
        .layer(CorsLayer::new().allow_origin(Any));

    let addr: SocketAddr = "0.0.0.0:5000".parse().unwrap();

    let _ = Database::connect("postgres://postgres:postgres@database:5432/postgres")
        .await
        .unwrap();

    Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn root() -> &'static str {
    "Hello World"
}
