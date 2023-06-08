import type { PageServerLoad } from "./$types";

import prisma from "$lib/prisma";

const codigosComentariosEjemplos: string[] = [
  "1b10dcf2-2d8d-4522-b8cc-817c499657b2",
  "a32940d8-42fc-4075-8a0f-682b7c794f25",
  "0ca33e5e-d6e7-4497-88ff-d11a1333a51f",
  "f0500f86-9df4-4bc2-a039-b26fe683a292",
  "ecf37334-486e-4b69-bd13-5d93156557d3",
  "71af0e57-3e1c-4455-87fd-6f2533324625",
];

export const load = (async () => {
	const promises = codigosComentariosEjemplos.map(async (codigo) => {
		return prisma.comentario.findUniqueOrThrow({
			where: { codigo },
      include: {
        docente: {
          select: {
            nombre: true,
          },
        }
      },
		});
	});

	const comentariosConDocente = await Promise.all(promises);

	return { comentarios: comentariosConDocente };
}) satisfies PageServerLoad;
