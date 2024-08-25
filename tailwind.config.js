/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ"], // maybe add "./**/*.go"
  theme: {
    extend: {},
  },
  plugins: [require('@tailwindcss/typography'), require("daisyui"), require('tailwind-scrollbar')],
  daisyui: {
    themes: ["aqua", "dark"]
  }
}

