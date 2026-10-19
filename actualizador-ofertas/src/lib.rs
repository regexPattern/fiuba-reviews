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
    let catalogo_materias = catalogo::extraer_materias(&ofertas_carreras);

    info!(
        carreras = ofertas_carreras.len(),
        materias = catalogo_materias.len(),
        "catalogo de materias generado"
    );

    let reconciliadas = catalogo::reconciliar_materias(&catalogo_materias, db)
        .await
        .unwrap();

    info!(materias = reconciliadas.len(), "materias reconciliadas");

    let total_carreras = ofertas_carreras.len();
    let ofertas_materias = siu::consolidar_ofertas_materias(ofertas_carreras);

    info!(
        carreras = total_carreras,
        materias = ofertas_materias.len(),
        "ofertas consolidadas"
    );

    vec![]
}
