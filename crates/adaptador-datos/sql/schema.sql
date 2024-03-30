CREATE TABLE IF NOT EXISTS materia (
    codigo integer PRIMARY KEY,
    nombre text NOT NULL,
    codigo_equivalencia integer REFERENCES materia(codigo)
);

CREATE TABLE IF NOT EXISTS catedra (
    codigo uuid PRIMARY KEY,
    codigo_materia integer REFERENCES materia(codigo) NOT NULL
);

CREATE TABLE IF NOT EXISTS docente (
    codigo uuid PRIMARY KEY,
    nombre text NOT NULL,
    codigo_materia integer REFERENCES materia(codigo) NOT NULL,
    resumen_comentarios text,
    comentarios_ultimo_resumen int DEFAULT 0 NOT NULL,
    UNIQUE (nombre, codigo_materia)
);

CREATE TABLE IF NOT EXISTS catedra_docente (
    codigo_catedra uuid REFERENCES catedra(codigo) ON DELETE CASCADE NOT NULL,
    codigo_docente uuid REFERENCES docente(codigo) NOT NULL,
    PRIMARY KEY (codigo_catedra, codigo_docente)
);

CREATE TABLE IF NOT EXISTS calificacion (
    codigo uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    codigo_docente uuid REFERENCES docente(codigo) NOT NULL,
    acepta_critica double precision NOT NULL,
    asistencia double precision NOT NULL,
    buen_trato double precision NOT NULL,
    claridad double precision NOT NULL,
    clase_organizada double precision NOT NULL,
    cumple_horarios double precision NOT NULL,
    fomenta_participacion double precision NOT NULL,
    panorama_amplio double precision NOT NULL,
    responde_mails double precision NOT NULL
);

CREATE TABLE IF NOT EXISTS cuatrimestre (
    nombre text PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS comentario (
    codigo uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    codigo_docente uuid REFERENCES docente(codigo) NOT NULL,
    cuatrimestre text REFERENCES cuatrimestre(nombre) NOT NULL,
    contenido text NOT NULL,
    es_de_dolly boolean default false NOT NULL
);