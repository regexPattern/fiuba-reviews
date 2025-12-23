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

//go:embed patch/select-docentes-con-estado.sql
var DocentesConEstado string

//go:embed patch/INSERT-nuevo-docente.sql
var crearNuevoDocenteQuery string

//go:embed patch/UPDATE-asociar-docente-existente.sql
var asociarDocenteExistenteQuery string

//go:embed patch/UPDATE-desactivar-catedras-materia.sql
var desactivarCatedrasMateriaQuery string

//go:embed patch/UPSERT-catedras-resueltas.sql
var upsertCatedrasResueltasQuery string
