# Goose Tasks: project-check

## Reihenfolge

Tasks können teilweise parallel laufen:

```
A (scaffold) → B + C + D + E parallel → F (release)
```

Task A muss zuerst laufen (go.mod + cmd-Skelett).
Tasks B, C, D, E können parallel laufen sobald A fertig ist.
Task F (release) läuft zuletzt.

## Goose starten

```bash
# Task A zuerst
goose run --no-session --with-builtin developer --max-turns 60 \
  -i .goose/task_A_scaffold.md

# Dann B+C+D+E parallel
for task in task_B_llm task_C_rules task_D_agent task_E_report; do
  goose run --no-session --with-builtin developer --max-turns 60 \
    -i .goose/${task}.md > /tmp/goose_${task}.log 2>&1 &
  sleep 4
done
wait

# Build check
go build ./...

# Task F
goose run --no-session --with-builtin developer --max-turns 60 \
  -i .goose/task_F_release.md
```

## Nach jedem Goose-Lauf prüfen

```bash
cd /Users/anton.feldmann/Projects/priv/pkg/project-check
go build ./...
go vet ./...
```

## Spec
`/Users/anton.feldmann/Projects/priv/pkg/docs/superpowers/specs/2026-03-24-project-check-design.md`
