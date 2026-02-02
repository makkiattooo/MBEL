# MBEL — Modern, Deterministic Localization

MBEL is a domain-specific language and toolset for internationalization that treats translations as code. It reduces ambiguity, prevents common CI failures, and provides deterministic context for automated translation workflows (LLMs, CAT tools, and human reviewers).

Key goals:
- Reduce merge conflicts and ambiguous keys
- Provide rich, structured context for translators and LLMs
- Offer a CI-friendly toolchain (lint, fmt, compile)
- Support production-grade deployment and monitoring

Quick links
- Quick Start: [docs/en/QUICKSTART.md](docs/en/QUICKSTART.md)
- Full documentation suite (per-language): see each language "Full Suite" link below

Supported documentation languages
| Language | Quick Start | Manual | Full Suite |
| :--- | :--- | :--- | :--- |
| 🇬🇧 English | 🚀 [Quick Start](docs/en/QUICKSTART.md) | [Manual](docs/en/Manual.md) | [Full Suite](docs/en/SUITE.md) |
| 🇵🇱 Polski | [Szybki start](docs/pl/Manual.md) | [Manual](docs/pl/Manual.md) | [Full Suite](docs/pl/SUITE.md) |
| 🇩🇪 Deutsch | [Handbuch](docs/de/Manual.md) | [Manual](docs/de/Manual.md) | [Full Suite](docs/de/SUITE.md) |
| 🇫🇷 Français | [Manuel](docs/fr/Manual.md) | [Manual](docs/fr/Manual.md) | [Full Suite](docs/fr/SUITE.md) |
| 🇪🇸 Español | [Manual](docs/es/Manual.md) | [Manual](docs/es/Manual.md) | [Full Suite](docs/es/SUITE.md) |
| 🇮🇹 Italiano | [Manuale](docs/it/Manual.md) | [Manual](docs/it/Manual.md) | [Full Suite](docs/it/SUITE.md) |
| 🇷🇺 Русский | [Руководство](docs/ru/Manual.md) | [Manual](docs/ru/Manual.md) | [Full Suite](docs/ru/SUITE.md) |
| 🇨🇳 中文 | [官方手册](docs/zh/Manual.md) | [Manual](docs/zh/Manual.md) | [Full Suite](docs/zh/SUITE.md) |
| 🇯🇵 日本語 | [マニュアル](docs/ja/Manual.md) | [Manual](docs/ja/Manual.md) | [Full Suite](docs/ja/SUITE.md) |

---

Why MBEL

- Deterministic AI context: attach structured metadata (tone, constraints, examples) directly to keys so automated translation yields reliable outputs.
- Programmable logic: plurals, ranges, gender and arbitrary matchers live in the DSL, not in application code.
- CI-friendly: `mbel lint` and `mbel fmt` make translations part of your engineering workflow.
- Production-ready: optional HTML escaping, lazy-loading of large locale sets, sourcemaps for debugging, and metrics for observability.

Getting started

Install the CLI and SDK:
```bash
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
go get github.com/makkiattooo/MBEL
```

Run the quickstart guide to compile example locales and generate sourcemaps:
```bash
mbel compile examples -o examples_out.json -sourcemap
```

Where to go next

- Read the Quick Start: [docs/en/QUICKSTART.md](docs/en/QUICKSTART.md)
- Explore the Full Suite for your language (see table above)
- For CI: add `mbel lint` and `mbel fmt` to pre-merge checks

Contributing

Contributions welcome — please follow the project's CONTRIBUTING guidelines. When contributing translations, add `SUITE.md` entries under the appropriate `docs/<lang>/` folder so users can discover localized resources.

---

If you want, I will now create `SUITE.md` files for each supported language that list available docs (and point to English fallbacks where translations are missing).
