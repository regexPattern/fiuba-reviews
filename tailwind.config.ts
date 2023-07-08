import type { Config } from "tailwindcss";
import { fontFamily } from "tailwindcss/defaultTheme";

const config: Config = {
	darkMode: "class",
	content: ["./src/**/*.{html,js,svelte,ts}"],
	theme: {
		container: {
			center: true,
		},
		extend: {
			colors: {
				foreground: "rgba(var(--foreground), <alpha-value>)",
				background: "rgba(var(--background), <alpha-value>)",
				border: "rgba(var(--border), <alpha-value>)",
				fiuba: "rgba(var(--fiuba), <alpha-value>)"
			},
			fontFamily: {
				sans: ["Inter", ...fontFamily.sans]
			},
			screens: {
				xs: "375px"
			}
		}
	},
	plugins: []
};

export default config;
