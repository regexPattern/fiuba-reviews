import db from "$lib/db";
import { catedra, materia } from "$lib/db/schema";
import { ImageResponse, type ImageResponseOptions } from "@ethercorps/sveltekit-og";
import type { RequestHandler } from "@sveltejs/kit";
import { eq } from "drizzle-orm";

import OG from "./OG.svelte";

export const GET: RequestHandler = async ({ url, params }) => {
	const materiaYCatedra = await db
		.select({
			materia: {
				codigo: materia.codigo,
				nombre: materia.nombre
			},
			catedra: {
				codigo: catedra.codigo,
				nombre: catedra.nombre
			}
		})
		.from(catedra)
		.where(eq(catedra.codigo, params.codigoCatedra || ""));

	const fontFile = await fetch(`${url.origin}/fonts/Geist-Medium.otf`);
	const fontData = await fontFile.arrayBuffer();

	const imageOpts: ImageResponseOptions = {
		fonts: [
			{
				name: "Geist",
				data: fontData,
				style: "normal"
			}
		]
	};

	return new ImageResponse(OG as any, imageOpts, {
		title: "Testing this"
	});
};
