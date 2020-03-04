package trace

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

func TestNew(t *testing.T) {
	tests := map[string][]Option{
		"unexpected error":            []Option{WithMessage("unexpected error")},
		"unexpected error: component": []Option{WithMessage("unexpected error: %s", "component")},
		"file:10 unexpected error":    []Option{WithMessage("unexpected error"), WithExpression(hcl.StaticExpr(cty.StringVal("prop"), hcl.Range{Filename: "file", Start: hcl.Pos{Line: 10}}))},
		"file:10 ":                    []Option{WithExpression(hcl.StaticExpr(cty.StringVal("prop"), hcl.Range{Filename: "file", Start: hcl.Pos{Line: 10}}))},
	}

	for expected, options := range tests {
		err := New(options...)
		if err.Error() != expected {
			t.Errorf("unexpected result %s, expected %s", err, expected)
		}
	}
}
