/** @type {import("prettier").Config} */
const config = {
  useTabs: false,
  tabWidth: 2,
  singleQuote: false,
  trailingComma: "none",
  printWidth: 100,
  plugins: ["@ianvs/prettier-plugin-sort-imports", "prettier-plugin-tailwindcss"],
  importOrder: ["^(react|next)(/.*)?$", "<TYPES>", "<THIRD_PARTY_MODULES>", "^@/", "^[./]"],
  tailwindStylesheet: "./src/app/globals.css"
};

export default config;
