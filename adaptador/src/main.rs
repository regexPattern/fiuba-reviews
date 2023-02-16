use std::{fs::File, io::Write};

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let query_sql = adaptador::query_sql().await.unwrap();
    let mut archivo = File::create("init.sql").unwrap();
    archivo.write_all(query_sql.as_bytes()).unwrap();
}
