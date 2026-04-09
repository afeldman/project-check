# project-check — Code Standards

## Go Code Style

- Maximum file length is 150 lines. Files exceeding this limit must be split into smaller, focused files.
- Functions must not exceed 40 lines. Long functions should be decomposed.
- No deeply nested code: maximum nesting depth is 3 levels.
- All exported functions, types, and constants must have a doc comment.
- Error values must be wrapped with context using `fmt.Errorf("...: %w", err)`.
- Do not use `os.Exit` outside of `main` or `cmd` package.

## Package Structure

- Business logic belongs in `internal/`. No logic in `cmd/` beyond flag parsing and wiring.
- Each file has a single, clearly scoped responsibility.
- No circular imports.

## Testing

- Every non-trivial function in `internal/` must have a corresponding `_test.go` file.
- Tests must not depend on external services (no real LLM calls, no network).
- Table-driven tests are preferred for functions with multiple input/output cases.
- Test file names must match the file they test (e.g. `schema_test.go` for `schema.go`).

## Error Handling

- Never silently ignore errors. Either return them or log them explicitly.
- Do not use `panic` for recoverable errors.

## Security

- File write operations must validate that the target path does not escape the working directory.
- Shell command execution must use an explicit allowlist (no arbitrary command execution).
- No secrets or tokens in source code.
