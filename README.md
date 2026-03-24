# lqt-lynter

Ein Tool zum hierarchischen Auslesen von Confluence-Seiten und Generierung von Lint-Dateien (Coding Rules).

## 📚 Documentation

**[Full documentation available](https://lynqtech.github.io/lqt-lynter/)** with detailed guides and examples.

Quick links:

- [Introduction & Architecture](docs/introduction.rst)
- [Installation Guide](docs/installation.rst)
- [CLI Usage & Examples](docs/usage.rst)
- [API Reference](docs/api.rst)
- [Testing Guide](docs/testing.rst)
- [Contributing](docs/contributing.rst)

## 📋 Table of Contents

- [Beschreibung](#beschreibung)
- [Installation](#installation)
- [Task Runner](#task-runner)
- [Schnellstart](#schnellstart)
- [Verwendung](#verwendung)
- [Pipeline](#pipeline)
- [Testing](#testing)
- [Pre-Commit Hooks](#git-pre-commit-hook)
- [Entwicklung](#entwicklung)

## Beschreibung

lqt-lynter verbindet sich mit Confluence, liest hierarchisch alle Seiten aus und extrahiert mithilfe von **any-llm (Mozilla)** **strikte Coding- und Linting-Regeln** für Programmiersprachen. Diese werden dann in einer YAML-Datei gespeichert, die als Lint-Datei für Ihr Projekt verwendet werden kann.

### Was wird extrahiert? ✅

Das Tool extrahiert **ausschließlich** code-bezogene Regeln:

- **Code-Style**: Formatierung, Naming Conventions, Indentation
- **Code-Qualität**: Komplexitätslimits, Funktionslänge, Code Smells
- **Sprach-spezifische Best Practices**: Type Hints, Error Handling, Imports
- **Linting-Regeln**: Unused Variables, Deprecated Features
- **Testing-Anforderungen**: Coverage, Test Patterns
- **Dokumentations-Standards**: Docstrings, Kommentare
- **Dependencies**: Package Management
- **Build-Einstellungen**: Compilation Settings

Unterstützte Sprachen: **Python**, **Node.js/JavaScript**, **TypeScript**, **Go** und andere

### Was wird NICHT extrahiert? ❌

Das Tool filtert automatisch heraus:

- ❌ Governance-Regeln (MFA, Authentication, User Requirements)
- ❌ Infrastruktur-Regeln (AWS, Kubernetes, Database Access)
- ❌ Security-Policies (Access Control, Permissions)
- ❌ Team-Prozesse (Jira Tickets, Approvals)
- ❌ Deployment-Regeln
- ❌ Alle nicht-code-bezogenen Anforderungen

## Installation

### Voraussetzungen

- Python 3.11+
- [uv](https://github.com/astral-sh/uv) package manager
- [task](https://taskfile.dev/) task runner (optional, aber empfohlen)

### Mit uv

```bash
# Alle Dependencies installieren
uv sync --all-extras

# Oder nur die Basis-Dependencies
uv sync
```

### Task Runner installieren

```bash
# macOS mit Homebrew
brew install go-task

# Linux
sh -c '$(curl --location https://taskfile.dev/install.sh)' -- -d -b /usr/local/bin

# Windows mit Scoop
scoop install task
```

## Task Runner

Das Projekt nutzt [Task](https://taskfile.dev/) für Automatisierung. Alle verfügbaren Tasks anzeigen:

```bash
task --list
```

### Wichtigste Tasks

**Installation:**

```bash
task install           # Basis-Dependencies
task install:all       # Mit allen Extras (dev, test, docs)
task install:dev       # Development-Dependencies
task install:test      # Test-Dependencies
task install:docs      # Documentation-Dependencies
```

**Testing:**

```bash
task test              # Alle Tests (Unit + BDD)
task test:unit         # Nur Unit-Tests
task test:bdd          # Nur BDD-Tests
task test:coverage     # Mit Coverage-Report
task test:watch        # Watch-Mode (continuous testing)
task test:fast         # Tests ohne Coverage
```

**Linting:**

```bash
task lint:all          # Alle Checks (ruff + mypy)
task lint:ruff         # Nur ruff
task lint:ruff:fix     # Auto-Fix mit ruff
task lint:format       # Code formatieren
task lint:mypy         # Type checking
```

**Documentation:**

```bash
task docs:build        # Sphinx-Dokumentation bauen
task docs:open         # Dokumentation im Browser öffnen
task docs:serve        # Live-Server mit Auto-Reload
task docs:clean        # Build-Artifacts löschen
task docs:rebuild      # Clean + Build
```

**CLI:**

```bash
task cli:help          # Haupthilfe anzeigen
task cli:confluence    # Confluence-Subcommand-Hilfe
task cli:pipeline      # Pipeline-Subcommand-Hilfe
```

**Development:**

```bash
task dev:check         # Kompletter Check (lint + test)
task dev:pre-commit    # Pre-commit Hooks ausführen
```

**Pipeline:**

```bash
task pipeline:ci       # CI-Pipeline (lint + test)
task pipeline:full     # Vollständige Pipeline (lint + test + build)
```

**Build & Cleanup:**

```bash
task build             # Distribution packages bauen
task clean             # Generated files löschen
task clean:all         # Alles löschen (inkl. venv)
```

## Schnellstart

### 1. Environment Variables einrichten

Kopiere `.env.example` zu `.env`:

```bash
cp .env.example .env
```

Bearbeite `.env` und füge deine Credentials ein:

```bash
# Confluence
CONFLUENCE_USERNAME=deine-email@domain.com
CONFLUENCE_API_TOKEN=dein-confluence-api-token

# LLM (OpenAI via any-llm)
LQT_LLM_SERVER=https://api.openai.com/v1
LQT_LLM_MODEL=gpt-4
LQT_LLM_KEY=sk-proj-...
```

### 2. Confluence API Token erstellen

1. Gehe zu [Atlassian API Tokens](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Klicke auf "Create API token"
3. Gib einen Namen ein und kopiere den Token

### 3. Page-ID finden

1. Öffne die Confluence-Seite
2. Die URL sieht so aus: `https://domain.atlassian.net/wiki/spaces/SPACE/pages/123456/Page+Title`
3. Die Zahl `123456` ist deine Page-ID

### 4. Tool ausführen

```bash
# Environment laden
source .env

# Rules extrahieren
uv run lqt-lynter confluence extract \
    --confluence-url https://your-domain.atlassian.net/wiki \
    --page-id 123456 \
    --output .company-standard.yml
```

## Verwendung

### CLI-Befehle

Das Tool bietet ein modernes CLI mit Subcommands:

**Haupthilfe:**

```bash
uv run lqt-lynter --help
```

**Confluence Subcommand:**

```bash
# Hilfe
uv run lqt-lynter confluence extract --help

# Rules extrahieren
uv run lqt-lynter confluence extract \
    --confluence-url https://your-domain.atlassian.net/wiki \
    --page-id 123456 \
    --output .company-standard.yml \
    --max-depth 10
```

**Pipeline Subcommand:**

```bash
# Hilfe
uv run lqt-lynter pipeline run --help

# Vollständige Pipeline
uv run lqt-lynter pipeline run \
    --confluence-url https://your-domain.atlassian.net/wiki \
    --page-id 123456
```

### Parameter

| Parameter           | Beschreibung            | Default                 | Env-Variable           |
| ------------------- | ----------------------- | ----------------------- | ---------------------- |
| `--confluence-url`  | Confluence URL          | -                       | -                      |
| `--page-id`         | Root Page ID            | -                       | -                      |
| `--output, -o`      | Output-Datei            | `.company-standard.yml` | -                      |
| `--max-depth`       | Traversal-Tiefe         | `10`                    | -                      |
| `--username`        | Confluence-Benutzername | -                       | `CONFLUENCE_USERNAME`  |
| `--api-token`       | Confluence API-Token    | -                       | `CONFLUENCE_API_TOKEN` |
| `--llm-server`      | LLM Server URL          | -                       | `LQT_LLM_SERVER`       |
| `--llm-model`       | LLM Model Name          | -                       | `LQT_LLM_MODEL`        |
| `--llm-key`         | LLM API-Key             | -                       | `LQT_LLM_KEY`          |
| `--skip-extraction` | Extraction überspringen | `false`                 | -                      |

## Pipeline

Die Pipeline automatisiert das komplette Setup:

### Was macht die Pipeline?

```
┌─────────────────────────────────────────────────┐
│  1. Confluence Extraktion                       │
│     └─> Liest hierarchisch Coding-Regeln       │
│                                                 │
│  2. Config-Generierung                          │
│     ├─> pyproject.toml (Python/ruff)           │
│     ├─> .eslintrc.json (JavaScript/TypeScript) │
│     ├─> tsconfig.json (TypeScript)             │
│     └─> .golangci.yml (Go)                     │
│                                                 │
│  3. Hook-Installation                           │
│     └─> Pre-commit Framework Setup             │
└─────────────────────────────────────────────────┘
```

### Pipeline ausführen

```bash
# Mit Task Runner (empfohlen)
task pipeline:full

# Oder direkt
uv run lqt-lynter pipeline run \
    --confluence-url https://lynqtech.atlassian.net/wiki \
    --page-id 339859479
```

### Pipeline-Optionen

```bash
# Nur Config-Generierung (Extraktion überspringen)
uv run lqt-lynter pipeline run \
    --confluence-url https://domain.atlassian.net/wiki \
    --page-id 123456 \
    --skip-extraction
```

## Testing

Das Projekt hat umfassende Test-Coverage mit Unit- und BDD-Tests.

### Test-Architektur

- **Unit Tests**: `pytest` für funktionale Tests (19 Tests)
- **BDD Tests**: `behave` für behavior-driven Szenarien (6 Szenarien, 27 Steps)
- **Coverage**: 29% gesamt, 82% CLI, 42% Confluence-Modul

### Tests ausführen

```bash
# Alle Tests
task test

# Nur Unit-Tests
task test:unit

# Nur BDD-Tests
task test:bdd

# Mit Coverage-Report
task test:coverage

# Watch-Mode
task test:watch

# Schnelle Tests (ohne Coverage)
task test:fast
```

### Spezifische Tests

```bash
# Einzelne Test-Datei
uv run pytest tests/test_confluence_to_rules.py -v

# Einzelne Test-Klasse
uv run pytest tests/test_confluence_to_rules.py::TestConfluenceRulesExtractor -v

# Einzelner Test
uv run pytest tests/test_cli.py::TestCLI::test_main_help -v

# Mit pdb Debugger
uv run pytest tests/ --pdb
```

### Coverage-Report

```bash
# HTML-Report generieren
task test:coverage

# Report öffnen
open htmlcov/index.html
```

### Test-Struktur

```
tests/
├── conftest.py                    # pytest fixtures
├── test_confluence_to_rules.py    # 11 Unit-Tests
├── test_cli.py                    # 8 Unit-Tests
└── features/
    ├── confluence.feature         # 3 BDD-Szenarien
    ├── cli.feature                # 3 BDD-Szenarien
    └── steps/
        └── confluence_steps.py    # 27 Step-Implementierungen
```

## Git Pre-Commit Hook

Das Projekt nutzt das [pre-commit Framework](https://pre-commit.com), das automatisch Linting-Regeln auf Code-Änderungen anwendet.

### Installation

```bash
# Dependencies installieren
uv sync

# Hooks installieren
uv run pre-commit install
uv run pre-commit install --hook-type commit-msg

# Optional: Alle Hooks einmal ausführen
uv run pre-commit run --all-files
```

### Was wird geprüft?

Der Hook führt automatisch die passenden Linter aus:

| Sprache    | Linter        | Konfiguration    |
| ---------- | ------------- | ---------------- |
| Python     | ruff          | `pyproject.toml` |
| JavaScript | eslint        | `.eslintrc.json` |
| TypeScript | eslint        | `.eslintrc.json` |
| Go         | golangci-lint | `.golangci.yml`  |
| YAML       | yamllint      | `.yamllint.yml`  |

### Verwendung

```bash
# Normal committen - Hooks laufen automatisch
git add .
git commit -m "feat: neue Funktion"

# Manuell auf allen Dateien
uv run pre-commit run --all-files

# Hook-Konfiguration aktualisieren
uv run pre-commit autoupdate

# Hooks temporär überspringen (nicht empfohlen)
git commit --no-verify
```

### Test-Beispiele

Das Projekt enthält positive und negative Beispiele zum Testen der Hooks:

```bash
# Gute Beispiele (sollten bestehen)
uv run pre-commit run --files examples/good_python.py
uv run pre-commit run --files examples/good_javascript.js
uv run pre-commit run --files examples/good_typescript.ts
uv run pre-commit run --files examples/good_go.go

# Schlechte Beispiele (sollten fehlschlagen)
uv run pre-commit run --files examples/bad_python.py
uv run pre-commit run --files examples/bad_javascript.js
uv run pre-commit run --files examples/bad_typescript.ts
uv run pre-commit run --files examples/bad_go.go
```

## Funktionsweise

1. **Verbindung**: Das Tool verbindet sich mit Confluence über die API
2. **Hierarchisches Auslesen**: Beginnend bei der angegebenen Root-Seite werden alle Kind-Seiten rekursiv durchsucht
3. **Extraktion**: Der Inhalt jeder Seite wird von HTML zu Text konvertiert
4. **KI-Analyse**: any-llm (Mozilla) analysiert den Inhalt und extrahiert **nur** strikte Coding-Regeln
5. **Merge**: Alle Regeln werden zusammengeführt und nach Sprachen kategorisiert
6. **Output**: Die Regeln werden in eine YAML-Datei geschrieben

## Ausgabeformat

Die generierte YAML-Datei folgt dieser Struktur:

```yaml
python:
  _sources:
    - "Python Best Practices"
    - "Development Guidelines"
  require_type_hints: true
  max_line_length: 100
  docstring_style: "google"
  forbid_print: true
  require_logging_library: "loguru"
  max_function_length: 50
  max_complexity: 10

nodejs:
  _sources:
    - "Node.js Guidelines"
  prefer_const: true
  no_var: true
  require_jsdoc: true
  max_line_length: 120

typescript:
  _sources:
    - "TypeScript Guidelines"
  strict_mode: true
  no_implicit_any: true
  require_explicit_return_types: true

go:
  _sources:
    - "Go Best Practices"
  require_comments: true
  max_line_length: 100
  use_alpine_images: true
  prefer_multi_stage_builds: true
```

### Sources Tracking

Jedes Sprach-Regelset enthält einen `_sources`-Key, der die Confluence-Seiten auflistet, von denen die Regeln extrahiert wurden. Dies ermöglicht Traceability und hilft bei der Wartung.

## Entwicklung

### Setup

```bash
# Development-Dependencies installieren
task install:dev

# Oder manuell
uv sync --extra dev
```

### Code Quality

```bash
# Linting
task lint:all

# Auto-Fix
task lint:ruff:fix

# Formatierung
task lint:format

# Type Checking
task lint:mypy

# Kompletter Check (Lint + Tests)
task dev:check
```

### Dokumentation

```bash
# Dokumentation bauen
task docs:build

# Im Browser öffnen
task docs:open

# Live-Server (mit Auto-Reload)
task docs:serve

# Dokumentation neu bauen
task docs:rebuild
```

### Contributing

Bitte lies [Contributing Guide](docs/contributing.rst) für Details zum Entwicklungsprozess, Code-Standards und Pull-Request-Prozess.

## Beispiel-Aufruf

```bash
# LYNQTECH Engineering Principles auslesen
source .env
uv run lqt-lynter confluence extract \
    --confluence-url https://lynqtech.atlassian.net/wiki \
    --page-id 339859479 \
    --output .company-standard.yml \
    --max-depth 10
```

**Ergebnis**: Nur strikte Coding-Regeln werden extrahiert, alle Governance/AWS/Infrastruktur-Regeln werden automatisch gefiltert.

## Technologie-Stack

- **Python**: 3.11+
- **CLI Framework**: [Typer](https://typer.tiangolo.com/)
- **LLM Integration**: [any-llm-sdk](https://github.com/mozilla/any-llm) (Mozilla)
- **Confluence API**: [atlassian-python-api](https://github.com/atlassian-api/atlassian-python-api)
- **Testing**: [pytest](https://pytest.org/) + [behave](https://behave.readthedocs.io/)
- **Linting**: [ruff](https://github.com/astral-sh/ruff) + [mypy](https://mypy.readthedocs.io/)
- **Documentation**: [Sphinx](https://www.sphinx-doc.org/) + [Read the Docs Theme](https://sphinx-rtd-theme.readthedocs.io/)
- **Task Runner**: [Task](https://taskfile.dev/)
- **Package Manager**: [uv](https://github.com/astral-sh/uv)

## Lizenz

Proprietär - LynQ Technologies

## Support & Kontakt

- **Documentation**: [https://lynqtech.github.io/lqt-lynter/](https://lynqtech.github.io/lqt-lynter/)
- **Issues**: Interne Issue-Tracker
- **Team**: LynQ DevOps Team
