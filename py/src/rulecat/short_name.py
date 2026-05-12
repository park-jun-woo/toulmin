def short_name(full: str) -> str:
    hash_idx = full.find("#")
    if hash_idx >= 0:
        base = full[:hash_idx]
        spec = full[hash_idx:]
    else:
        base = full
        spec = ""
    dot_idx = base.rfind(".")
    if dot_idx >= 0:
        base = base[dot_idx + 1:]
    return base + spec
