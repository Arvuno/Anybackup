SUPPORTED_LOCALES = ("zh-CN", "en-US")
DEFAULT_LOCALE = "zh-CN"


def negotiate_locale(accept_language: str | None) -> str:
    if accept_language is None or not accept_language.strip():
        return DEFAULT_LOCALE

    candidates: list[tuple[float, str]] = []
    for item in accept_language.split(","):
        parts = item.strip().split(";")
        if not parts[0]:
            continue
        quality = 1.0
        for param in parts[1:]:
            if param.strip().startswith("q="):
                try:
                    quality = float(param.strip()[2:])
                except ValueError:
                    quality = 0.0
        candidates.append((quality, _normalize_locale(parts[0])))

    for _, locale in sorted(candidates, key=lambda value: value[0], reverse=True):
        if locale in SUPPORTED_LOCALES:
            return locale
    return DEFAULT_LOCALE


def _normalize_locale(value: str) -> str:
    normalized = value.strip().replace("_", "-").lower()
    if normalized in {"zh", "zh-cn", "zh-hans"}:
        return "zh-CN"
    if normalized in {"en", "en-us"}:
        return "en-US"
    return value
