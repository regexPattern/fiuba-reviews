import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

export const cuatrimestres = await prisma.cuatrimestre.findMany();

export default prisma;
