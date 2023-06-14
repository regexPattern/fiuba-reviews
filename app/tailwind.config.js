const colors = require("tailwindcss/colors");

/** @type {import('tailwindcss').Config} */
export default {
	content: ["./src/**/*.{html,js,svelte,ts}"],
	theme: {
		colors: {
			transparent: "transparent",
			current: "currentColor",
			black: colors.black,
			white: colors.white,
			gray: colors.slate,
			blue: colors.cyan
		},
		fontFamily: {
			sans: ["Inter, sans-serif"]
		},
		extend: {}
	},
	plugins: []
};
