mod docente;

pub use docente::{CandidatoDocente, DocentePendiente, DocenteResuelto};

use crate::siu::OfertaMateria;
use sqlx::PgConnection;

#[derive(Debug, PartialEq)]
pub struct PropuestaMateria {
    pub codigo: String,
    pub docentes_pendientes: Vec<DocentePendiente>,
    pub docentes_resueltos: Vec<DocenteResuelto>,
}

/// Propone una actualización por cada materia con docentes en la oferta del SIU, tanto los
/// pendientes como los ya resueltos.
///
/// Las propuestas salen de menor a mayor cantidad de docentes pendientes y, a igual cantidad,
/// por código de materia.
pub async fn proponer(
    ofertas_materias: &[OfertaMateria],
    conn: &mut PgConnection,
) -> Result<Vec<PropuestaMateria>, sqlx::Error> {
    let materias = docente::buscar(ofertas_materias, conn).await?;

    let mut propuestas: Vec<_> = materias
        .into_iter()
        .map(|(codigo, docentes)| PropuestaMateria {
            codigo,
            docentes_pendientes: docentes.pendientes,
            docentes_resueltos: docentes.resueltos,
        })
        .collect();

    propuestas.sort_by(|a, b| {
        (a.docentes_pendientes.len(), &a.codigo).cmp(&(b.docentes_pendientes.len(), &b.codigo))
    });

    Ok(propuestas)
}
