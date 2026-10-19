#![allow(dead_code)]

pub use crate::oferta::OfertaCarrera;
use crate::{oferta::OfertaMateria, patch::Patch};

use tracing::info;

mod oferta;
mod patch;

pub async fn armar_patches(ofertas_carreras: Vec<OfertaCarrera>) -> Vec<Patch> {
    let total_carreras = ofertas_carreras.len();
    let ofertas_materias: Vec<OfertaMateria> =
        oferta::consolidar_ofertas_materias(ofertas_carreras);

    info!(
        carreras = total_carreras,
        materias = ofertas_materias.len(),
        "ofertas consolidadas"
    );

    vec![]
}
