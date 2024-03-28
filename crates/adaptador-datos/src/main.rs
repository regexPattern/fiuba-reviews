use clap::{Parser, Subcommand};

#[derive(Parser)]
struct Cli {
    #[clap(subcommand)]
    command: Command,
}

#[derive(Subcommand)]
enum Command {
    Init,
    Update {
        #[arg(long)]
        dump: bool,
    },
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    tracing_subscriber::fmt::init();

    match cli.command {
        Command::Init => {
            let query = adaptador_datos::init_query().await?;
            println!("{query}");
        }
        Command::Update { dump: _dump } => todo!(),
    };

    Ok(())
}
