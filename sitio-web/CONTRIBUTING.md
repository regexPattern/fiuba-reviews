# Queries SQL Útiles

## Materia en plan vigente sin cuatrimestre de última actualización, con equivalencias que sí tienen cátedras

```sql
select
  m.codigo as codigo_materia_vigente,
  m.nombre as materia_vigente,
  e.codigo_materia_plan_anterior as codigo_materia_plan_anterior,
  count(distinct c.codigo) as cantidad_catedras_plan_anterior
from plan p
join plan_materia pm on pm.codigo_plan = p.codigo
join materia m on m.codigo = pm.codigo_materia
join equivalencia e on e.codigo_materia_plan_vigente = m.codigo
join catedra c on c.codigo_materia = e.codigo_materia_plan_anterior
where p.esta_vigente = true
  and m.cuatrimestre_ultima_actualizacion is null
group by m.codigo, m.nombre, e.codigo_materia_plan_anterior
order by m.codigo, e.codigo_materia_plan_anterior
limit 1;
```

## Materia en plan vigente sin cuatrimestre de última actualización, con equivalencias que no tienen cátedras

```sql
select
  m.codigo as codigo_materia_vigente,
  m.nombre as materia_vigente,
  e.codigo_materia_plan_anterior as codigo_materia_plan_anterior
from plan p
join plan_materia pm on pm.codigo_plan = p.codigo
join materia m on m.codigo = pm.codigo_materia
join equivalencia e on e.codigo_materia_plan_vigente = m.codigo
left join catedra c on c.codigo_materia = e.codigo_materia_plan_anterior
where p.esta_vigente = true
  and m.cuatrimestre_ultima_actualizacion is null
group by m.codigo, m.nombre, e.codigo_materia_plan_anterior
having count(c.codigo) = 0
order by m.codigo, e.codigo_materia_plan_anterior
limit 1;
```

## Materia y código de cátedra donde ningún docente tiene calificaciones

```sql
select
  m.codigo as codigo_materia,
  m.nombre as materia,
  c.codigo as codigo_catedra
from catedra c
join materia m on m.codigo = c.codigo_materia
where exists (
  select 1
  from catedra_docente cd
  where cd.codigo_catedra = c.codigo
)
and not exists (
  select 1
  from catedra_docente cd
  join calificacion_dolly cal on cal.codigo_docente = cd.codigo_docente
  where cd.codigo_catedra = c.codigo
)
order by m.codigo, c.codigo
limit 1;
```
