#[tokio::main]
async fn main() {
    servidor::iniciar().await.unwrap();
}
