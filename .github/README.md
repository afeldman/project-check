# GitHub Actions Configuration

Dieses Projekt nutzt GitHub Actions für CI/CD. Die folgenden Workflows sind konfiguriert:

## Workflows

### 1. CI Pipeline (`.github/workflows/ci.yml`)

**Trigger:**

- Push auf `main` oder `develop` Branch
- Pull Requests auf `main` oder `develop`
- Manuell über `workflow_dispatch`

**Jobs:**

- **lint**: Code-Qualitätsprüfung
  - ruff linting
  - ruff formatting check
  - mypy type checking
- **test**: Test-Suite auf Python 3.11 und 3.12
  - Unit-Tests mit pytest
  - BDD-Tests mit behave
  - Coverage-Reports
- **build**: Distribution-Pakete erstellen
  - Wheel (`.whl`)
  - Source Distribution (`.tar.gz`)
- **docs**: Sphinx-Dokumentation bauen

### 2. Release (`.github/workflows/release.yml`)

**Trigger:**

- Push eines Tags mit Pattern `v*.*.*` (z.B. `v0.1.0`)
- Manuell über `workflow_dispatch`

**Jobs:**

- **release**: Build und Release
  - Tests ausführen
  - Distribution bauen
  - GitHub Release erstellen mit Artifacts
  - Optional: Zu PyPI publishen

### 3. Documentation (`.github/workflows/docs.yml`)

**Trigger:**

- Push auf `main` Branch
- Manuell über `workflow_dispatch`

**Jobs:**

- **docs**: Sphinx-Dokumentation zu GitHub Pages deployen
  - Dokumentation bauen
  - Zu GitHub Pages deployen

### 4. Dependency Update (`.github/workflows/dependency-update.yml`)

**Trigger:**

- Scheduled: Jeden Montag um 02:00 UTC
- Manuell über `workflow_dispatch`

**Jobs:**

- **update-dependencies**: Automatisches Dependency-Update
  - `uv.lock` aktualisieren
  - Tests ausführen
  - Pull Request erstellen (bei Erfolg)
  - Issue erstellen (bei Fehler)

## Erforderliche GitHub Secrets

Die folgenden Secrets müssen in den GitHub Repository Settings konfiguriert werden:

### Settings → Secrets and variables → Actions → New repository secret

| Secret Name     | Beschreibung                          | Verwendung     | Erforderlich |
| --------------- | ------------------------------------- | -------------- | ------------ |
| `CODECOV_TOKEN` | Token für Codecov-Integration         | ci.yml         | Optional     |
| `PYPI_TOKEN`    | API Token für PyPI Publishing         | release.yml    | Optional     |
| `GITHUB_TOKEN`  | Automatisch von GitHub bereitgestellt | Alle Workflows | Automatisch  |

### Secret-Details

#### CODECOV_TOKEN

**Zweck:** Upload von Test-Coverage-Reports zu Codecov

**Erstellen:**

1. Gehe zu [codecov.io](https://codecov.io/)
2. Verbinde dein GitHub-Repository
3. Kopiere den Upload-Token
4. Füge als Secret hinzu: `CODECOV_TOKEN`

**Optional:** Wenn nicht gesetzt, läuft der Workflow trotzdem durch (coverage wird übersprungen)

#### PYPI_TOKEN

**Zweck:** Automatisches Publishing zu PyPI bei Releases

**Erstellen:**

1. Gehe zu [PyPI Account Settings](https://pypi.org/manage/account/)
2. Navigiere zu "API tokens"
3. Erstelle neuen Token mit Scope "Entire account" oder spezifischem Projekt
4. Kopiere den Token (beginnt mit `pypi-`)
5. Füge als Secret hinzu: `PYPI_TOKEN`

**Optional:** Wenn nicht gesetzt, wird der PyPI-Publishing-Step übersprungen

#### GITHUB_TOKEN

**Zweck:** Authentifizierung für GitHub API-Operationen

**Automatisch:** Wird von GitHub automatisch für jeden Workflow bereitgestellt, keine Konfiguration nötig

## Permissions

Die Workflows benötigen folgende Permissions (bereits in den YAML-Dateien konfiguriert):

**ci.yml:**

- Keine speziellen Permissions (default: read)

**release.yml:**

```yaml
permissions:
  contents: write # Für GitHub Releases
```

**docs.yml:**

```yaml
permissions:
  contents: write # Für Pages
  pages: write # Für Pages-Deployment
  id-token: write # Für Pages-Authentifizierung
```

**dependency-update.yml:**

```yaml
permissions:
  contents: write # Für Commits
  pull-requests: write # Für PR-Erstellung
```

## Setup GitHub Pages

Für die Dokumentation muss GitHub Pages aktiviert werden:

1. Gehe zu **Settings** → **Pages**
2. **Source**: GitHub Actions
3. **Branch**: (wird automatisch von Workflow gesetzt)
4. **Custom domain**: Optional

Nach dem ersten erfolgreichen Workflow-Run ist die Dokumentation verfügbar unter:

```
https://<username>.github.io/lqt-lynter/
```

## Workflow-Status Badges

Füge diese Badges in dein README ein:

```markdown
![CI Pipeline](https://github.com/lynqtech/lqt-lynter/workflows/CI%20Pipeline/badge.svg)
![Release](https://github.com/lynqtech/lqt-lynter/workflows/Release/badge.svg)
![Documentation](https://github.com/lynqtech/lqt-lynter/workflows/Publish%20Documentation/badge.svg)
[![codecov](https://codecov.io/gh/lynqtech/lqt-lynter/branch/main/graph/badge.svg)](https://codecov.io/gh/lynqtech/lqt-lynter)
```

## Manuelles Triggern

Alle Workflows können manuell getriggert werden:

1. Gehe zu **Actions** Tab
2. Wähle den gewünschten Workflow
3. Klicke auf **Run workflow**
4. Wähle Branch
5. Klicke **Run workflow**

## Release-Prozess

Um ein neues Release zu erstellen:

```bash
# 1. Version in pyproject.toml aktualisieren
# 2. Commit und push
git add pyproject.toml
git commit -m "chore: bump version to 0.2.0"
git push

# 3. Tag erstellen und pushen
git tag v0.2.0
git push origin v0.2.0
```

Der Release-Workflow wird automatisch getriggert und:

- Tests ausführen
- Distribution bauen
- GitHub Release mit Artifacts erstellen
- (Optional) Zu PyPI publishen

## Troubleshooting

### Workflow schlägt fehl

1. **Linting-Fehler**: Führe lokal `task lint:all` aus
2. **Test-Fehler**: Führe lokal `task test` aus
3. **Build-Fehler**: Führe lokal `task build` aus

### Secrets funktionieren nicht

1. Überprüfe Secret-Namen (case-sensitive!)
2. Stelle sicher, dass Secrets auf Repository-Ebene (nicht Environment) gesetzt sind
3. Re-run des Workflows nach Secret-Änderungen

### GitHub Pages 404

1. Überprüfe, ob Pages in Settings aktiviert ist
2. Warte ~5 Minuten nach dem ersten Deploy
3. Stelle sicher, dass `docs.yml` erfolgreich durchlief

## Lokale Simulation

GitHub Actions können lokal mit [act](https://github.com/nektos/act) getestet werden:

```bash
# act installieren (macOS)
brew install act

# Workflow lokal ausführen
act -j lint
act -j test
act -j build
```

**Hinweis:** Secrets müssen in `.secrets` Datei bereitgestellt werden.
