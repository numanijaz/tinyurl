import path from "path";
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import react from "@vitejs/plugin-react"

export default defineConfig(({mode}) => {
  return {
    plugins: [
      react(), tailwindcss(),
    ],
    build: {
      outDir: "../build",
      emptyOutDir: true,
      sourcemap: mode === "dev" ? "inline" : false,
    },
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src")
      }
    }
  }
})
