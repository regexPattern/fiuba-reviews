/** @type {import('tailwindcss').Config} */
export default {
	content: ["./src/**/*.{html,js,svelte,ts}"],
	theme: {
		fontFamily: {
			sans: ["Inter, sans-serif"]
		},
		extend: {
			colors: {
				fiuba: "#0194DB"
			}
		}
	},
	plugins: []
};
