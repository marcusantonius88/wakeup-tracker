import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}", "./lib/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#17211b",
        moss: "#31573f",
        mint: "#c7ead5",
        paper: "#fbfaf6",
        line: "#dfe7dc",
        coral: "#d9694f",
        amber: "#d89a35"
      }
    }
  },
  plugins: []
};

export default config;

