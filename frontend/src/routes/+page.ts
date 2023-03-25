import type { PageLoad } from "./$types";

export const load = (async ({ fetch }) => {
	const res = await fetch("/api");
	const data = await res.text();
	return { data }
}) satisfies PageLoad;
