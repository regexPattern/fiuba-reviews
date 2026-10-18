import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";
import { parseMetadata, parseOfertas } from "./index.js";

describe("parseOfertas", () => {
  it("genera una oferta por período lectivo", () => {
    const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: ÁLGEBRA LINEAL (CB002)
Período lectivo: 2026 - 2do Cuatrimestre
Actividad: FÍSICA I (CB003)`;

    const ofertas = parseOfertas(texto);
    const [oferta1, oferta2] = ofertas;

    expect(ofertas).toHaveLength(2);
    expect(oferta1).toMatchObject({
      carrera: "Ingeniería en Informática",
      anio: 2026,
      cuatrimestre: 1
    });
    expect(oferta2).toMatchObject({
      carrera: "Ingeniería en Informática",
      anio: 2026,
      cuatrimestre: 2
    });
  });

  it("agrupa materias y comisiones", () => {
    const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: ÁLGEBRA LINEAL (CB002)
Comisión: CURSO: 01A
Comisión: CURSO: 01B
Actividad: ANÁLISIS MATEMÁTICO II (CB001)
Comisión: CURSO: 1`;

    const [oferta] = parseOfertas(texto);

    expect(oferta?.materias.map((m) => [m.codigo, m.nombre])).toEqual([
      ["CB002", "ÁLGEBRA LINEAL"],
      ["CB001", "ANÁLISIS MATEMÁTICO II"]
    ]);
    expect(oferta?.materias[0]?.comisiones.map((c) => c.nombre)).toEqual([
      "CURSO: 01A",
      "CURSO: 01B"
    ]);
    expect(oferta?.materias[1]?.comisiones.map((c) => c.nombre)).toEqual(["CURSO: 1"]);
  });

  it("parsea docentes con su rol", () => {
    const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: ÁLGEBRA LINEAL (CB002)
Comisión: CURSO: 01A
Docentes: VARGAS GABRIELA (Profesor/a Adjunto/a), SOSA MARINA (Jefe/a Trabajos Practicos)`;

    const [oferta] = parseOfertas(texto);

    expect(oferta?.materias[0]?.comisiones[0]?.docentes).toEqual([
      { nombre: "VARGAS GABRIELA", rol: "Profesor/a Adjunto/a" },
      { nombre: "SOSA MARINA", rol: "Jefe/a Trabajos Practicos" }
    ]);
  });

  it("parsea 'Sin docentes' como lista vacía", () => {
    const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: ÁLGEBRA LINEAL (CB002)
Comisión: CURSO: 01B
Docentes: Sin docentes`;

    const [oferta] = parseOfertas(texto);

    expect(oferta?.materias[0]?.comisiones[0]?.docentes).toEqual([]);
  });

  it("parsea las clases y descarta las sin definir", () => {
    const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: ÁLGEBRA LINEAL (CB002)
Comisión: CURSO: 01A
Tipo de clase\tDía\tHorario\tAula
Teórica\tLunes\t09:00 a 11:00\t403 - PC
Práctica\tMiércoles\t11:00 a 13:00\t403 - PC
Comisión: CURSO: 01B
Tipo de clase\tDía\tHorario\tAula
Sin definir \tSin definir\tSin definir \tAula a determinar`;

    const [oferta] = parseOfertas(texto);

    expect(oferta?.materias[0]?.comisiones[0]?.clases).toEqual([
      { tipo: "Teórica", dia: "Lunes", inicio: "09:00", fin: "11:00", aula: "403 - PC" },
      { tipo: "Práctica", dia: "Miércoles", inicio: "11:00", fin: "13:00", aula: "403 - PC" }
    ]);
    expect(oferta?.materias[0]?.comisiones[1]?.clases).toEqual([]);
  });

  it("deja el aula en null si está vacía", () => {
    const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: ANÁLISIS MATEMÁTICO II (CB001)
Comisión: CURSO: 1
Tipo de clase\tDía\tHorario\tAula
Teórico Práctica Obligatoria\tSábado\t08:00 a 12:00\t`;

    const [oferta] = parseOfertas(texto);

    expect(oferta?.materias[0]?.comisiones[0]?.clases[0]?.aula).toBeNull();
  });

  it("soporta el formato con los datos de la comisión en una sola línea", () => {
    const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: ÁLGEBRA LINEAL (CB002)
Comisión: CONDICIONALES
\tTurno: Sin definir Cupo / insc.: 62 / 64 \tInstancias: Promoción Ubicación: Sede Unica Docentes: VARGAS GABRIELA (Profesor/a Adjunto/a), GRYNBERG SEBASTIAN (Profesor/a Adjunto/a)
Tipo de clase\tDía\tHorario\tAula
Teórica \tLunes\t09:00 a 11:00 \t403 - PC`;

    const [oferta] = parseOfertas(texto);

    expect(oferta?.materias[0]?.comisiones[0]).toEqual({
      nombre: "CONDICIONALES",
      docentes: [
        { nombre: "VARGAS GABRIELA", rol: "Profesor/a Adjunto/a" },
        { nombre: "GRYNBERG SEBASTIAN", rol: "Profesor/a Adjunto/a" }
      ],
      clases: [{ tipo: "Teórica", dia: "Lunes", inicio: "09:00", fin: "11:00", aula: "403 - PC" }]
    });
  });

  it("falla si el texto no es del SIU", () => {
    const texto = "hola mundo";

    expect(() => parseOfertas(texto)).toThrow();
  });

  it("falla si el texto no contiene períodos", () => {
    const texto = ` 0
 CARLOS CASTILLO
Propuesta: Ingeniería en Informática
Oferta de comisiones
Período lectivo:
Todos
Actividad:
Todas`;

    expect(() => parseOfertas(texto)).toThrow(/período/i);
  });

  describe("casos borde", () => {
    it("ignora los filtros del encabezado que se parecen a períodos y actividades", () => {
      const texto = ` 0
 CARLOS CASTILLO
Propuesta: Ingeniería en Informática
Oferta de comisiones
Período lectivo:
Todos
Actividad:
Todas
Período lectivo: 2026 - 1er Cuatrimestre`;

      const [oferta] = parseOfertas(texto);

      expect(oferta?.materias).toEqual([]);
    });

    it("docente sin rol en el medio de la lista", () => {
      const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: FÍSICA (CB003)
Comisión: CURSO 1
Docentes: A B (Profesor/a Adjunto/a), C D, E F (Ayudante 1ro/a)`;

      const [oferta] = parseOfertas(texto);

      expect(oferta?.materias[0]?.comisiones[0]?.docentes).toEqual([
        { nombre: "A B", rol: "Profesor/a Adjunto/a" },
        { nombre: "C D", rol: null },
        { nombre: "E F", rol: "Ayudante 1ro/a" }
      ]);
    });

    it("las observaciones no pisan a los docentes", () => {
      const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: FÍSICA (CB003)
Comisión: CURSO 1
Docentes: A B (Profesor/a Adjunto/a)
Observaciones: consultar Docentes: otro (Jefe/a)
Tipo de clase\tDía\tHorario\tAula
Teórica\tLunes\t09:00 a 11:00\t403`;

      const [oferta] = parseOfertas(texto);

      expect(oferta?.materias[0]?.comisiones[0]?.docentes).toEqual([
        { nombre: "A B", rol: "Profesor/a Adjunto/a" }
      ]);
    });

    it("descarta clases con día definido pero horario sin definir", () => {
      const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: FÍSICA (CB003)
Comisión: CURSO 1
Tipo de clase\tDía\tHorario\tAula
Teórica\tLunes\tSin definir\tAula a determinar
Práctica\tMartes\t09:00 a 11:00\tAula a determinar`;

      const [oferta] = parseOfertas(texto);

      expect(oferta?.materias[0]?.comisiones[0]?.clases).toEqual([
        {
          tipo: "Práctica",
          dia: "Martes",
          inicio: "09:00",
          fin: "11:00",
          aula: "Aula a determinar"
        }
      ]);
    });

    it("clase sin la columna de aula", () => {
      const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: FÍSICA (CB003)
Comisión: CURSO 1
Tipo de clase\tDía\tHorario
Teórica\tLunes\t09:00 a 11:00`;

      const [oferta] = parseOfertas(texto);

      expect(oferta?.materias[0]?.comisiones[0]?.clases[0]?.aula).toBeNull();
    });

    it("nombre de materia con paréntesis", () => {
      const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: INGLÉS (NIVEL 1) (CB010)
Comisión: CURSO 1
Docentes: Sin docentes`;

      const [oferta] = parseOfertas(texto);

      expect(oferta?.materias[0]).toMatchObject({ nombre: "INGLÉS (NIVEL 1)", codigo: "CB010" });
    });

    it("comisión sin tabla de clases", () => {
      const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: FÍSICA (CB003)
Comisión: CURSO 1
Docentes: A B (Profesor/a Adjunto/a)`;

      const [oferta] = parseOfertas(texto);

      expect(oferta?.materias[0]?.comisiones[0]?.clases).toEqual([]);
    });

    it("materia sin comisiones", () => {
      const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: FÍSICA (CB003)
Actividad: QUÍMICA (CB004)
Comisión: CURSO 1
Docentes: Sin docentes`;

      const [oferta] = parseOfertas(texto);

      expect(oferta?.materias.map((m) => [m.codigo, m.comisiones.length])).toEqual([
        ["CB003", 0],
        ["CB004", 1]
      ]);
    });
  });

  describe("casos particulares en ofertas reales", () => {
    it("las comisiones CONDICIONALES sin docentes quedan en su materia", () => {
      const texto = readFileSync(
        new URL("./fixtures/petroleo-2C-2025-y-1C-2026.txt", import.meta.url),
        "utf8"
      );

      const ofertas = parseOfertas(texto);
      const codigosPorPeriodo = ofertas.map((oferta) =>
        oferta.materias
          .filter((m) =>
            m.comisiones.some((c) => c.nombre === "CONDICIONALES" && c.docentes.length === 0)
          )
          .map((m) => m.codigo)
      );

      expect(codigosPorPeriodo).toEqual([
        ["CB102", "TB061", "TA040"],
        ["TB015", "TB061", "TA040"]
      ]);
    });

    it("las comisiones CONDICIONALES sin docentes ni horario se conservan", () => {
      const texto = readFileSync(
        new URL("./fixtures/mecanica-1C-2026.txt", import.meta.url),
        "utf8"
      );

      const [oferta] = parseOfertas(texto);
      const materia = oferta?.materias.find((m) => m.codigo === "TB123");
      const condicionales = materia?.comisiones.find((c) => c.nombre === "CONDICIONALES");

      expect(condicionales).toEqual({ nombre: "CONDICIONALES", docentes: [], clases: [] });
    });

    it("comisión con observaciones conserva docentes y clases", () => {
      const texto = readFileSync(new URL("./fixtures/civil-1C-2026.txt", import.meta.url), "utf8");

      const [oferta] = parseOfertas(texto);
      const comision = oferta?.materias
        .flatMap((m) => m.comisiones)
        .find((c) => c.nombre === "Koutsovitis");

      expect(comision?.docentes.map((d) => d.nombre)).toEqual([
        "KOUTSOVITIS MARIA EVA",
        "GOYENECHE MATIAS"
      ]);
      expect(comision?.clases).toEqual([
        {
          tipo: "Teórico Práctica Obligatoria",
          dia: "Lunes",
          inicio: "15:00",
          fin: "18:00",
          aula: "Aula a determinar"
        },
        {
          tipo: "Teórico Práctica Obligatoria",
          dia: "Miércoles",
          inicio: "15:00",
          fin: "18:00",
          aula: "Aula a determinar"
        }
      ]);
    });

    it("distintos tipos de clase y días", () => {
      const textos = ["civil-1C-2026", "industrial-1C-2026"].map((nombre) =>
        readFileSync(new URL(`./fixtures/${nombre}.txt`, import.meta.url), "utf8")
      );

      const ofertas = textos.flatMap((texto) => parseOfertas(texto));
      const clases = ofertas.flatMap((o) =>
        o.materias.flatMap((m) => m.comisiones.flatMap((c) => c.clases))
      );
      const tipos = new Set(clases.map((c) => c.tipo));

      expect(tipos).toContain("Desarrollo y Consulta");
      expect(tipos).toContain("Teórico-Práctica");
      expect(tipos).toContain("Teórica Asincrónica");
      expect(clases.some((c) => c.dia === "Sábado")).toBe(true);
    });

    it("todas las clases tienen un horario válido", () => {
      const textos = [
        "civil-1C-2026",
        "industrial-1C-2026",
        "informatica-1C-2026",
        "informatica-2C-2026",
        "mecanica-1C-2026",
        "petroleo-2C-2025-y-1C-2026"
      ].map((nombre) => readFileSync(new URL(`./fixtures/${nombre}.txt`, import.meta.url), "utf8"));

      const ofertas = textos.flatMap((texto) => parseOfertas(texto));
      const clases = ofertas.flatMap((o) =>
        o.materias.flatMap((m) => m.comisiones.flatMap((c) => c.clases))
      );

      expect(clases.length).toBeGreaterThan(0);
      for (const clase of clases) {
        expect(clase.inicio).toMatch(/^\d{2}:\d{2}$/);
        expect(clase.fin).toMatch(/^\d{2}:\d{2}$/);
      }
    });
  });
});

describe("parseMetadata", () => {
  it("devuelve carrera, año y cuatrimestre de cada período", () => {
    const texto = `Propuesta: Ingeniería en Informática
Período lectivo: 2026 - 1er Cuatrimestre
Actividad: ÁLGEBRA LINEAL (CB002)
Período lectivo: 2026 - 2do Cuatrimestre
Actividad: FÍSICA I (CB003)`;

    const metadata = parseMetadata(texto);

    expect(metadata).toEqual([
      { carrera: "Ingeniería en Informática", anio: 2026, cuatrimestre: 1 },
      { carrera: "Ingeniería en Informática", anio: 2026, cuatrimestre: 2 }
    ]);
  });
});
