mod tablas;

use std::{collections::HashMap, fmt::format};

use uuid::Uuid;

use crate::remoto::{CalificacionRemoto, ComentarioRemoto};

#[derive(Debug)]
pub struct AdaptadorMateria {
    pub entidad_materia: EntidadMateria,
    pub entidades_catedras: Vec<EntidadCatedra>,
    pub entidades_docentes: Vec<EntidadDocente>,
}

#[derive(Debug)]
pub struct EntidadMateria {
    pub codigo: u32,
    pub nombre: String,
}

#[derive(Hash, Eq, PartialEq, Debug)]
pub struct EntidadCatedra {
    pub codigo: Uuid,
    pub nombre: String,
    pub codigo_materia: u32,
}

#[derive(Debug)]
pub struct EntidadDocente {
    pub codigo: Uuid,
    pub calificacion: CalificacionRemoto,
    pub codigo_catedra: Uuid,
}

#[derive(Debug)]
pub struct EntidadComentario {
    pub cuatrimestre: String,
    pub contenido: String,
    pub codigo_docente: Uuid,
}

trait SQL {
    fn insertar(&self) -> String;
}

impl SQL for EntidadMateria {
    fn insertar(&self) -> String {
        format!(
            "INSERT INTO materias (codigo, nombre) \
                VALUES ({}, '{}');",
            self.codigo,
            self.nombre.replace("'", "''")
        )
    }
}

impl SQL for EntidadCatedra {
    fn insertar(&self) -> String {
        format!(
            "INSERT INTO catedras (codigo, nombre, codigo_materia) \
            VALUES ('{}', '{}', {});",
            self.codigo,
            self.nombre.replace("'", "''"),
            self.codigo_materia,
        )
    }
}

impl SQL for EntidadDocente {
    fn insertar(&self) -> String {
        let CalificacionRemoto {
            nombre,
            respuestas,
            acepta_critica,
            asistencia,
            buen_trato,
            claridad,
            clase_organizada,
            cumple_horarios,
            fomenta_participacion,
            panorama_amplio,
            responde_mails,
        } = &self.calificacion;

        format!(
            "INSERT INTO docentes (\
                codigo, \
                nombre, \
                respuestas, \
                acepta_critica, \
                asistencia, \
                buen_trato, \
                claridad, \
                clase_organizada, \
                cumple_horarios, \
                fomenta_participacion, \
                panorama_amplio, \
                codigo_catedra\
            ) VALUES ('{}', '{}', '{}', {}, {}, {}, {}, {}, {}, {}, {}, '{}');",
            self.codigo,
            nombre.replace("'", "''"),
            respuestas,
            acepta_critica.unwrap_or_default(),
            asistencia.unwrap_or_default(),
            buen_trato.unwrap_or_default(),
            claridad.unwrap_or_default(),
            clase_organizada.unwrap_or_default(),
            cumple_horarios.unwrap_or_default(),
            fomenta_participacion.unwrap_or_default(),
            panorama_amplio.unwrap_or_default(),
            self.codigo_catedra,
        )
    }
}

impl SQL for EntidadComentario {
    fn insertar(&self) -> String {
        format!(
            "INSERT INTO comentarios (cuatrimestre, contenido, codigo_docente) \
            VALUES ('{}', '{}', '{}');",
            self.cuatrimestre, self.contenido, self.codigo_docente
        )
    }
}

pub fn exportar_query_sql(
    materias: Vec<EntidadMateria>,
    catedras: Vec<EntidadCatedra>,
    docentes: Vec<EntidadDocente>,
    comentarios: Vec<EntidadComentario>,
) -> String {
    let mut secciones = vec![
        self::tablas::TABLA_MATERIAS.to_owned(),
        self::tablas::TABLA_CATEDRAS.to_owned(),
        self::tablas::TABLA_DOCENTES.to_owned(),
        self::tablas::TABLA_COMENTARIOS.to_owned(),
    ];

    secciones.push(agrupar_queries(materias));
    secciones.push(agrupar_queries(catedras));
    secciones.push(agrupar_queries(docentes));
    secciones.push(agrupar_queries(comentarios));

    secciones.join("\n\n")
}

fn agrupar_queries<E, Q>(queries: Q) -> String
where
    E: SQL,
    Q: IntoIterator<Item = E>,
{
    queries
        .into_iter()
        .map(|elemento| elemento.insertar())
        .collect::<Vec<_>>()
        .join("\n")
}
