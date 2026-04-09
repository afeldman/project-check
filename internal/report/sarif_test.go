package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/afeldman/project-check/internal/agent"
	"github.com/afeldman/project-check/internal/rules"
)

func testRuleSet() rules.RuleSet {
	return rules.RuleSet{
		Version: 1,
		Rules: []rules.Rule{
			{ID: "R001", Description: "No globals", Severity: "error", Languages: []string{"go"}},
			{ID: "R002", Description: "Add comments", Severity: "warning", Languages: []string{"go"}},
			{ID: "R003", Description: "Use logger", Severity: "info", Languages: []string{"go"}},
		},
	}
}

func TestWriteSARIF_CreatesFile(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.sarif")
	err := WriteSARIF(out, nil, testRuleSet(), "1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(out); os.IsNotExist(err) {
		t.Fatal("SARIF file was not created")
	}
}

func TestWriteSARIF_ValidJSON(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.sarif")
	findings := []agent.Finding{
		{RuleID: "R001", File: "main.go", Line: 10, Message: "global var found"},
	}
	if err := WriteSARIF(out, findings, testRuleSet(), "0.9.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestWriteSARIF_ToolVersion(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.sarif")
	if err := WriteSARIF(out, nil, testRuleSet(), "2.3.4"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	var v map[string]interface{}
	json.Unmarshal(data, &v)
	runs := v["runs"].([]interface{})
	tool := runs[0].(map[string]interface{})["tool"].(map[string]interface{})
	driver := tool["driver"].(map[string]interface{})
	if driver["version"] != "2.3.4" {
		t.Errorf("expected version 2.3.4, got %v", driver["version"])
	}
}

func TestWriteSARIF_SeverityMapping(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.sarif")
	findings := []agent.Finding{
		{RuleID: "R001", File: "a.go", Line: 1, Message: "err"},
		{RuleID: "R002", File: "b.go", Line: 2, Message: "warn"},
		{RuleID: "R003", File: "c.go", Line: 3, Message: "info"},
	}
	if err := WriteSARIF(out, findings, testRuleSet(), "1.0.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	var v map[string]interface{}
	json.Unmarshal(data, &v)
	runs := v["runs"].([]interface{})
	results := runs[0].(map[string]interface{})["results"].([]interface{})

	want := []string{"error", "warning", "note"}
	for i, r := range results {
		level := r.(map[string]interface{})["level"].(string)
		if level != want[i] {
			t.Errorf("result %d: expected level %q, got %q", i, want[i], level)
		}
	}
}

func TestWriteSARIF_EmptyFindings(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.sarif")
	if err := WriteSARIF(out, []agent.Finding{}, testRuleSet(), "1.0.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	var v map[string]interface{}
	json.Unmarshal(data, &v)
	runs := v["runs"].([]interface{})
	results := runs[0].(map[string]interface{})["results"].([]interface{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
