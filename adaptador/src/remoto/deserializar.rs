use std::{collections::HashMap, fmt};

use base64::Engine;
use serde::{
    de::{Error, IgnoredAny, MapAccess, SeqAccess, Visitor},
    Deserialize, Deserializer,
};
use serde_json::Value;

use crate::remoto::CalificacionRemoto;

use super::CatedraRemoto;

pub fn codigo<'de, D>(deserializer: D) -> Result<u32, D::Error>
where
    D: Deserializer<'de>,
{
    struct U32Visitor;

    impl<'de> Visitor<'de> for U32Visitor {
        type Value = u32;

        fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
            formatter.write_str("an unsigned 32-integer as a string")
        }

        fn visit_str<E>(self, number: &str) -> Result<Self::Value, E>
        where
            E: Error,
        {
            number.parse().map_err(Error::custom)
        }
    }

    deserializer.deserialize_any(U32Visitor)
}

pub fn comentarios<'de, D>(deserializer: D) -> Result<Vec<String>, D::Error>
where
    D: Deserializer<'de>,
{
    struct CommentsVisitor;

    impl<'de> Visitor<'de> for CommentsVisitor {
        type Value = Vec<String>;

        fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
            formatter.write_str("array of base64-enconded strings")
        }

        fn visit_seq<A>(
            self,
            mut comentarios_base64: A,
        ) -> std::result::Result<Self::Value, A::Error>
        where
            A: SeqAccess<'de>,
        {
            let mut comentarios = vec![];

            while let Some(elemento) = comentarios_base64.next_element::<Value>()? {
                if let Value::String(comentario) = elemento {
                    let comentario = match base64::engine::general_purpose::STANDARD.decode(comentario) {
                        Ok(comentario) => comentario,
                        _ => continue,
                    };

                    let comentario = String::from_utf8(comentario)
                        .map_err(|_| Error::custom("invalid utf8-encoded characters in comment"))?;

                    comentarios.push(comentario);
                }
            }

            Ok(comentarios)
        }
    }

    deserializer.deserialize_seq(CommentsVisitor)
}

impl<'de> Deserialize<'de> for CatedraRemoto {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        enum Field {
            Docentes,
            Ignored,
        }

        impl<'de> Deserialize<'de> for Field {
            fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
            where
                D: Deserializer<'de>,
            {
                struct FieldVisitor;

                impl<'de> Visitor<'de> for FieldVisitor {
                    type Value = Field;

                    fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
                        formatter.write_str("field of struct `Catedra`")
                    }

                    fn visit_str<E>(self, field: &str) -> Result<Self::Value, E>
                    where
                        E: Error,
                    {
                        match field {
                            "docentes" => Ok(Field::Docentes),
                            _ => Ok(Field::Ignored),
                        }
                    }
                }

                deserializer.deserialize_identifier(FieldVisitor)
            }
        }

        struct CatedraVisitor;

        impl<'de> Visitor<'de> for CatedraVisitor {
            type Value = Vec<CalificacionRemoto>;

            fn expecting(&self, formatter: &mut std::fmt::Formatter) -> std::fmt::Result {
                formatter.write_str("struct `Catedra`")
            }

            fn visit_map<A>(self, mut map: A) -> std::result::Result<Self::Value, A::Error>
            where
                A: MapAccess<'de>,
            {
                let mut docentes = None;

                while let Some(key) = map.next_key()? {
                    match key {
                        Field::Docentes => {
                            if docentes.is_some() {
                                return Err(Error::duplicate_field("docentes"));
                            }

                            docentes = Some(
                                map.next_value::<HashMap<String, CalificacionRemoto>>()?
                                    .into_iter()
                                    .map(|(_, docente)| docente)
                                    .collect::<Vec<_>>(),
                            );
                        }
                        Field::Ignored => {
                            map.next_value::<IgnoredAny>()?;
                        }
                    }
                }

                let docentes = docentes.ok_or_else(|| Error::missing_field("docentes"))?;

                if docentes.len() > 0 {
                    Ok(docentes)
                } else {
                    Err(Error::custom("catedra without docentes is invalid"))
                }
            }
        }

        const FIELDS: &'static [&'static str] = &["docentes"];
        let mut calificaciones =
            deserializer.deserialize_struct("Catedra", FIELDS, CatedraVisitor)?;

        calificaciones.sort_by(|a, b| a.nombre.cmp(&b.nombre));

        let nombre = calificaciones
            .iter()
            .map(|calificacion| calificacion.nombre.clone())
            .collect::<Vec<_>>()
            .join("-");

        Ok(Self {
            nombre,
            calificaciones,
        })
    }
}

#[cfg(test)]
mod tests {
    use crate::remoto::ComentarioRemoto;

    use super::*;

