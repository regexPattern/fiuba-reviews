#![allow(unused, dead_code)]

mod db;

use db::materias::Entity as Materia;
use sea_orm::{Database, EntityTrait};

pub async fn iniciar() -> anyhow::Result<()> {
    let db = Database::connect("postgres://postgres:postgres@localhost:5432/postgres")
        .await
        .unwrap();

    let materia = Materia::find().one(&db).await?;
    dbg!(materia);

    Ok(())
}
