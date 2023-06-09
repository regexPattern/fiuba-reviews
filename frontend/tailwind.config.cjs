const colors = require("tailwindcss/colors");
const defaultTheme = require("tailwindcss/defaultTheme");

/** @type {import('tailwindcss').Config}*/
const config = {
	content: ["./src/**/*.{html,js,svelte,ts}"],

	theme: {
		extend: {
			colors: {
				"light": colors.slate["50"],
				"light-hover": colors.slate["100"],
				"light-border": colors.slate["200"],
				"dark": colors.slate["950"],
				"dark-hover": colors.slate["800"],
				"dark-border": colors.slate["800"],
			},
			fontFamily: {
				sans: ["Inter", ...defaultTheme.fontFamily.sans],
			},
		},
	},

	plugins: [],
};

module.exports = config;
