package agent

import (
	"strings"
	"testing"
)

func TestUnifiedDiff_NoChange(t *testing.T) {
	result := unifiedDiff("foo.go", "same\n", "same\n")
	if !strings.Contains(result, "no changes") {
		t.Errorf("expected no-changes message, got: %s", result)
	}
}

func TestUnifiedDiff_NewFile(t *testing.T) {
	result := unifiedDiff("foo.go", "", "hello\nworld\n")
	if !strings.Contains(result, "+ hello") {
		t.Errorf("expected added lines in diff, got:\n%s", result)
	}
	// "\n- " marks an actual removed line (not the "---" header)
	if strings.Contains(result, "\n- ") {
		t.Errorf("unexpected removed lines for new file, got:\n%s", result)
	}
}

func TestUnifiedDiff_DeletedLines(t *testing.T) {
	result := unifiedDiff("foo.go", "a\nb\nc\n", "a\nc\n")
	if !strings.Contains(result, "\n- ") {
		t.Errorf("expected removed line marker, got:\n%s", result)
	}
}

func TestUnifiedDiff_AddedLines(t *testing.T) {
	result := unifiedDiff("foo.go", "a\nc\n", "a\nb\nc\n")
	if !strings.Contains(result, "\n+ ") {
		t.Errorf("expected added line marker, got:\n%s", result)
	}
}

func TestUnifiedDiff_Header(t *testing.T) {
	result := unifiedDiff("my/file.go", "old\n", "new\n")
	if !strings.Contains(result, "--- my/file.go") {
		t.Error("missing --- header")
	}
	if !strings.Contains(result, "+++ my/file.go") {
		t.Error("missing +++ header")
	}
}

func TestDiffLines_Identical(t *testing.T) {
	lines := []string{"a", "b", "c"}
	ops := diffLines(lines, lines)
	for _, op := range ops {
		if op.op != ' ' {
			t.Errorf("expected all context ops, got %c for line %q", op.op, op.line)
		}
	}
}

func TestDiffLines_AllAdded(t *testing.T) {
	ops := diffLines(nil, []string{"x", "y"})
	for _, op := range ops {
		if op.op != '+' {
			t.Errorf("expected all added ops, got %c", op.op)
		}
	}
}

func TestDiffLines_AllRemoved(t *testing.T) {
	ops := diffLines([]string{"x", "y"}, nil)
	for _, op := range ops {
		if op.op != '-' {
			t.Errorf("expected all removed ops, got %c", op.op)
		}
	}
}
