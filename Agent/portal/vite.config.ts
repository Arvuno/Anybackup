import { defineConfig } from "vitest/config"
import { loadEnv } from "vite"
import react from "@vitejs/plugin-react"
import path from "path"
import { fileURLToPath } from "url"
import { viteBaseFromBasePath } from "./src/config/base-path"

const __dirname = path.dirname(fileURLToPath(import.meta.url))

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "")
  const authServiceProxyTarget = env.VITE_AUTH_SERVICE_PROXY_TARGET?.trim()
  const conversationServiceProxyTarget = env.VITE_CONVERSATION_SERVICE_PROXY_TARGET?.trim()

  const proxy: Record<string, { target: string; changeOrigin: boolean; secure: boolean }> = {}

  if (authServiceProxyTarget) {
    proxy["/api/auth_service"] = {
      target: authServiceProxyTarget,
      changeOrigin: true,
      secure: false,
    }
  }

  if (conversationServiceProxyTarget) {
    proxy["/api/conversation_service"] = {
      target: conversationServiceProxyTarget,
      changeOrigin: true,
      secure: false,
    }
  }

  return {
    base: viteBaseFromBasePath(env.VITE_APP_BASE_PATH),
    plugins: [react()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    server: Object.keys(proxy).length > 0 ? { proxy } : undefined,
    test: {
      environment: "jsdom",
      setupFiles: ["./src/test/rtl-wrapper-mock.ts", "./src/test/setup.ts"],
    },
  }
})
