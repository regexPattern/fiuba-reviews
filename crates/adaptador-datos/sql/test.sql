-- 1. BOOSTRAPING
--
-- Se descargan los datos iniciales de Dolly. La query está simplificada en
-- este ejemplo ya que las cláusulas de conflicto no son necesarias en el
-- bootstraping de la base de datos.

INSERT INTO materia (codigo, nombre, codigo_equivalencia)
VALUES (6103, 'ANALISIS MATEMATICO II A', NULL);

INSERT INTO catedra (codigo, codigo_materia)
VALUES ('193f595c-fdd4-4a15-b465-864c763c7f66', 6103);

INSERT INTO docente (codigo, nombre, codigo_materia)
VALUES ('36fe1719-fab0-4258-97be-4a6716272343', 'Fernández', 6103);

INSERT INTO catedra_docente (codigo_catedra, codigo_docente)
VALUES ('193f595c-fdd4-4a15-b465-864c763c7f66', '36fe1719-fab0-4258-97be-4a6716272343');

INSERT INTO calificacion (codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
VALUES ('36fe1719-fab0-4258-97be-4a6716272343', 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0);

INSERT INTO cuatrimestre (nombre)
VALUES ('1Q2024');

INSERT INTO comentario (codigo_docente, cuatrimestre, contenido, es_de_dolly)
VALUES ('36fe1719-fab0-4258-97be-4a6716272343', '1Q2023', 'Muy buen profe.', true);



-- 2. REGISTROS PROPIOS DE FIUBA-REVIEWS
--
-- Ya con la aplicación andando, se agregan calificaciones y comentarios a
-- los docentes. Estas nuevas entradas no están disponibles en Dolly.
--
-- Hacemos la inserción de los datos insertando explícitamente el UUID tanto de
-- la calificación como del comentario agregado, para que después podamos
-- verificar que los mismos siguen en la base de datos luego de la traída de
-- nuevos datos de Dolly.
--

INSERT INTO calificacion (codigo, codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
VALUES ('840254d0-c3f6-4a86-a73f-a0c066871abf', '36fe1719-fab0-4258-97be-4a6716272343', 2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0);

-- En el caso de los comentarios, lo que determina si se conserva en la base de
-- datos o no, es el valor de su campo `es_de_dolly`. Si el mismo es false,
-- signfica que fue ingresado a través de fiuba-reviews, por lo que es un
-- comentario que se va a conservar.

INSERT INTO comentario (codigo, codigo_docente, cuatrimestre, contenido, es_de_dolly)
VALUES ('c205e767-cd21-46b8-b989-1c05591e1df1', '36fe1719-fab0-4258-97be-4a6716272343', '1Q2024', 'Excelente, me gustó mucho cursar con él.', false);



-- 3. TRAYENDO DATOS NUEVOS DE DOLLY
--
-- Nuevamente se hace el mismo proceso de inserción que se hizo en el paso 1,
-- donde vienen los mismos datos que ya se habían descargado antes, y los datos
-- más recientes agregados a Dolly posteriores a la descarga anterior. Sin
-- embargo, la query está diseñada de tal forma que solo se inserten los datos
-- más recientes de Dolly que no están ya insertados en la base de datos.
--
-- Esta vez si se dejan las cláusalas que se encargan de esta inseción con
-- conflictos.

INSERT INTO materia (codigo, nombre, codigo_equivalencia)
VALUES
	(6103, 'ANALISIS MATEMATICO II A', NULL),
	(6107, 'MATEMATICA DISCRETA', NULL)
ON CONFLICT (codigo)
DO NOTHING;

-- Las cátedras no son más que un conjunto de docentes que imparten determinada
-- materia, entonces con tal que de que se conserven las materias y docentes
-- que ya estaban, el eliminar las cátedras no implica ningún problema, porque
-- igualmente esas agrupaciones se van a volver a establecer, puede que con
-- otro código, pero tampoco me interesa que el código de una cátedra en
-- particular no cambie, porque como digo, estas no son más que una relación de
-- docentes. Sin tras la actualización, una cátedra tienen un nuevo código pero
-- agrupa a los mismos docentes, a efectos prácticos sigue siendo la misma
-- cátedra.

DELETE FROM catedra;

INSERT INTO catedra (codigo, codigo_materia)
VALUES
	('193f595c-fdd4-4a15-b465-864c763c7f66', 6103),
	('27720503-8559-48bb-b29b-d6c750c74208', 6107);

INSERT INTO docente (codigo, nombre, codigo_materia)
VALUES
	-- Agrego nuevamente a Fernández, docente de Análisis Matemático A. En este
	-- caso tengo un nuevo código, pero como voy a verificar a continuación,
	-- esta entrada realmente no se inserta porque ya tenía un docente llamado
	-- Fernández en esta misma materia.

	(gen_random_uuid (), 'Fernández', 6103),
	(gen_random_uuid (), 'Ocampos', 6103),
	(gen_random_uuid (), 'Álvarez', 6107)
ON CONFLICT (nombre, codigo_materia)
DO NOTHING;

-- Corriendo la siguiente query se comprueba que sigue estando el docente
-- 'Fernández' y conserva código con el que se insertó en el paso 1
-- ('36fe1719-fab0-4258-97be-4a6716272343'), es decir, es el mismo docente, no
-- se reinsertó. Si se insertaron 'Ocampos' y 'Álvarez'.

SELECT codigo
FROM docente;

-- Si o si voy a necesitar la tabla de (codigo_materia, nombre_docente) ->
-- codigo_docente existente para poder asignar los docentes ya existentes a las
-- nuevas cátedras.

SELECT codigo_materia, nombre, codigo
FROM docente;

INSERT INTO catedra_docente (codigo_catedra, codigo_docente)
VALUES ();
