use std::net::SocketAddr;

use axum::{routing::get, Router, Server};

#[tokio::main]
async fn main() {
    let app = Router::new().route("/", get(root));
    let addr: SocketAddr = "0.0.0.0:5000".parse().unwrap();

    Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn root() -> &'static str {
    "Another message"
}
