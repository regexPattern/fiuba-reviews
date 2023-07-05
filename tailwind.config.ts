import type { Config } from "tailwindcss";
import { fontFamily } from "tailwindcss/defaultTheme";

const config: Config = {
	darkMode: "class",
	content: ["./src/**/*.{html,js,svelte,ts}"],
	theme: {
		extend: {
			colors: {
				foreground: "hsl(var(--foreground) / <alpha-value>)",
				background: "hsl(var(--background) / <alpha-value>)",
				border: "hsl(var(--border) / <alpha-value>)",
			},
			fontFamily: {
				sans: ["Inter", ...fontFamily.sans]
			}
		}
	},
	plugins: []
};

export default config;
