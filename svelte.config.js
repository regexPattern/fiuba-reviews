import adapter from "@sveltejs/adapter-auto";
import { vitePreprocess } from "@sveltejs/kit/vite";

/** @type {import("@sveltejs/kit").Config} */
const config = {
	preprocess: [vitePreprocess({})],
	kit: {
		adapter: adapter(),
		csrf: {
			checkOrigin: false
		},
		prerender: {
			handleHttpError: "fail"
		}
	}
};

export default config;
