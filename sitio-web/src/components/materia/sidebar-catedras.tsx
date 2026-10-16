import { cacheLife, cacheTag } from "next/cache";
import { getCatedrasMateria } from "@/server/db/queries/catedras-materia";
import { ListaCatedras } from "./lista-catedras";

export async function SidebarCatedras({ codigo }: { codigo: string }) {
  "use cache";
  cacheLife("max");
  cacheTag(`materia:${codigo}`);

  const catedras = await getCatedrasMateria(codigo);

  return <ListaCatedras codigoMateria={codigo} catedras={catedras} />;
}
