use std::{
    collections::{HashMap, HashSet},
    hash::{Hash, Hasher},
};

use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use uuid::Uuid;

use crate::{materias::Materia, sql::Sql};

pub const CREACION_TABLA_CATEDRAS: &str = r#"
CREATE TABLE IF NOT EXISTS Catedra(
    codigo         TEXT PRIMARY KEY,
    codigo_materia INTEGER REFERENCES Materia(codigo) NOT NULL
);
"#;

pub const CREACION_TABLA_DOCENTES: &str = r#"
CREATE TABLE IF NOT EXISTS Docente(
    codigo      TEXT PRIMARY KEY,
    nombre      TEXT NOT NULL,
    descripcion TEXT,

    -- Cantidad de comentarios del docente al momento de la ultima actualizacion de la descripcion.
    comentarios_ultima_descripcion INT DEFAULT 0 NOT NULL
);
"#;

pub const CREACION_TABLA_CATEDRA_DOCENTE: &str = r#"
CREATE TABLE IF NOT EXISTS CatedraDocente(
    codigo_catedra TEXT REFERENCES Catedra(codigo),
    codigo_docente TEXT REFERENCES Docente(codigo),
    CONSTRAINT catedra_docente_pkey PRIMARY KEY (codigo_catedra, codigo_docente)
);
"#;

pub const CREACION_TABLA_CALIFICACION: &str = r#"
CREATE TABLE IF NOT EXISTS Calificacion(
    codigo                UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    codigo_docente        TEXT REFERENCES Docente(codigo) NOT NULL,
    acepta_critica        DOUBLE PRECISION NOT NULL,
    asistencia            DOUBLE PRECISION NOT NULL,
    buen_trato            DOUBLE PRECISION NOT NULL,
    claridad              DOUBLE PRECISION NOT NULL,
    clase_organizada      DOUBLE PRECISION NOT NULL,
    cumple_horarios       DOUBLE PRECISION NOT NULL,
    fomenta_participacion DOUBLE PRECISION NOT NULL,
    panorama_amplio       DOUBLE PRECISION NOT NULL,
    responde_mails        DOUBLE PRECISION NOT NULL
);
"#;

const URL_DESCARGA_CATEDRAS: &str = "https://dollyfiuba.com/analitics/cursos";

#[derive(Debug)]
pub struct Catedra {
    pub codigo: Uuid,
    pub nombre: String,
    pub docentes: HashMap<String, Calificacion>,
}

#[derive(Deserialize, Default, Debug)]
pub struct Calificacion {
    respuestas: usize,

    acepta_critica: Option<f64>,
    asistencia: Option<f64>,
    buen_trato: Option<f64>,
    claridad: Option<f64>,
    clase_organizada: Option<f64>,
    cumple_horarios: Option<f64>,
    fomenta_participacion: Option<f64>,
    panorama_amplio: Option<f64>,
    responde_mails: Option<f64>,
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

impl Materia {
    pub async fn catedras(
        &self,
        http: &ClientWithMiddleware,
    ) -> anyhow::Result<impl Iterator<Item = Catedra>> {
        #[derive(Deserialize)]
        struct Catedras {
            #[serde(alias = "opciones")]
            catedras: Vec<CatedraDolly>,
        }

        #[derive(Deserialize)]
        struct CatedraDolly {
            pub nombre: String,
            pub docentes: HashMap<String, Calificacion>,
        }

        tracing::info!("descargando catedras de materia {}", self.codigo);

        let res = http
            .get(format!("{}/{}", URL_DESCARGA_CATEDRAS, self.codigo))
            .send()
            .await?;

        let data = res.text().await?;

        let Catedras { mut catedras } =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        for catedra in &mut catedras {
            let mut nombres_docentes: Vec<_> = catedra.nombre.split('-').collect();
            nombres_docentes.sort();
            catedra.nombre = nombres_docentes.join("-").to_uppercase();
        }

        let catedras: HashSet<_> = catedras
            .into_iter()
            .map(|catedra| Catedra {
                codigo: Uuid::new_v4(),
                nombre: catedra.nombre,
                docentes: catedra.docentes,
            })
            .collect();

        Ok(catedras.into_iter())
    }
}

impl Catedra {
    pub fn query_sql(&self, codigo_materia: u32) -> String {
        format!(
            r#"
INSERT INTO Catedra(codigo, codigo_materia)
VALUES ('{}', {});
"#,
            self.codigo, codigo_materia,
        )
    }

    pub fn relacion_con_docente_query_sql(&self, codigo_docente: &Uuid) -> String {
        format!(
            r#"
INSERT INTO CatedraDocente(codigo_catedra, codigo_docente)
VALUES ('{}', '{}');
"#,
            self.codigo, codigo_docente
        )
    }
}

impl Calificacion {
    pub fn query_sql(&self, nombre_docente: &str, codigo_docente: Uuid) -> String {
        let docente = format!(
            r#"
INSERT INTO Docente(codigo, nombre) VALUES ('{codigo_docente}', '{}');
"#,
            nombre_docente.sanitizar()
        );

        let acepta_critica = self.acepta_critica.unwrap_or_default();
        let asistencia = self.asistencia.unwrap_or_default();
        let buen_trato = self.buen_trato.unwrap_or_default();
        let claridad = self.claridad.unwrap_or_default();
        let clase_organizada = self.clase_organizada.unwrap_or_default();
        let cumple_horarios = self.cumple_horarios.unwrap_or_default();
        let fomenta_participacion = self.fomenta_participacion.unwrap_or_default();
        let panorama_amplio = self.panorama_amplio.unwrap_or_default();
        let responde_mails = self.responde_mails.unwrap_or_default();

        let mut calificaciones = String::new();

        for _ in 0..self.respuestas {
            calificaciones.push_str(&format!(
                r#"
INSERT INTO Calificacion(codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
VALUES ('{}', {}, {}, {}, {}, {}, {}, {}, {}, {});
"#,
            codigo_docente,
            acepta_critica,
            asistencia,
            buen_trato,
            claridad,
            clase_organizada,
            cumple_horarios,
            fomenta_participacion,
            panorama_amplio,
            responde_mails,
        ));
        }

        docente + &calificaciones
    }
}
