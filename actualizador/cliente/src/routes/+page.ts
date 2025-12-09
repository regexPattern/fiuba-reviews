import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
	const res = await fetch(`http://localhost:8080`);
	const json = await res.json();

	console.log(json);

	return {};
};
