pub const TABLA_MATERIAS: &'static str = "\
CREATE TABLE IF NOT EXISTS materias (
    codigo INTEGER PRIMARY KEY,
    nombre TEXT NOT NULL
);";

pub const TABLA_CATEDRAS: &'static str = "\
CREATE TABLE IF NOT EXISTS catedras (
    codigo         UUID PRIMARY KEY,
    nombre         TEXT NOT NULL,
    codigo_materia INTEGER REFERENCES materias(codigo) NOT NULL
);";

pub const TABLA_DOCENTES: &'static str = "\
CREATE TABLE IF NOT EXISTS docentes (
    codigo                UUID PRIMARY KEY,
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
    responde_mails        DOUBLE PRECISION,
    codigo_catedra        UUID REFERENCES catedras(codigo) NOT NULL
);";

pub const TABLA_COMENTARIOS: &'static str = "\
CREATE TABLE IF NOT EXISTS comentarios (
    codigo         SERIAL PRIMARY KEY,
    cuatrimestre   TEXT NOT NULL,
    contenido      TEXT NOT NULL,
    codigo_docente UUID REFERENCES docentes(codigo) NOT NULL
);";
