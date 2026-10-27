pub use crate::propuesta::{CandidatoDocente, DocentePendiente, DocenteResuelto, PropuestaMateria};
pub use crate::siu::OfertaCarrera;

use sqlx::PgConnection;
use tracing::info;

mod catalogo;
mod propuesta;
mod siu;

pub mod gui;

pub async fn proponer_actualizaciones(
    mut ofertas_carreras: Vec<OfertaCarrera>,
    db: &mut PgConnection,
) -> Vec<PropuestaMateria> {
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

    propuesta::proponer(&ofertas_materias, db).await.unwrap()
}
