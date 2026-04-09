package agent

import (
	"fmt"
	"strings"
)

type diffOp struct {
	op   byte // ' ', '+', '-'
	line string
}

func lcsTable(a, b []string) [][]int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	return dp
}

func diffLines(a, b []string) []diffOp {
	dp := lcsTable(a, b)
	var ops []diffOp
	i, j := len(a), len(b)
	for i > 0 || j > 0 {
		switch {
		case i > 0 && j > 0 && a[i-1] == b[j-1]:
			ops = append(ops, diffOp{' ', a[i-1]})
			i--
			j--
		case j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]):
			ops = append(ops, diffOp{'+', b[j-1]})
			j--
		default:
			ops = append(ops, diffOp{'-', a[i-1]})
			i--
		}
	}
	// reverse
	for l, r := 0, len(ops)-1; l < r; l, r = l+1, r-1 {
		ops[l], ops[r] = ops[r], ops[l]
	}
	return ops
}

// unifiedDiff returns a unified-diff string between oldContent and newContent.
func unifiedDiff(path, oldContent, newContent string) string {
	if oldContent == newContent {
		return fmt.Sprintf("(no changes to %s)", path)
	}
	splitLines := func(s string) []string {
		s = strings.TrimRight(s, "\n")
		if s == "" {
			return nil
		}
		return strings.Split(s, "\n")
	}
	old := splitLines(oldContent)
	neu := splitLines(newContent)
	ops := diffLines(old, neu)

	const ctx = 3
	var b strings.Builder
	fmt.Fprintf(&b, "--- %s (original)\n", path)
	fmt.Fprintf(&b, "+++ %s (new)\n", path)

	// emit with context lines
	type hunkLine struct {
		op   byte
		line string
	}
	var hunk []hunkLine
	flush := func() {
		if len(hunk) == 0 {
			return
		}
		fmt.Fprintln(&b, "@@ ... @@")
		for _, hl := range hunk {
			fmt.Fprintf(&b, "%c %s\n", hl.op, hl.line)
		}
		hunk = hunk[:0]
	}

	pending := make([]diffOp, 0, ctx)
	for _, op := range ops {
		if op.op == ' ' {
			if len(hunk) > 0 {
				pending = append(pending, op)
				if len(pending) > ctx {
					flush()
					pending = pending[len(pending)-ctx:]
					for _, p := range pending {
						hunk = append(hunk, hunkLine{p.op, p.line})
					}
					pending = pending[:0]
				} else {
					hunk = append(hunk, hunkLine{op.op, op.line})
				}
			} else {
				pending = append(pending, op)
				if len(pending) > ctx {
					pending = pending[1:]
				}
			}
		} else {
			for _, p := range pending {
				hunk = append(hunk, hunkLine{p.op, p.line})
			}
			pending = pending[:0]
			hunk = append(hunk, hunkLine{op.op, op.line})
		}
	}
	flush()
	return b.String()
}
