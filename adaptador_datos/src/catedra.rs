pub mod docente;

use std::{
    collections::{HashMap, HashSet},
    hash::{Hash, Hasher},
};

use docente::Docente;
use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use uuid::Uuid;

const URL_DESCARGA_BASE: &'static str = "https://dollyfiuba.com/analitics/cursos";

pub const TABLA: &'static str = r#"
CREATE TABLE IF NOT EXISTS catedra (
    codigo         TEXT PRIMARY KEY,
    nombre         TEXT NOT NULL,
    codigo_materia INTEGER REFERENCES materia(codigo) NOT NULL,
    promedio       DOUBLE PRECISION NOT NULL
);
"#;

pub const TABLA_RELACION_CATEDRA_DOCENTE: &'static str = r#"
CREATE TABLE IF NOT EXISTS catedra_docente (
    codigo_catedra TEXT REFERENCES catedra(codigo),
    codigo_docente TEXT REFERENCES docente(codigo),
    CONSTRAINT catedra_docente_pkey PRIMARY KEY (codigo_catedra, codigo_docente)
);
"#;

pub struct Catedra {
    pub codigo: Uuid,
    pub nombre: String,
    pub docentes: HashMap<String, Docente>,
    pub promedio: f64,
}

impl PartialEq for Catedra {
    fn eq(&self, other: &Self) -> bool {
        self.nombre == other.nombre
    }
}

impl Eq for Catedra {}

impl Hash for Catedra {
    fn hash<H: Hasher>(&self, state: &mut H) {
        self.nombre.hash(state);
    }
}

impl Catedra {
    pub async fn descargar_para_materia(
        client: &ClientWithMiddleware,
        codigo_materia: u32,
    ) -> anyhow::Result<HashSet<Catedra>> {
        #[derive(Deserialize)]
        struct WrapperRespuesta {
            #[serde(alias = "opciones")]
            catedras: Vec<CatedraRespuesta>,
        }

        #[derive(Deserialize)]
        struct CatedraRespuesta {
            pub nombre: String,
            pub docentes: HashMap<String, Docente>,
        }

        tracing::info!("descargando catedras de materia {}", codigo_materia);

        let respuesta = client
            .get(format!("{}/{}", URL_DESCARGA_BASE, codigo_materia))
            .send()
            .await?
            .text()
            .await?;

        let WrapperRespuesta { mut catedras } =
            serde_json::from_str(&respuesta).map_err(|err| SerdeError::new(respuesta, err))?;

        for catedra in &mut catedras {
            let mut nombres_docentes: Vec<_> = catedra.nombre.split("-").collect();
            nombres_docentes.sort();
            catedra.nombre = nombres_docentes.join("-");
        }

        let catedras = catedras.into_iter().map(|catedra| {
            let acumulado: f64 = catedra
                .docentes
                .values()
                .map(|docente| docente.promedio_calificaciones())
                .sum();

            let cantidad_docentes = catedra.docentes.len();
            let promedio = if cantidad_docentes > 0 {
                acumulado / cantidad_docentes as f64
            } else {
                0.0
            };

            Catedra {
                codigo: Uuid::new_v4(),
                nombre: catedra.nombre,
                docentes: catedra.docentes,
                promedio,
            }
        });

        Ok(catedras.collect())
    }

    pub fn sql(&self, codigo_materia: u32) -> String {
        format!(
            r#"
INSERT INTO catedra (codigo, codigo_materia, nombre, promedio)
VALUES ('{}', {}, '{}', {});
"#,
            self.codigo,
            codigo_materia,
            self.nombre.replace("'", "''"),
            self.promedio
        )
    }

    pub fn relacionar_docente_sql(&self, codigo_docente: &Uuid) -> String {
        format!(
            r#"
INSERT INTO catedra_docente (codigo_catedra, codigo_docente)
VALUES ('{}', '{}');
"#,
            self.codigo, codigo_docente
        )
    }
}
