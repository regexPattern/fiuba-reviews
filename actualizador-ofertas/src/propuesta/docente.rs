use std::collections::HashMap;

use sqlx::{PgConnection, types::Uuid};

use crate::siu::OfertaMateria;

const SIMILITUD_MINIMA: f32 = 0.5;
const MAX_CANDIDATOS: i64 = 5;

/// Docente del SIU que todavía no está vinculado a ningún docente de la DB.
///
/// Sin candidatos, lo más probable es que se trate de un docente nuevo.
#[derive(Debug, PartialEq)]
pub struct DocentePendiente {
    pub nombre_siu: String,
    pub candidatos: Vec<CandidatoDocente>,
}

/// Docente del SIU ya vinculado a un docente de la DB. Se conservan los candidatos por si se
/// quiere cambiar el vínculo.
#[derive(Debug, PartialEq)]
pub struct DocenteResuelto {
    pub nombre_siu: String,
    pub vinculado: CandidatoDocente,
    pub candidatos: Vec<CandidatoDocente>,
}

/// Docente de la DB que podría corresponder a un docente del SIU de su misma materia.
#[derive(Debug, PartialEq)]
pub struct CandidatoDocente {
    pub codigo: Uuid,
    pub nombre: String,
    /// Nombre del SIU con el que ya está vinculado, si lo está.
    pub nombre_siu_actual: Option<String>,
    pub score: f32,
}

/// Docentes del SIU de una materia, separados según estén o no vinculados.
#[derive(Debug, Default, PartialEq)]
pub(super) struct PropuestaDocentes {
    pub pendientes: Vec<DocentePendiente>,
    pub resueltos: Vec<DocenteResuelto>,
}

/// Fila de `candidatos.sql`: un docente del SIU junto a uno de sus candidatos de la DB.
///
/// Las columnas del candidato (`codigo`, `nombre`, `nombre_siu_actual`, `score`, `vinculado`)
/// son `None` cuando el docente del SIU no tiene vínculo ni candidatos.
#[derive(Debug, sqlx::FromRow)]
struct FilaCandidato {
    codigo_materia: String,
    nombre_siu: String,
    codigo: Option<Uuid>,
    nombre: Option<String>,
    nombre_siu_actual: Option<String>,
    score: Option<f32>,
    vinculado: Option<bool>,
}

/// Propueta de los docentes de cada materia, indexados por código de materia.
///
/// Un docente está resuelto si ya existe u docente de su materia ya que tiene su mismo
/// `nombre_siu`.
pub(super) async fn buscar(
    ofertas_materias: &[OfertaMateria],
    conn: &mut PgConnection,
) -> Result<HashMap<String, PropuestaDocentes>, sqlx::Error> {
    let (codigos, nombres): (Vec<&str>, Vec<&str>) = ofertas_materias
        .iter()
        .flat_map(|o| o.nombres_docentes().into_iter().map(|n| (o.codigo(), n)))
        .unzip();

    let filas: Vec<FilaCandidato> = sqlx::query_as(include_str!("docente/candidatos.sql"))
        .bind(codigos)
        .bind(nombres)
        .bind(SIMILITUD_MINIMA)
        .bind(MAX_CANDIDATOS)
        .fetch_all(conn)
        .await?;

    Ok(agrupar(filas))
}

/// Docente del SIU mientras se leen las filas de sus candidatos.
#[derive(Debug)]
struct DocenteSiu {
    nombre_siu: String,
    vinculado: Option<CandidatoDocente>,
    candidatos: Vec<CandidatoDocente>,
}

/// Agrupa las filas (ordenadas por docente del SIU, con el vinculado primero).
fn agrupar(filas: Vec<FilaCandidato>) -> HashMap<String, PropuestaDocentes> {
    let mut por_materia: HashMap<String, Vec<DocenteSiu>> = HashMap::new();

    for fila in filas {
        let docentes = por_materia.entry(fila.codigo_materia).or_default();

        if docentes
            .last()
            .is_none_or(|d| d.nombre_siu != fila.nombre_siu)
        {
            docentes.push(DocenteSiu {
                nombre_siu: fila.nombre_siu,
                vinculado: None,
                candidatos: Vec::new(),
            });
        }
        let docente = docentes.last_mut().unwrap();

        if let (Some(codigo), Some(nombre), Some(score)) = (fila.codigo, fila.nombre, fila.score) {
            let candidato = CandidatoDocente {
                codigo,
                nombre,
                nombre_siu_actual: fila.nombre_siu_actual,
                score,
            };
            if fila.vinculado == Some(true) && docente.vinculado.is_none() {
                docente.vinculado = Some(candidato);
            } else {
                docente.candidatos.push(candidato);
            }
        }
    }

    por_materia
        .into_iter()
        .map(|(codigo_materia, docentes)| {
            let mut materia = PropuestaDocentes::default();
            for d in docentes {
                match d.vinculado {
                    Some(vinculado) => materia.resueltos.push(DocenteResuelto {
                        nombre_siu: d.nombre_siu,
                        vinculado,
                        candidatos: d.candidatos,
                    }),
                    None => materia.pendientes.push(DocentePendiente {
                        nombre_siu: d.nombre_siu,
                        candidatos: d.candidatos,
                    }),
                }
            }
            (codigo_materia, materia)
        })
        .collect()
}

