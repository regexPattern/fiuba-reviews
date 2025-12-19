SELECT
    oc.codigo_carrera,
    lower(unaccent (carr.nombre)) AS nombre_carrera,
    json_build_object('numero', cuat.numero, 'anio', cuat.anio) AS cuatrimestre,
    oc.contenido
FROM
    oferta_comisiones oc
    INNER JOIN cuatrimestre cuat ON cuat.codigo = oc.codigo_cuatrimestre
    INNER JOIN carrera carr ON carr.codigo = oc.codigo_carrera
ORDER BY
    cuat.codigo DESC;

