# MBEL Changelog

## [1.2.1] - 2026-01-31
### Fixed
- **Module Path**: Corrected `go.mod` and all internal imports to match the official GitHub repository (`github.com/makkiattooo/MBEL`). This fixes `go get` compatibility.
- **Documentation**: Updated all 9 localized manuals with the correct installation paths.

## [1.2.0] - 2026-01-31
### Added
- **Global Documentation**: Full technical suites for 9 languages (EN, PL, DE, FR, ES, IT, RU, ZH, JA).
- **Interactive CLI**: `mbel init` wizard for easy project startup.
- **Enterprise Repository**: `Repository` interface for custom data sources (SQL, Redis).
- **AI Context**: First-class support for `@AI_Context` and `@AI_Tone`.
- **CI/CD Tools**: `mbel lint`, `fmt`, `diff`, and `watch` for professional pipelines.
- **Integration Tests**: Comprehensive test suite in `pkg/mbel/tests`.

### Fixed
- Improved flag handling (`-v`, `-h`).
- Grammar consistency in DSL examples.
- Documentation root cleanup (removed outdated Polish-only files).