#[cfg(test)]
mod tests {
    use super::*;
    use sqlx::PgPool;

    fn oferta(codigo: &str, docentes: &[&str]) -> OfertaMateria {
        let docentes: Vec<_> = docentes
            .iter()
            .map(|nombre| serde_json::json!({"nombre": nombre, "rol": null}))
            .collect();

        serde_json::from_value(serde_json::json!({
            "carreras": ["Ingeniería Civil"],
            "anio": 2026,
            "cuatrimestre": 1,
            "materia": {
                "codigo": codigo,
                "nombre": "ÁLGEBRA LINEAL",
                "comisiones": [{"nombre": "A", "docentes": docentes, "clases": []}]
            }
        }))
        .unwrap()
    }

    async fn insertar_materia(pool: &PgPool, codigo: &str) {
        sqlx::query("INSERT INTO materia (codigo, nombre) VALUES ($1, 'ÁLGEBRA LINEAL')")
            .bind(codigo)
            .execute(pool)
            .await
            .unwrap();
    }

    async fn insertar_docente(
        pool: &PgPool,
        materia: &str,
        nombre: &str,
        nombre_siu: Option<&str>,
    ) {
        sqlx::query("INSERT INTO docente (nombre, codigo_materia, nombre_siu) VALUES ($1, $2, $3)")
            .bind(nombre)
            .bind(materia)
            .bind(nombre_siu)
            .execute(pool)
            .await
            .unwrap();
    }

    fn fila_candidato(
        siu: &str,
        docente: Option<(u128, &str)>,
        score: f32,
        vinculado: bool,
    ) -> FilaCandidato {
        FilaCandidato {
            codigo_materia: "CB002".into(),
            nombre_siu: siu.into(),
            codigo: docente.map(|(codigo, _)| Uuid::from_u128(codigo)),
            nombre: docente.map(|(_, nombre)| nombre.into()),
            nombre_siu_actual: None,
            score: docente.map(|_| score),
            vinculado: docente.map(|_| vinculado),
        }
    }

    #[test]
    fn agrupar_junta_las_filas_de_un_mismo_docente_del_siu() {
        let agrupado = agrupar(vec![
            fila_candidato("PALACIOS PATRICIA", Some((1, "Palacios")), 1.0, false),
            fila_candidato("PALACIOS PATRICIA", Some((2, "Palacio")), 0.7, false),
            fila_candidato("SOSA MARINA", None, 0.0, false),
        ]);

        let PropuestaDocentes {
            pendientes,
            resueltos,
        } = &agrupado["CB002"];
        assert!(resueltos.is_empty());
        assert_eq!(pendientes.len(), 2);
        assert_eq!(pendientes[0].nombre_siu, "PALACIOS PATRICIA");
        assert_eq!(pendientes[0].candidatos.len(), 2);
        assert_eq!(pendientes[0].candidatos[0].codigo, Uuid::from_u128(1));
        assert_eq!(pendientes[1].nombre_siu, "SOSA MARINA");
        assert!(pendientes[1].candidatos.is_empty());
    }

    #[test]
    fn agrupar_separa_el_vinculado_de_los_candidatos() {
        let agrupado = agrupar(vec![
            fila_candidato("VARGAS GABRIELA", Some((1, "Vargas")), 1.0, true),
            fila_candidato("VARGAS GABRIELA", Some((2, "Vargas G")), 0.8, false),
        ]);

        let materia = &agrupado["CB002"];
        assert!(materia.pendientes.is_empty());
        let resuelto = &materia.resueltos[0];
        assert_eq!(resuelto.vinculado.codigo, Uuid::from_u128(1));
        assert_eq!(resuelto.candidatos.len(), 1);
        assert_eq!(resuelto.candidatos[0].codigo, Uuid::from_u128(2));
    }

    #[sqlx::test]
    async fn un_docente_con_el_mismo_nombre_siu_esta_resuelto(pool: PgPool) {
        insertar_materia(&pool, "CB002").await;
        insertar_docente(&pool, "CB002", "Vargas", Some("VARGAS GABRIELA")).await;
        let ofertas = [oferta("CB002", &["VARGAS GABRIELA", "SOSA MARINA"])];

        let mut conn = pool.acquire().await.unwrap();
        let materias = buscar(&ofertas, &mut conn).await.unwrap();

        let materia = &materias["CB002"];
        let pendientes: Vec<_> = materia.pendientes.iter().map(|p| &*p.nombre_siu).collect();
        assert_eq!(pendientes, ["SOSA MARINA"]);
        assert_eq!(materia.resueltos.len(), 1);
        assert_eq!(materia.resueltos[0].nombre_siu, "VARGAS GABRIELA");
        assert_eq!(materia.resueltos[0].vinculado.nombre, "Vargas");
    }

