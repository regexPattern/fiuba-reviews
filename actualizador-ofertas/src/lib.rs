#![allow(dead_code)]

pub use crate::oferta::OfertaCarrera;
use crate::{oferta::OfertaMateria, patch::Patch};

mod oferta;
mod patch;

pub async fn armar_patches(ofertas_carreras: Vec<OfertaCarrera>) -> Vec<Patch> {
    let ofertas_materias: Vec<OfertaMateria> = oferta::consolidar_ofertas_materias(ofertas_carreras);

    dbg!(ofertas_materias);

    vec![]
}
