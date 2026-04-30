from __future__ import annotations

import re
from pathlib import Path


TOKEN_RE = re.compile(r"{{\s*([^}]+?)\s*}}")


def render_template_string(template: str, context: dict[str, object]) -> str:
    def replace(match: re.Match[str]) -> str:
        value: object = context
        for part in match.group(1).split("."):
            if not isinstance(value, dict):
                raise KeyError(part)
            value = value[part]
        return str(value)

    return TOKEN_RE.sub(replace, template)


def extract_template_references(template: str) -> dict[str, set[str]]:
    references: dict[str, set[str]] = {"skills": set(), "bkns": set(), "agents": set()}
    for match in TOKEN_RE.finditer(template):
        parts = [part.strip() for part in match.group(1).split(".") if part.strip()]
        if len(parts) < 2:
            continue
        root, name = parts[0], parts[1]
        if root in references:
            references[root].add(name)
    return references


def extract_template_references_from_file(template_path: str | Path) -> dict[str, set[str]]:
    source = Path(template_path)
    return extract_template_references(source.read_text(encoding="utf-8"))


def render_template_file(template_path: str | Path, output_path: str | Path, context: dict[str, object]) -> Path:
    source = Path(template_path)
    target = Path(output_path)
    rendered = render_template_string(source.read_text(encoding="utf-8"), context)
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(rendered, encoding="utf-8")
    return target
