package flatmap

import (
	"fmt"
	"reflect"
	"strings"
)

/*
flatmap package originally from https://github.com/nextmv-io/sdk and is licensed under the Apache 2.0 License.

Do takes a nested map and flattens it into a single level map. The flattening
follows the [JSONPath] standard. Please see the example to understand how the
flattened output looks like.

[JSONPath]: https://datatracker.ietf.org/doc/html/rfc9535
*/
func Do(nested map[string]any, options Options) map[string]any {
	flattened := map[string]any{}
	for childKey, childValue := range nested {
		rootKey := fmt.Sprintf("$.%s", childKey)
		setChildren(flattened, rootKey, childValue)
	}

	if !options.JSONPath {
		return prefixRemover(flattened)
	}
	return flattened
}

// setChildren is a helper function for flatten. It is invoked recursively on a
// child value. If the child is not a map or a slice, then the value is simply
// set on the flattened map. If the child is a map or a slice, then the
// function is invoked recursively on the child's values, until a
// non-map-non-slice value is hit.
func setChildren(flattened map[string]any, parentKey string, parentValue any) {
	newKey := fmt.Sprintf(".%s", parentKey)
	split := strings.Split(parentKey, "")
	if len(split) > 1 {
		firstTwo := strings.Join(split[0:2], "")
		if firstTwo == "$." {
			newKey = parentKey
		}
	}

	if reflect.TypeOf(parentValue) == nil {
		flattened[newKey] = parentValue
		return
	}

	if reflect.TypeOf(parentValue).Kind() == reflect.Map {
		children := parentValue.(map[string]any)
		for childKey, childValue := range children {
			newKey = fmt.Sprintf("%s.%s", parentKey, childKey)
			setChildren(flattened, newKey, childValue)
		}
		return
	}

	if reflect.TypeOf(parentValue).Kind() == reflect.Slice {
		children := parentValue.([]any)
		if len(children) == 0 {
			flattened[newKey] = children
			return
		}

		for childIndex, childValue := range children {
			newKey = fmt.Sprintf("%s[%v]", parentKey, childIndex)
			setChildren(flattened, newKey, childValue)
		}
		return
	}

	flattened[newKey] = parentValue
}

func prefixRemover(flattened map[string]any) map[string]any {
	removed := map[string]any{}
	for key, value := range flattened {
		newKey := strings.TrimPrefix(key, "$.")
		removed[newKey] = value
	}
	return removed
}

type Options struct {
	JSONPath bool
}
