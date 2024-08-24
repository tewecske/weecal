/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ"], // maybe add "./**/*.go"
  theme: {
    extend: {},
  },
  plugins: [require('@tailwindcss/typography'), require("daisyui")],
  daisyui: {
    themes: ["aqua", "dark"]
  }
}

