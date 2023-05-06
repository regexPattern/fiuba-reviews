use uuid::Uuid;

use super::{Calificacion, Catedra, ComentariosDocentePorCuatri, Materia};

pub mod create_tables {
    pub const MATERIAS: &'static str = "\
CREATE TABLE IF NOT EXISTS materias (
    codigo INTEGER PRIMARY KEY,
    nombre TEXT NOT NULL
);";

    pub const CATEDRAS: &'static str = "\
CREATE TABLE IF NOT EXISTS catedras (
    codigo         TEXT PRIMARY KEY,
    codigo_materia INTEGER REFERENCES materias(codigo) NOT NULL
);";

    pub const DOCENTES: &'static str = "\
CREATE TABLE IF NOT EXISTS docentes (
    codigo                TEXT PRIMARY KEY,
    nombre                TEXT NOT NULL,
    respuestas            INTEGER NOT NULL,
    acepta_critica        DOUBLE PRECISION,
    asistencia            DOUBLE PRECISION,
    buen_trato            DOUBLE PRECISION,
    claridad              DOUBLE PRECISION,
    clase_organizada      DOUBLE PRECISION,
    cumple_horarios       DOUBLE PRECISION,
    fomenta_participacion DOUBLE PRECISION,
    panorama_amplio       DOUBLE PRECISION,
    responde_mails        DOUBLE PRECISION
);";

    pub const COMENTARIOS: &'static str = "\
CREATE TABLE IF NOT EXISTS comentarios (
    codigo         TEXT PRIMARY KEY,
    codigo_docente TEXT REFERENCES docentes(codigo) NOT NULL,
    cuatrimestre   TEXT NOT NULL,
    contenido      TEXT NOT NULL
);";

    pub const CATEDRA_DOCENTE: &'static str = "\
CREATE TABLE IF NOT EXISTS catedra_docente (
    codigo_catedra TEXT REFERENCES catedras(codigo),
    codigo_docente TEXT REFERENCES docentes(codigo),
    CONSTRAINT catedra_docente_pkey PRIMARY KEY (codigo_catedra, codigo_docente)
);";
}

impl Materia {
    pub fn insert_query(&self) -> String {
        format!(
            "INSERT INTO materias (codigo, nombre) VALUES ({}, '{}');",
            self.codigo,
            self.nombre.replace("'", "''")
        )
    }
}

impl Catedra {
    pub fn insert_query(codigo_catedra: Uuid, codigo_materia: u32) -> String {
        format!(
            "INSERT INTO catedras (codigo, codigo_materia) \
VALUES ('{codigo_catedra}', {codigo_materia});"
        )
    }
}

impl Calificacion {
    pub fn insert_query(&self, codigo_docente: &Uuid, nombre_docente: String) -> String {
        let acepta_critica = self.acepta_critica.unwrap_or(0.0);
        let asistencia = self.asistencia.unwrap_or(0.0);
        let buen_trato = self.buen_trato.unwrap_or(0.0);
        let claridad = self.claridad.unwrap_or(0.0);
        let clase_organizada = self.clase_organizada.unwrap_or(0.0);
        let cumple_horarios = self.cumple_horarios.unwrap_or(0.0);
        let fomenta_participacion = self.fomenta_participacion.unwrap_or(0.0);
        let panorama_amplio = self.panorama_amplio.unwrap_or(0.0);
        let responde_mails = self.responde_mails.unwrap_or(0.0);
        let respuestas = self.respuestas.unwrap_or(0.0);

        format!("INSERT INTO docentes (codigo, nombre, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails, respuestas) \
VALUES ('{codigo_docente}', '{}', {acepta_critica}, {asistencia}, {buen_trato}, {claridad}, {clase_organizada}, {cumple_horarios}, {fomenta_participacion}, {panorama_amplio}, {responde_mails}, {respuestas});", nombre_docente.replace("'", "''"))
    }
}

impl ComentariosDocentePorCuatri {
    pub fn insert_query(&self, codigo_docente: &Uuid) -> String {
        let mut buffer = vec![];

        for contenido in &self.entradas {
            buffer.push(format!(
                "INSERT INTO comentarios (codigo, codigo_docente, cuatrimestre, contenido) \
VALUES ('{}', '{codigo_docente}', '{}', '{}');",
                Uuid::new_v4(),
                self.cuatrimestre,
                contenido.replace("'", "''")
            ));
        }

        buffer.join("\n")
    }
}

pub fn catedra_docente_rel_query(codigo_catedra: &Uuid, codigo_docente: &Uuid) -> String {
    format!(
        "INSERT INTO catedra_docente (codigo_catedra, codigo_docente) \
VALUES ('{codigo_catedra}', '{codigo_docente}');"
    )
}
