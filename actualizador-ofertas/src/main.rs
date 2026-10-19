use tracing_subscriber::{EnvFilter, fmt::time::ChronoLocal};

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env())
        .with_timer(ChronoLocal::new("%H:%M:%S".to_string()))
        .init();

    let archivos_planes = vec![
        "./tests/fixtures/ofertas/civil-1C-2026.json",
        "./tests/fixtures/ofertas/industrial-1C-2026.json",
        "./tests/fixtures/ofertas/informatica-1C-2026.json",
        "./tests/fixtures/ofertas/informatica-2C-2026.json",
    ];

    let ofertas_carreras = archivos_planes
        .iter()
        .map(|&archivo| {
            let source = std::fs::read_to_string(archivo).unwrap();
            serde_json::from_str(&source).unwrap()
        })
        .collect();

    actualizador_ofertas::armar_patches(ofertas_carreras).await;
}
