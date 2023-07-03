import type { RequestHandler } from "./$types";

const GET = (() => {
	return new Response(JSON.stringify(""));
}) satisfies RequestHandler;
