import db from "$lib/db";
import { materia } from "$lib/db/schema";

import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
  const materias = await db.select().from(materia);

  return { materias };
};
