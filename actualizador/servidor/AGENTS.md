# Idea general

Tengo un proyecto en el que busco sincronizar las ofertas de las materias obtenidas del SIU en la base de datos. Una oferta de una materia es simplemente un listado de las catedras de la materia y los docenetes que la componen. Una oferta puede estar mas actualizada que otra si corresponde a un cuatrimestre mas reciente.

Los docentes de la base de datos fueron agregados antes de tener informacion sobre los docentes del SIU, por lo que se debe hacer un proceso de matching o asociacion manual en el que se vincula un docente del siu con su contraparte en la base de datos que ya esta registrado, o se decide que el docente del siu corresponde a un docente nuevo que debe ser creado.

Esta asociacion se hace de manera manual a traves de una aplicacion web, razon por la que este servidor responde peticiones http para listar y resolver (aplicar los patches).

# Estructura del repositorio

El repositorio es bastante simple y se compone de los siguientes archivos:

- main.go: Archivo principal de entrada en la que se define el flujo general de la aplicacion, que es: iniciar la conexion con la base de datos, obtener las ofertas disponibles actualmente, sincronizar las materias de la base de datos cuyos datos no estan en linea con los del siu, iniciar el servidor para resolverlas.
- oferta.go: Se define la logica para obtener las ofertas mas recientes desde el siu.
- patch.go: Se define la logica para armar los patches de las materias que requieren actualizacion (ver seccion mas adelante).
- router.go: Define las rutas del servidor http.
- sync.go: Logica de sincronizacion de materias del siu con la base de datos (por si los codigos difieren y esas cosas).
- queries/: Carpeta que contiene las queries utilizadas en el proyecto. No es necesario que leas estas queries, a menos que sea de mucha utilidad para saber que hacer determinada parte del programa, pero solo leelas si es necesario. Las queries tienen nombres descriptivos. Cuando escribas una query sql agregala aca y segui la convencion de nomenclatura para embeber queries al programa.

# Patches

La idea es que este programa se ejecute todos los cuatrimestres. Para no repetir el proceso de asociacion una y otra vez, la tabla `docente` en la base de datos tiene un campo `nombre_siu`. Cuando un docente se asocie o se cree uno nuevo, es decir, cuando un docente finalmente se resuelva, va tener su `nombre_siu` marcado. Se asume que un docente que tiene esta campo con un valor ya esta resuelto. 

Los docentes tambien tienen un campo `nombre`, que es el "display name" del docente.

Las catedras no son mas que agrupaciones de docentes. Es decir, si los docentes "Castillo Carlos", "Fernandez Edmundo" y "Cardona Irene" tienen un catedra juntos, su catedra va a ser la agrupacion de estos docentes asociada a la materia que imparten juntos. Una catedra puede mantenerse de cuatrimestre a cuatrimestre. Para identificar una catedra, ademas de su id, que solo sirve para razones de indexado, lo que se puede hacer es concatenar los nombres de sus docentes en orden alfabetico, separados por guion, la concatenacion puede ser por nombre de display o nombre siu, a como veamos conveniente.

Una catedra se considera resuelta cuando todos sus docentes estan resueltos.

Un docente no esta registrado para otra materia. Se asume que la combinacion docente materia es unica. Aunque haya un docente con el mismo nombre en otra materia, se toma esta combinacion como dos docentes distintos si estan en dos materias distintas.

# Herramientas

Tenes acceso a la herramienta necesaria para hacer lecturas en la base de datos.

# Pruebas

No tenes necesidad de ejecutar pruebas. Si queres simplemente podes intentar hacer el build corriendo el comando `make build`. No es necesario ejecutar el servidor ni nada por el estilo, las pruebas las voy a ejecutar yo manualmente.
