import { describe, expect, it } from "vitest"
import { publicAssetUrlFromBasePath, routerBasenameFromBasePath, viteBaseFromBasePath } from "@/config/base-path"

describe("base path configuration", () => {
  it("uses root paths by default", () => {
    expect(viteBaseFromBasePath(undefined)).toBe("/")
    expect(routerBasenameFromBasePath(undefined)).toBeUndefined()
  })

  it("normalizes a sub-path for Vite assets and React Router", () => {
    expect(viteBaseFromBasePath("agent-web-codex")).toBe("/agent-web-codex/")
    expect(viteBaseFromBasePath("/agent-web-codex/")).toBe("/agent-web-codex/")
    expect(routerBasenameFromBasePath("/agent-web-codex/")).toBe("/agent-web-codex")
  })

  it("prefixes public assets with the configured base path", () => {
    expect(publicAssetUrlFromBasePath("/", "/images/logo-icon.png")).toBe("/images/logo-icon.png")
    expect(publicAssetUrlFromBasePath("/agent-web-codex/", "/images/logo-icon.png")).toBe(
      "/agent-web-codex/images/logo-icon.png",
    )
  })
})