    #[test]
    fn deserializando_comentarios_en_base64() {
        let input = r#"
{
	"mat": 6103,
	"doc": "Acero",
	"cuat": "1Q2015",
	"editado": 0,
    "comentarios": [
        "QTAgZXMgbG8gbcOhcyBncmFuZGUgcXVlIGhheSBlbiBlbCBtdW5kbw==",
        "Tm8gaW5zcGlyYSBnYW5hcyBkZSBwcmVndW50YXIgZW4gY2xhc2UsIHBvciBtaWVkbyBhIHF1ZSB0ZSBjb250ZXN0ZSBtYWwgeSB0ZSByaWRpY3VsaWNlLg=="
    ]
}"#;

        let comentarios: ComentarioRemoto = serde_json::from_str(input).unwrap();

        assert_eq!(comentarios, ComentarioRemoto {
            codigo_materia: 6103,
            nombre_docente: "Acero".to_string(),
            cuatrimestre: "1Q2015".to_string(),
            contenido_comentarios: vec![
                "A0 es lo más grande que hay en el mundo".to_string(),
                "No inspira ganas de preguntar en clase, por miedo a que te conteste mal y te ridiculice.".to_string(),
            ],
        });
    }

    #[test]
    fn deserializando_comentarios_con_contenidos_nulos() {
        let input = r#"
{
	"mat": 6103,
	"doc": "Acero",
	"cuat": "1Q2015",
	"editado": 0,
    "comentarios": [
        null,
        "QTAgZXMgbG8gbcOhcyBncmFuZGUgcXVlIGhheSBlbiBlbCBtdW5kbw==",
        null,
        "Tm8gaW5zcGlyYSBnYW5hcyBkZSBwcmVndW50YXIgZW4gY2xhc2UsIHBvciBtaWVkbyBhIHF1ZSB0ZSBjb250ZXN0ZSBtYWwgeSB0ZSByaWRpY3VsaWNlLg=="
    ]
}"#;

        let comentarios: ComentarioRemoto = serde_json::from_str(input).unwrap();

        assert_eq!(comentarios.contenido_comentarios,
            vec![
                "A0 es lo más grande que hay en el mundo".to_string(),
                "No inspira ganas de preguntar en clase, por miedo a que te conteste mal y te ridiculice.".to_string(),
            ],
        );
    }

    #[test]
    fn deserializando_calificacion() {
        let input = r#"
{
    "nombre": "Suarez",
    "respuestas": 38,
    "acepta_critica": 2.92,
    "asistencia": 3.87,
    "buen_trato": 3.66,
    "claridad": 2.97,
    "clase_organizada": 3.05,
    "cumple_horarios": 3.82,
    "fomenta_participacion": 2.79,
    "panorama_amplio": 3.0,
    "responde_mails": 2.74
}"#;

        let docente: CalificacionRemoto = serde_json::from_str(input).unwrap();

        assert_eq!(
            docente,
            CalificacionRemoto {
                nombre: "Suarez".to_string(),
                respuestas: 38,
                acepta_critica: Some(2.92),
                asistencia: Some(3.87),
                buen_trato: Some(3.66),
                claridad: Some(2.97),
                clase_organizada: Some(3.05),
                cumple_horarios: Some(3.82),
                fomenta_participacion: Some(2.79),
                panorama_amplio: Some(3.0),
                responde_mails: Some(2.74),
            }
        );
    }

    #[test]
    fn deserializando_catedra() {
        let input = r#"
{
    "nombre": "...",
    "promedio": "...",
    "docentes": {
        "Suarez": {
            "nombre": "Suarez",
            "respuestas": 38
        },
        "Sanchez": {
            "nombre": "Sanchez",
            "respuestas": 22
        }
    }
}"#;

        let catedra: CatedraRemoto = serde_json::from_str(input).unwrap();

        // Los docentes tienen que ser verificados utilizando `Vec::contains` ya que al momento de
        // deserializar el hashmap de docentes, estos se agregan al vector sin un orden concreto.

        assert!(
            catedra.nombre == "Suarez-Sanchez".to_string()
                || catedra.nombre == "Sanchez-Suarez".to_string()
        );

        assert_eq!(catedra.calificaciones.len(), 2);

        assert!(catedra.calificaciones.contains(&CalificacionRemoto {
            nombre: "Suarez".to_string(),
            respuestas: 38,
            ..Default::default()
        }));

        assert!(catedra.calificaciones.contains(&CalificacionRemoto {
            nombre: "Sanchez".to_string(),
            respuestas: 22,
            ..Default::default()
        }));
    }

    #[test]
    fn intentando_deserializar_catedra_sin_llave_docentes() {
        let input = r#"{}"#;

        let err = serde_json::from_str::<CatedraRemoto>(input).unwrap_err();

        assert_eq!(
            err.to_string(),
            "missing field `docentes` at line 1 column 2"
        );
    }

    #[test]
    fn intentando_deserializar_catedra_sin_docentes() {
        let input = r#"
{
    "docentes": {}
}"#;

        let err = serde_json::from_str::<CatedraRemoto>(input).unwrap_err();

        assert_eq!(
            err.to_string(),
            "catedra without docentes is invalid at line 4 column 1"
        );
    }

    #[test]
    fn intentando_deserializar_catedra_con_llave_docentes_duplicada() {
        let input = r#"
{
    "docentes": {
        "Suarez": {
            "nombre": "Suarez",
            "respuestas": 38
        }
    },
    "docentes": {
        "Sanchez": {
            "nombre": "Sanchez",
            "respuestas": 22
        }
    }
}"#;

        let err = serde_json::from_str::<CatedraRemoto>(input).unwrap_err();

        assert_eq!(
            err.to_string(),
            "duplicate field `docentes` at line 9 column 14"
        );
    }
}
