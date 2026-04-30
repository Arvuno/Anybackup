function normalizeBasePath(value: string | undefined): string {
  const trimmed = value?.trim()

  if (!trimmed || trimmed === "/" || trimmed === "./") return "/"

  const withoutTrailingSlash = trimmed.replace(/\/+$/, "")
  return withoutTrailingSlash.startsWith("/") ? withoutTrailingSlash : `/${withoutTrailingSlash}`
}

export function viteBaseFromBasePath(value: string | undefined): string {
  const basePath = normalizeBasePath(value)
  return basePath === "/" ? "/" : `${basePath}/`
}

export function routerBasenameFromBasePath(value: string | undefined): string | undefined {
  const basePath = normalizeBasePath(value)
  return basePath === "/" ? undefined : basePath
}

export function publicAssetUrlFromBasePath(basePathValue: string | undefined, assetPath: string): string {
  const basePath = viteBaseFromBasePath(basePathValue)
  const normalizedAssetPath = assetPath.replace(/^\/+/, "")
  return `${basePath}${normalizedAssetPath}`
}
