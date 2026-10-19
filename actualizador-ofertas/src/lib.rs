#![allow(dead_code)]

pub use crate::siu::OfertaCarrera;
use crate::patch::Patch;

use sqlx::{Pool, Postgres};
use tracing::info;

mod catalogo;
mod patch;
mod siu;

pub async fn armar_patches(
    mut ofertas_carreras: Vec<OfertaCarrera>,
    db: &mut Pool<Postgres>,
) -> Vec<Patch> {
    siu::sanitizar_ofertas_carreras(&mut ofertas_carreras);

    let catalogo_materias = catalogo::extraer_materias(&ofertas_carreras);
    catalogo::reconciliar_materias(&catalogo_materias, db).await;

    let total_carreras = ofertas_carreras.len();
    let ofertas_materias = siu::consolidar_ofertas_materias(ofertas_carreras);

    // Bueno, entonces ahora tenemos que hacer el siguiente paso, que es actualizar los codigos de
    // las materias que se han encontrado que no tenian codigo (prefijo 'COD'), y para estas
    // materias tenemos que hacer la "migracion" de los docentes y comentarios y calificaciones de
    // las materias anteriores. Es decir, si esta materia tenia equivalencias en un plan anterior,
    // se copian los docentes de estas equivalencias anteriores y se les pone el codigo de esta
    // nueva materia. Hay que considerar el caso en el que haya multiples equivalencias con los
    // mismos docente (por nombre, no tienen nombre_siu estos aun porque son de los planes
    // anteriores).

    info!(
        carreras = total_carreras,
        materias = ofertas_materias.len(),
        "ofertas consolidadas"
    );

    vec![]
}