    #[sqlx::test]
    async fn el_vinculo_se_conserva_aunque_los_nombres_no_se_parezcan(pool: PgPool) {
        insertar_materia(&pool, "CB002").await;
        insertar_docente(&pool, "CB002", "Zeta", Some("SOSA MARINA")).await;
        insertar_docente(&pool, "CB002", "Sosa", None).await;
        let ofertas = [oferta("CB002", &["SOSA MARINA"])];

        let mut conn = pool.acquire().await.unwrap();
        let materias = buscar(&ofertas, &mut conn).await.unwrap();

        let resuelto = &materias["CB002"].resueltos[0];
        assert_eq!(resuelto.vinculado.nombre, "Zeta");
        let candidatos: Vec<_> = resuelto.candidatos.iter().map(|c| &*c.nombre).collect();
        assert_eq!(candidatos, ["Sosa"]);
    }

    #[sqlx::test]
    async fn el_nombre_siu_se_compara_sin_acentos_mayusculas_ni_espacios_repetidos(pool: PgPool) {
        insertar_materia(&pool, "CB002").await;
        insertar_docente(
            &pool,
            "CB002",
            "Ponce",
            Some("ponce  de leon   robles fabián"),
        )
        .await;
        let ofertas = [oferta("CB002", &["PONCE DE LEON ROBLES FABIAN"])];

        let mut conn = pool.acquire().await.unwrap();
        let materias = buscar(&ofertas, &mut conn).await.unwrap();

        assert!(materias["CB002"].pendientes.is_empty());
        assert_eq!(materias["CB002"].resueltos.len(), 1);
    }

    #[sqlx::test]
    async fn los_candidatos_se_ordenan_por_similitud_y_solo_son_de_su_materia(pool: PgPool) {
        for materia in ["CB002", "CB003"] {
            insertar_materia(&pool, materia).await;
        }
        insertar_docente(&pool, "CB002", "Palacios", None).await;
        insertar_docente(&pool, "CB002", "Palacios Marina", None).await;
        insertar_docente(&pool, "CB002", "Sosa", None).await;
        insertar_docente(&pool, "CB003", "Palacios", None).await;
        let ofertas = [oferta("CB002", &["PALACIOS PATRICIA ELIZABETH"])];

        let mut conn = pool.acquire().await.unwrap();
        let materias = buscar(&ofertas, &mut conn).await.unwrap();

        let candidatos = &materias["CB002"].pendientes[0].candidatos;
        let nombres: Vec<_> = candidatos.iter().map(|c| &*c.nombre).collect();
        assert_eq!(nombres, ["Palacios", "Palacios Marina"]);
        assert!(candidatos[0].score >= candidatos[1].score);
    }

    #[sqlx::test]
    async fn los_candidatos_ya_vinculados_a_otro_nombre_siu_se_incluyen(pool: PgPool) {
        insertar_materia(&pool, "CB002").await;
        insertar_docente(&pool, "CB002", "Vargas", Some("VARGAS GABRIELA")).await;
        let ofertas = [oferta("CB002", &["VARGAS GABRIELA INES"])];

        let mut conn = pool.acquire().await.unwrap();
        let materias = buscar(&ofertas, &mut conn).await.unwrap();

        let candidato = &materias["CB002"].pendientes[0].candidatos[0];
        assert_eq!(candidato.nombre, "Vargas");
        assert_eq!(
            candidato.nombre_siu_actual.as_deref(),
            Some("VARGAS GABRIELA")
        );
    }

    #[sqlx::test]
    async fn un_docente_sin_similares_queda_pendiente_sin_candidatos(pool: PgPool) {
        insertar_materia(&pool, "CB002").await;
        insertar_docente(&pool, "CB002", "Vargas", None).await;
        let ofertas = [oferta("CB002", &["PERALTA NORA BEATRIZ"])];

        let mut conn = pool.acquire().await.unwrap();
        let materias = buscar(&ofertas, &mut conn).await.unwrap();

        assert_eq!(
            materias["CB002"],
            PropuestaDocentes {
                pendientes: vec![DocentePendiente {
                    nombre_siu: "PERALTA NORA BEATRIZ".into(),
                    candidatos: vec![]
                }],
                resueltos: vec![],
            }
        );
    }

    #[sqlx::test]
    async fn se_limita_la_cantidad_de_candidatos(pool: PgPool) {
        insertar_materia(&pool, "CB002").await;
        for i in 0..=MAX_CANDIDATOS {
            insertar_docente(&pool, "CB002", &format!("Sosa {i}"), None).await;
        }
        let ofertas = [oferta("CB002", &["SOSA MARINA"])];

        let mut conn = pool.acquire().await.unwrap();
        let materias = buscar(&ofertas, &mut conn).await.unwrap();

        assert_eq!(
            materias["CB002"].pendientes[0].candidatos.len() as i64,
            MAX_CANDIDATOS
        );
    }
}
