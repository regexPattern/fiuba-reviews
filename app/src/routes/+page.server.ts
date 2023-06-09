import type { PageServerLoad } from "./$types";

const COMENTARIOS = [
	{
		contenido:
			"Le pone toda la onda para que aprendas, incluso fuera del horario de clase. La teóricas las arma el, incluso te enseña a programar para hacer cálculos algebraicos. Te da clase hasta los feriados(algo bueno o malo según como lo veas), siempre cumple con los recreos y es respetuoso con el tiempo. Particularmente a mí no me pareció un tipo que explique muy bien ya que usa un lenguaje extremadamente matemático pero el hecho que sepa tanto y se dedique completamente a explicarte los temas hace que ese punto se compense. Además apunta  a ejercicios clave buscando que aprendas y puedas aprobar el parcial. Siempre dispuesto a enseñarte.",
		docente: { nombre: "Grynberg" },
		cuatrimestre: "2Q2021"
	},
	{
		contenido:
			"Un crack con todas las letras, muy buenas filminas de la clases teóricas, muy buena onda, muy bien explicando. Las clases se hacen bastante amenas con él.",
		docente: { nombre: "Essaya" },
		cuatrimestre: "2Q2020"
	},
	{
		contenido:
			"Marcelo se hace querer. Me gustó mucho cursar Física I con él. Al principio me costó llevarle el ritmo a la clase, pero con el tiempo me fui sintiendo más cómodo (explica muy rápido). Usa power points en la clase y los sube al campus, por lo que recomiendo no matarse tomando apuntes sino más bien escuchar la explicación de Marcelo.",
		docente: { nombre: "Fontana" },
		cuatrimestre: "2Q2022"
	},
	{
		contenido:
			"Buen profe Mariano, como siempre, se le entiende todo y lo baja a tierra, sin las complejidades innecesarias que agregan la mayoría de los docentes. La materia es jodida, tiene orales del TP de JOS cada algunas semanas en las últimas clases. El TP Jos, que es de a 2, la hace jodida en sí, y tiene además 2 tps individuales antes de empezar con JOS.",
		docente: { nombre: "Mendez" },
		cuatrimestre: "1Q2021"
	},
	{
		contenido:
			"Nora es lo más, sabe una banda, le re gusta la materia y tiene en claro cuales son los errores más comunes, caes re armado a rendir. La práctica con herrmann y marti es buenisima.",
		docente: { nombre: "Peralta" },
		cuatrimestre: "2Q2019"
	},
	{
		contenido:
			"que mujer del bien claudia. no vayas a otra, está es tu primera opción siempre. no te vas a arrepentir",
		docente: { nombre: "Ureña" },
		cuatrimestre: "1Q2022"
	}
];

export const load = (() => {
	return { comentarios: COMENTARIOS };
}) satisfies PageServerLoad;
