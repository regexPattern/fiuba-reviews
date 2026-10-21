use std::collections::HashSet;

use sqlx::PgConnection;
use tracing::{debug, info};

use crate::OfertaCarrera;

#[derive(Eq, PartialEq, Hash, Debug)]
struct MetadatosMateria {
    codigo: String,
    nombre: String,
}

/// Materia vigente cuyo código provisional fue reemplazado por el oficial del SIU.
#[derive(Debug, Eq, PartialEq, sqlx::FromRow)]
pub struct MateriaReconciliada {
    codigo: String,
    codigo_anterior: String,
    nombre: String,
    codigos_equivalencias: Vec<String>,
    docentes_migrados: i32,
    calificaciones_migradas: i32,
    comentarios_migrados: i32,
}

fn extraer_materias(ofertas_carreras: &[OfertaCarrera]) -> Vec<MetadatosMateria> {
    let metadatos: HashSet<_> = ofertas_carreras
        .iter()
        .flat_map(|o| &o.materias)
        .map(|m| MetadatosMateria {
            codigo: m.codigo.clone(),
            nombre: m.nombre.clone(),
        })
        .collect();

    metadatos.into_iter().collect()
}

/// Extrae el catálogo de todas las ofertas sanitizadas y reconcilia sus materias en la DB.
///
/// No abre ni cierra transacciones: quien llama decide si confirmar o descartar los cambios.
/// Incluye el refresco de las materialized views.
pub async fn reconciliar_materias(
    ofertas_carreras: &[OfertaCarrera],
    conn: &mut PgConnection,
) -> Result<Vec<MateriaReconciliada>, sqlx::Error> {
    let catalogo = extraer_materias(ofertas_carreras);

    debug!(
        carreras = ofertas_carreras.len(),
        materias = catalogo.len(),
        "catalogo de materias generado"
    );

    let nombres: Vec<&str> = catalogo.iter().map(|m| m.nombre.as_str()).collect();
    let codigos: Vec<&str> = catalogo.iter().map(|m| m.codigo.as_str()).collect();

    let reconciliadas: Vec<MateriaReconciliada> =
        sqlx::query_as(include_str!("catalogo/reconciliar_materias.sql"))
            .bind(nombres)
            .bind(codigos)
            .fetch_all(&mut *conn)
            .await?;

    if !reconciliadas.is_empty() {
        sqlx::query("SELECT refresh_materialized_views()")
            .execute(&mut *conn)
            .await?;
    }

    for materia in &reconciliadas {
        debug!(
            materia = %materia.codigo,
            equivalencias = ?materia.codigos_equivalencias,
            docentes = materia.docentes_migrados,
            calificaciones = materia.calificaciones_migradas,
            comentarios = materia.comentarios_migrados,
            "materia reconciliada"
        );
    }

    info!(materias = reconciliadas.len(), "materias reconciliadas");

    Ok(reconciliadas)
}
