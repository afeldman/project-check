package rules

import "testing"

func validRule() Rule {
	return Rule{
		ID:          "R001",
		Description: "Example rule",
		Severity:    "error",
		Languages:   []string{"go"},
	}
}

func TestValidate_OK(t *testing.T) {
	rs := RuleSet{Version: 1, Rules: []Rule{validRule()}}
	if err := Validate(rs); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidate_WrongVersion(t *testing.T) {
	rs := RuleSet{Version: 2, Rules: []Rule{validRule()}}
	if err := Validate(rs); err == nil {
		t.Error("expected error for version 2")
	}
}

func TestValidate_EmptyID(t *testing.T) {
	r := validRule()
	r.ID = ""
	rs := RuleSet{Version: 1, Rules: []Rule{r}}
	if err := Validate(rs); err == nil {
		t.Error("expected error for empty ID")
	}
}

func TestValidate_EmptyDescription(t *testing.T) {
	r := validRule()
	r.Description = ""
	rs := RuleSet{Version: 1, Rules: []Rule{r}}
	if err := Validate(rs); err == nil {
		t.Error("expected error for empty description")
	}
}

func TestValidate_InvalidSeverity(t *testing.T) {
	cases := []string{"", "critical", "high", "low"}
	for _, sev := range cases {
		r := validRule()
		r.Severity = sev
		rs := RuleSet{Version: 1, Rules: []Rule{r}}
		if err := Validate(rs); err == nil {
			t.Errorf("expected error for severity %q", sev)
		}
	}
}

func TestValidate_ValidSeverities(t *testing.T) {
	for _, sev := range []string{"error", "warning", "info"} {
		r := validRule()
		r.Severity = sev
		rs := RuleSet{Version: 1, Rules: []Rule{r}}
		if err := Validate(rs); err != nil {
			t.Errorf("unexpected error for severity %q: %v", sev, err)
		}
	}
}

func TestValidate_EmptyLanguages(t *testing.T) {
	r := validRule()
	r.Languages = nil
	rs := RuleSet{Version: 1, Rules: []Rule{r}}
	if err := Validate(rs); err == nil {
		t.Error("expected error for empty languages")
	}
}

func TestValidate_UnknownLanguage(t *testing.T) {
	r := validRule()
	r.Languages = []string{"rust"}
	rs := RuleSet{Version: 1, Rules: []Rule{r}}
	if err := Validate(rs); err == nil {
		t.Error("expected error for unknown language")
	}
}

func TestValidate_KnownLanguages(t *testing.T) {
	for lang := range LanguageExtensions {
		r := validRule()
		r.Languages = []string{lang}
		rs := RuleSet{Version: 1, Rules: []Rule{r}}
		if err := Validate(rs); err != nil {
			t.Errorf("unexpected error for language %q: %v", lang, err)
		}
	}
}

func TestValidate_MultipleRules(t *testing.T) {
	r2 := validRule()
	r2.ID = "R002"
	r2.Severity = "warning"
	r2.Languages = []string{"python"}
	rs := RuleSet{Version: 1, Rules: []Rule{validRule(), r2}}
	if err := Validate(rs); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
