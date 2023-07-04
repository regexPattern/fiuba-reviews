CREATE TABLE IF NOT EXISTS materia(
    codigo              INTEGER PRIMARY KEY,
    nombre              TEXT NOT NULL,
    codigo_equivalencia INTEGER REFERENCES materia(codigo)
);

CREATE TABLE IF NOT EXISTS catedra(
    codigo         UUID PRIMARY KEY,
    codigo_materia INTEGER REFERENCES materia(codigo) NOT NULL
);

CREATE TABLE IF NOT EXISTS docente(
    codigo      UUID PRIMARY KEY,
    nombre      TEXT NOT NULL,
    descripcion TEXT,

    -- Cantidad de comentarios del docente al momento de la ultima actualizacion de la descripcion.
    comentarios_ultima_descripcion INT DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS cuatrimestre(
    nombre TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS comentario(
    codigo         UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    codigo_docente UUID REFERENCES docente(codigo) NOT NULL,
    cuatrimestre   TEXT REFERENCES cuatrimestre(nombre) NOT NULL,
    contenido      TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS catedra_docente(
    codigo_catedra UUID REFERENCES catedra(codigo),
    codigo_docente UUID REFERENCES docente(codigo),
    CONSTRAINT catedra_docente_pkey PRIMARY KEY (codigo_catedra, codigo_docente)
);

CREATE TABLE IF NOT EXISTS calificacion(
    codigo                UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    codigo_docente        UUID REFERENCES docente(codigo) NOT NULL,
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

