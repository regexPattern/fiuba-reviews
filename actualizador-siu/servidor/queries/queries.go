package queries

import _ "embed"

//go:embed oferta/select-ofertas-carreras.sql
var OfertasCarreras string

//go:embed sincronizacion/sincronizar-materias.sql
var SincronizarMaterias string

//go:embed sincronizacion/materias-no-registradas-en-db.sql
var MateriasNoRegistradasEnDb string

//go:embed patch/select-materias-candidatas.sql
var MateriasCandidatas string

//go:embed patch/select-docentes-pendientes.sql
var DocentesPendientes string

//go:embed patch/select-catedras-con-estado.sql
var CatedrasConEstado string

//go:embed patch/marcar-materia-sin-cambios.sql
var MarcarMateriaSinCambios string

//go:embed resolucion/select-docentes-con-estado.sql
var DocentesConEstado string

//go:embed resolucion/update-docentes-existentes.sql
var UpdateDocentes string

//go:embed resolucion/insert-docentes-nuevos.sql
var InsertDocentes string

//go:embed resolucion/upsert-catedras.sql
var UpsertCatedras string

//go:embed resolucion/update-cuatrimestre-ultima-actualizacion.sql
var UpdateCuatrimestreUltimaActualizacion string
