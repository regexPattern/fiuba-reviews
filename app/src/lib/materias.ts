import prisma from "./prisma";

export default await prisma.materia.findMany();
