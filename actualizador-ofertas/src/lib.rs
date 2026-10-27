use crate::patch::Patch;
pub use crate::siu::OfertaCarrera;

use sqlx::PgConnection;
use tracing::info;

mod catalogo;
mod patch;
mod siu;

pub async fn armar_patches(
    mut ofertas_carreras: Vec<OfertaCarrera>,
    db: &mut PgConnection,
) -> Vec<Patch> {
    siu::sanitizar_ofertas_carreras(&mut ofertas_carreras);
    catalogo::reconciliar_materias(&ofertas_carreras, db)
        .await
        .unwrap();

    let total_carreras = ofertas_carreras.len();
    let ofertas_materias = siu::consolidar_ofertas_materias(ofertas_carreras);

    info!(
        carreras = total_carreras,
        materias = ofertas_materias.len(),
        "ofertas consolidadas"
    );

    vec![]
}

pub mod gui;
