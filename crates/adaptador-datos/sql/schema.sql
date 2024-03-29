CREATE TABLE materia (
    codigo integer PRIMARY KEY,
    nombre text NOT NULL,
    codigo_equivalencia integer REFERENCES materia (codigo)
);

CREATE TABLE catedra (
    codigo uuid PRIMARY KEY,
    codigo_materia integer REFERENCES materia (codigo) NOT NULL
);

CREATE TABLE docente (
    codigo text PRIMARY KEY,
    nombre text NOT NULL,
    codigo_materia integer REFERENCES materia (codigo),
    resumen_comentarios text,
    comentarios_ultimo_resumen int DEFAULT 0 NOT NULL
);

CREATE TABLE catedra_docente (
    codigo_catedra uuid REFERENCES catedra (codigo),
    codigo_docente text REFERENCES docente (codigo),
    PRIMARY KEY (codigo_catedra, codigo_docente)
);

CREATE TABLE cuatrimestre (
    nombre text PRIMARY KEY
);

CREATE TABLE comentario (
    codigo uuid DEFAULT gen_random_uuid () PRIMARY KEY,
    codigo_docente text REFERENCES docente (codigo) NOT NULL,
    cuatrimestre text REFERENCES cuatrimestre (nombre) NOT NULL,
    contenido text NOT NULL,
    es_de_dolly boolean default false NOT NULL
);

CREATE TABLE calificacion (
    codigo uuid DEFAULT gen_random_uuid () PRIMARY KEY,
    codigo_docente text REFERENCES docente (codigo) NOT NULL,
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
