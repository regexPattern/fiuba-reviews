use std::collections::HashSet;

use sqlx::{Pool, Postgres};

use crate::OfertaCarrera;

#[derive(Eq, PartialEq, Hash, Debug)]
pub struct MetadatosMateria {
    codigo: String,
    nombre: String,
}

pub fn extraer_materias(ofertas_carreras: &[OfertaCarrera]) -> Vec<MetadatosMateria> {
    let metadatos: HashSet<_> = ofertas_carreras
        .iter()
        .flat_map(|oferta| &oferta.materias)
        .map(|materia| MetadatosMateria {
            codigo: materia.codigo.clone(),
            nombre: materia.nombre.clone(),
        })
        .collect();

    metadatos.into_iter().collect()
}

pub async fn reconciliar_materias(catalogo: &[MetadatosMateria], db: &mut Pool<Postgres>) {}
