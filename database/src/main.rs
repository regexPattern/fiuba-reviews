use std::{fs::File, io::Write};

#[tokio::main]
async fn main() {
    let mut outfile = File::create("init.sql").unwrap();
    outfile
        .write_all(
            "\
CREATE TABLE IF NOT EXISTS materias (
    codigo INTEGER PRIMARY KEY,
    nombre TEXT NOT NULL
);

INSERT INTO materias (codigo, nombre) VALUES (6103, 'ANALISIS MATEMATICO II A');
INSERT INTO materias (codigo, nombre) VALUES (6106, 'PROBABILIDAD Y ESTADISTICA A');
"
            .as_bytes(),
        )
        .unwrap();
}
