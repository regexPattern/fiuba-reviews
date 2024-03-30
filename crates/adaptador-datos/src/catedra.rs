use std::collections::{HashMap, HashSet};

use uuid::Uuid;

use crate::{docente::Calificacion, sql::Sql};

#[derive(Debug)]
#[cfg_attr(test, derive(Clone, PartialEq))]
pub struct Catedra {
    pub codigo: Uuid,
    pub docentes: HashMap<String, Calificacion>,
}

pub fn eliminar_repetidas(catedras: &mut Vec<Catedra>) {
    catedras.sort_by(|a, b| a.docentes.len().cmp(&b.docentes.len()).reverse());

    let mut nombres_docentes_catedras_unicas = Vec::with_capacity(catedras.len());
    let mut codigos_catedras_unicas = HashSet::with_capacity(catedras.len());

    for catedra in catedras.iter() {
        let nombres_docentes: HashSet<_> = catedra.docentes.keys().collect();
        if !nombres_docentes_catedras_unicas
            .iter()
            .any(|n| nombres_docentes.is_subset(n))
        {
            nombres_docentes_catedras_unicas.push(nombres_docentes);
            codigos_catedras_unicas.insert(catedra.codigo);
        }
    }

    catedras.retain(|c| codigos_catedras_unicas.contains(&c.codigo));
}

pub fn bulk_insert(insert_tuples: &Vec<String>) -> String {
    format!(
        "INSERT INTO catedra (codigo, codigo_materia)
VALUES
\t{};",
        insert_tuples.sanitize()
    )
}

impl Catedra {
    pub fn sql(&self, codigo_materia: i16) -> String {
        let codigo = self.codigo.sanitize();
        format!("({codigo}, {codigo_materia})")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn catedra_from_docentes(docentes: HashMap<String, Calificacion>) -> Catedra {
        Catedra {
            codigo: Uuid::new_v4(),
            docentes: docentes.into(),
        }
    }

    #[test]
    fn filtrado_de_catedras_con_los_mismos_docentes() {
        let docentes: HashMap<_, _> = [
            ("Garcia".to_string(), Default::default()),
            ("Husain Cerruti".to_string(), Default::default()),
        ]
        .into();

        let c1 = catedra_from_docentes(docentes.clone());
        let c2 = catedra_from_docentes(docentes.clone());
        let c3 = catedra_from_docentes([("Sassano".to_string(), Default::default())].into());

        let mut catedras = vec![c1.clone(), c2.clone(), c3.clone()];

        eliminar_repetidas(&mut catedras);

        assert_eq!(catedras.len(), 2);
        assert!(catedras.contains(&c1) || catedras.contains(&c2));
        assert!(catedras.contains(&c3));
    }

    #[test]
    fn filtrado_de_catedras_con_docentes_completamente_superpuestos() {
        let mut docentes: HashMap<_, _> = [
            ("Garcia".to_string(), Default::default()),
            ("Husain Cerruti".to_string(), Default::default()),
        ]
        .into();

        let c1 = catedra_from_docentes(docentes.clone());

        docentes.insert("Sassano".to_string(), Default::default());

        let c2 = catedra_from_docentes(docentes.clone());

        let mut catedras = vec![c1.clone(), c2.clone()];

        eliminar_repetidas(&mut catedras);

        assert_eq!(catedras.len(), 1);
        assert!(catedras.contains(&c2));
    }
}
