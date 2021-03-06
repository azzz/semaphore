package conditions

import (
	"testing"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/specs"
)

func TestNewEvaluableExpression(t *testing.T) {
	type test struct {
		raw    string
		params map[string]*specs.Property
	}

	tests := []test{
		{
			raw: "{{ id }} == 1",
			params: map[string]*specs.Property{
				"id": {
					Template: specs.Template{},
				},
			},
		},
		{
			raw: "{{ input:id }} == {{ input:id }}",
			params: map[string]*specs.Property{
				"input:id": {
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "id",
						},
					},
				},
			},
		},
		{
			raw: "({{ input:id }} == {{ input:id }}) || {{ input:name }}",
			params: map[string]*specs.Property{
				"input:id": {
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "id",
						},
					},
				},
				"input:name": {
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "name",
						},
					},
				},
			},
		},
		{
			raw: "({{ resource:id }} == {{ input:id }})",
			params: map[string]*specs.Property{
				"input:id": {
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "id",
						},
					},
				},
				"resource:id": {
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Resource: "resource",
							Path:     "id",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.raw, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			condition, err := NewEvaluableExpression(ctx, test.raw)
			if err != nil {
				t.Fatal(err)
			}

			if len(test.params) != len(condition.Params.Params) {
				t.Fatalf("expected number of params %d, actual %d", len(test.params), len(condition.Params.Params))
			}

			for key, param := range condition.Params.Params {
				expected, has := test.params[key]
				if !has {
					t.Fatalf("unexpected result, expected %s to be set", key)
				}

				if param.Label != expected.Label {
					t.Fatalf("unexpected label %q, expected %q", param.Label, expected.Label)
				}
				if expected.Reference != nil && param.Reference == nil {
					t.Fatalf("unexpected reference %s, reference not set", key)
				}

				if expected.Reference != nil {
					if param.Reference.Resource != expected.Reference.Resource {
						t.Fatalf("unexpected resource '%+v', expected '%+v'", param.Reference.Resource, expected.Reference.Resource)
					}

					if param.Reference.Path != expected.Reference.Path {
						t.Fatalf("unexpected path '%+v', expected '%+v'", param.Reference.Path, expected.Reference.Path)
					}
				}
			}
		})
	}
}

func TestInvalidExpressions(t *testing.T) {
	tests := []string{
		"( {{ input:id }}",
		"== {{ input:id }}",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			_, err := NewEvaluableExpression(ctx, test)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}

func TestInvalidReference(t *testing.T) {
	tests := []string{
		"{{ input:id.. }}",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			_, err := NewEvaluableExpression(ctx, test)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}
