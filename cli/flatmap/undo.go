package flatmap

import (
	"fmt"
	"strconv"
	"strings"
)

/*
flatmap package originally from https://github.com/nextmv-io/sdk and is licensed under the Apache 2.0 License.

Undo takes a flattened map and nests it into a multi-level map. The flattened
map should follow the [JSONPath] standard. Please see the example to understand
how the nested output looks like.

[JSONPath]: https://datatracker.ietf.org/doc/html/rfc9535
*/
func Undo(flattened map[string]any, options Options) (map[string]any, error) {
	// First, convert the flat map to a nested map. Then reshape the map into a
	// slice where appropriate.
	const magicSliceKey = "isSlice"
	nested := make(map[string]any)
	for key, value := range flattened {
		p, err := pathFrom(key, options)
		if err != nil {
			return nil, err
		}

		current := nested
		for i, k := range p {
			key := k.Key()
			if k.IsSlice() {
				current[magicSliceKey] = true
			}

			isLast := i == len(p)-1
			if isLast {
				current[key] = value
				break
			}

			if current[key] == nil {
				current[key] = make(map[string]any)
			}

			current = current[key].(map[string]any)
		}
	}

	// Convert maps to slices where appropriate using non recursive breadth
	// first search.
	queue := []map[string]any{nested}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for k, v := range current {
			m, ok := v.(map[string]any)
			if !ok {
				// Not a map, we reached the end of the tree.
				continue
			}

			if m[magicSliceKey] == nil {
				// Just a normal map, enqueue.
				queue = append(queue, m)
				continue
			}

			// A map that needs to be converted to a slice.
			delete(m, magicSliceKey)
			slice, err := toSlice(m)
			if err != nil {
				return nil, err
			}

			for _, x := range slice {
				if _, ok := x.(map[string]any); ok {
					// Enqueue all maps in the slice.
					queue = append(queue, x.(map[string]any))
				}
			}
			current[k] = slice
		}
	}
	return nested, nil
}

func toSlice(x map[string]any) ([]any, error) {
	slice := make([]any, len(x))
	for k, v := range x {
		idx, err := strconv.Atoi(k)
		if err != nil {
			return nil, err
		}

		if idx >= len(slice) || idx < 0 {
			return nil, fmt.Errorf("index %d out of bounds", idx)
		}

		slice[idx] = v
	}
	return slice, nil
}

type pathKey struct {
	name  string
	index int
}

func (p pathKey) IsSlice() bool {
	return p.index != -1
}

func (p pathKey) Key() string {
	if p.IsSlice() {
		return strconv.Itoa(p.index)
	}
	return p.name
}

type path []pathKey

func pathFrom(key string, options Options) (path, error) {
	var split []string
	if options.JSONPath {
		key = strings.TrimPrefix(key, "$")
		split = strings.Split(key[1:], ".")
	} else {
		split = strings.Split(key, ".")
	}
	p := make(path, 0, len(split))
	for _, s := range split {
		stops, err := pathKeysFrom(s)
		if err != nil {
			return path{}, err
		}

		p = append(p, stops...)
	}

	return p, nil
}

func pathKeysFrom(key string) ([]pathKey, error) {
	if strings.Contains(key, "[") {
		start := strings.Index(key, "[")
		end := strings.Index(key, "]")
		index, err := strconv.Atoi(key[start+1 : end])
		if err != nil {
			return []pathKey{}, err
		}

		return []pathKey{
			{name: key[:start], index: -1},
			{index: index},
		}, nil
	}

	return []pathKey{{name: key, index: -1}}, nil
}
