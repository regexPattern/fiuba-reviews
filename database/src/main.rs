use std::{fs::File, io::Write};

#[tokio::main]
async fn main() {
    let mut outfile = File::create("init.sql").unwrap();
    outfile
        .write_all("CREATE TABLE IF NOT EXISTS materias ();".as_bytes())
        .unwrap();
}
