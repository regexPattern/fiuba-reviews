import type { RequestHandler } from "./$types";

export const GET = (async ({ fetch }) => {
	const res = await fetch("http://127.0.0.1:5000");
	const data = await res.text();
	return new Response(data);
}) satisfies RequestHandler;
