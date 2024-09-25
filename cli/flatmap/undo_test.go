package flatmap_test

import (
	"reflect"
	"testing"

	"github.com/mtintes/configamajig/flatmap"
)

/*
flatmap package originally from https://github.com/nextmv-io/sdk and is licensed under the Apache 2.0 License.
*/

func Test_Undo(t *testing.T) {
	type args struct {
		flattened map[string]any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "flat",
			args: args{
				flattened: map[string]any{
					"$.a": "foo",
					"$.b": 2,
					"$.c": true,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": 2,
				"c": true,
			},
		},
		{
			name: "flat with nil",
			args: args{
				flattened: map[string]any{
					"$.a": "foo",
					"$.b": nil,
					"$.c": true,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": nil,
				"c": true,
			},
		},
		{
			name: "slice",
			args: args{
				flattened: map[string]any{
					"$.a":    "foo",
					"$.b[0]": "bar",
					"$.b[1]": 2,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": []any{
					"bar",
					2,
				},
			},
		},
		{
			name: "nested map",
			args: args{
				flattened: map[string]any{
					"$.a":   "foo",
					"$.b.c": "bar",
					"$.b.d": 2,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": map[string]any{
					"c": "bar",
					"d": 2,
				},
			},
		},
		{
			name: "slice with nested maps",
			args: args{
				flattened: map[string]any{
					"$.a":      "foo",
					"$.b[0].c": "bar",
					"$.b[0].d": 2,
					"$.b[1].c": "baz",
					"$.b[1].d": 3,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": []any{
					map[string]any{
						"c": "bar",
						"d": 2,
					},
					map[string]any{
						"c": "baz",
						"d": 3,
					},
				},
			},
		},
		{
			name: "slice with nested maps with nested slice",
			args: args{
				flattened: map[string]any{
					"$.a":         "foo",
					"$.b[0].c":    "bar",
					"$.b[0].d[0]": 2,
					"$.b[0].d[1]": true,
					"$.b[1].c":    "baz",
					"$.b[1].d[0]": 3,
					"$.b[1].d[1]": false,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": []any{
					map[string]any{
						"c": "bar",
						"d": []any{
							2,
							true,
						},
					},
					map[string]any{
						"c": "baz",
						"d": []any{
							3,
							false,
						},
					},
				},
			},
		},
		{
			name: "slice with nested maps with nested slice with nested map",
			args: args{
				flattened: map[string]any{
					"$.a":           "foo",
					"$.b[0].c":      "bar",
					"$.b[0].d[0].e": 2,
					"$.b[0].d[1]":   true,
					"$.b[1].c":      "baz",
					"$.b[1].d[0].e": 3,
					"$.b[1].d[1]":   false,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": []any{
					map[string]any{
						"c": "bar",
						"d": []any{
							map[string]any{
								"e": 2,
							},
							true,
						},
					},
					map[string]any{
						"c": "baz",
						"d": []any{
							map[string]any{
								"e": 3,
							},
							false,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := flatmap.Undo(tt.args.flattened, flatmap.Options{JSONPath: true}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nest() = %v, want %v", got, tt.want)
			}
		})
	}
}
