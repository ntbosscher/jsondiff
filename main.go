package jsondiff

import (
	"encoding/json"
	"errors"
	"reflect"
)

// Diff compares oldValue with newValue and returns a json tree of
// the changed values.
func Diff(oldValue interface{}, newValue interface{}) (json.RawMessage, error) {
	return DiffFormat(oldValue, newValue, DefaultFormat)
}

func DefaultFormat(oldValue interface{}, newValue interface{}) (outputValue interface{}) {
	return newValue
}

func NewValueFormat(oldValue interface{}, newValue interface{}) (outputValue interface{}) {
	return newValue
}

func OldValueFormat(oldValue interface{}, newValue interface{}) (outputValue interface{}) {
	return oldValue
}

func BothValuesAsMapFormat(oldValue interface{}, newValue interface{}) (outputValue interface{}) {
	return map[string]interface{}{
		"Old": oldValue,
		"New": newValue,
	}
}

func DiffOldNew(oldValue interface{}, newValue interface{}) (json.RawMessage, error) {
	return DiffFormat(oldValue, newValue, BothValuesAsMapFormat)
}

// Formatter controls how to represent the diff in the output json message
// e.g. to show only the newValue, this func would return newValue
// e.g. to show only the oldValue, this func would return oldValue
// e.g. to show a {old: <v>, new: v}, this func would return map[string]interface{}{ "old": oldValue, "new": newValue }
//
// oldValue & newValue will always be non-struct types
type Formatter func(oldValue interface{}, newValue interface{}) (outputValue interface{})

func getIgnoredKeys(typ reflect.Type, baseAddr []string, maxDepth int) [][]string {
	addrs := [][]string{}

	if maxDepth == 0 {
		return addrs
	}

	if typ.Kind() == reflect.Map {
		return addrs
	}

	for i := 0; i < typ.NumField(); i++ {
		tg := typ.Field(i).Tag.Get("jsondiff")
		if tg == "-" {
			addrs = append(addrs, append(baseAddr, typ.Field(i).Name))
		}

		child := baseType(typ.Field(i).Type)
		if child.Kind() == reflect.Struct {
			add := getIgnoredKeys(child, append(baseAddr, typ.Field(i).Name), maxDepth-1)
			addrs = append(addrs, add...)
		}
	}

	return addrs
}

func baseType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
		typ = typ.Elem()
	}

	return typ
}

func DiffFormat(oldValue interface{}, newValue interface{}, formatter Formatter) (json.RawMessage, error) {
	typ := reflect.TypeOf(oldValue)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct && typ.Kind() != reflect.Map {
		return nil, errors.New("jsondiff only supports structs (and pointers to structs)")
	}

	ignoreAddrs := getIgnoredKeys(typ, nil, 10)

	jsonOld, err := json.Marshal(oldValue)
	if err != nil {
		return nil, err
	}

	jsonNew, err := json.Marshal(newValue)
	if err != nil {
		return nil, err
	}

	oldMap := map[string]interface{}{}
	if err = json.Unmarshal(jsonOld, &oldMap); err != nil {
		return nil, err
	}

	newMap := map[string]interface{}{}
	if err = json.Unmarshal(jsonNew, &newMap); err != nil {
		return nil, err
	}

	diff := map[string]interface{}{}
	calculateDiff(diff, oldMap, newMap, formatter, nil, ignoreAddrs)

	return json.Marshal(diff)
}

// calculateDiff calculates the difference between the old and new maps
// and fills diffResult with the result
func calculateDiff(
	diffResult map[string]interface{},
	oldMap map[string]interface{},
	newMap map[string]interface{},
	formatter Formatter,
	addr []string,
	ignoreAddrs [][]string,
) {

	// iterate over keys
	for _, k := range allKeys(oldMap, newMap) {
		if containsAddr(ignoreAddrs, append(addr, k)) {
			delete(diffResult, k)
			continue
		}

		newProp := newMap[k]
		oldProp := oldMap[k]

		// check if the values are maps themselves
		mpOld, oldIsMap := oldProp.(map[string]interface{})
		mpNew, newIsMap := newProp.(map[string]interface{})

		// one is a map, the other is not, must be a change
		if oldIsMap != newIsMap {
			diffResult[k] = formatter(wrapJson(oldProp, append(addr, k), ignoreAddrs), wrapJson(newProp, append(addr, k), ignoreAddrs))
			continue
		}

		// both are maps, check subkeys for changes
		if oldIsMap && newIsMap {
			subResult := map[string]interface{}{}
			calculateDiff(subResult, mpOld, mpNew, formatter, append(addr, k), ignoreAddrs)

			// has subkey differences, add to diff
			if len(subResult) > 0 {
				diffResult[k] = subResult
			}

			continue
		}

		// use deepEquals to determine equality b/c we don't dive into array-diffing
		// we just show the entire array as changed
		if !deepEquals(oldProp, newProp) {
			diffResult[k] = formatter(wrapJson(oldProp, append(addr, k), ignoreAddrs), wrapJson(newProp, append(addr, k), ignoreAddrs))
		}
	}
}

func wrapJson(value interface{}, currentAddr []string, ignoreAddrs [][]string) interface{} {
	if value == nil {
		return value
	}

	mp, ok := value.(map[string]interface{})
	if !ok {

		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice, reflect.Array:
			list := []interface{}{}

			arr := reflect.ValueOf(value)
			for i := 0; i < arr.Len(); i++ {
				list = append(list, wrapJson(arr.Index(i).Interface(), currentAddr, ignoreAddrs))
			}

			return list
		}

		return value
	}

	keysToRemove := []string{}
	for k := range mp {
		if containsAddr(ignoreAddrs, append(currentAddr, k)) {
			keysToRemove = append(keysToRemove, k)
		}
	}

	for _, k := range keysToRemove {
		delete(mp, k)
	}

	for k, v := range mp {
		mp[k] = wrapJson(v, append(currentAddr, k), ignoreAddrs)
	}

	return mp
}

func containsAddr(addrs [][]string, test []string) bool {
	for _, addr := range addrs {
		if equalsAddr(test, addr) {
			return true
		}
	}

	return false
}

func equalsAddr(test []string, prefix []string) bool {
	if len(prefix) != len(test) {
		return false
	}

	for i := 0; i < len(prefix); i++ {
		if test[i] != prefix[i] {
			return false
		}
	}

	return true
}

func deepEquals(a interface{}, b interface{}) bool {
	mpA, aIsMap := a.(map[string]interface{})
	mpB, bIsMap := b.(map[string]interface{})

	// one is a map, the other is not, must be a change
	if aIsMap != bIsMap {
		return false
	}

	// both are maps, check if entries match
	if aIsMap && bIsMap {
		for _, k := range allKeys(mpA, mpB) {
			if !deepEquals(mpA[k], mpB[k]) {
				return false
			}
		}

		return true
	}

	arrA, aIsArr := a.([]interface{})
	arrB, bIsArr := b.([]interface{})

	// one is an array, the other isn't, must be change
	if aIsArr != bIsArr {
		return false
	}

	// both are arrays, check if entries match
	if aIsArr && bIsArr {
		if len(arrA) != len(arrB) {
			return false
		}

		for i := 0; i < len(arrB); i++ {
			if !deepEquals(arrA[i], arrB[i]) {
				return false
			}
		}

		return true
	}

	// primitive comparison
	return a == b
}

// allKeys returns the combined keys of a & b (without duplicates)
func allKeys(a map[string]interface{}, b map[string]interface{}) []string {
	keyMap := map[string]bool{}

	for k := range a {
		keyMap[k] = true
	}

	for k := range b {
		keyMap[k] = true
	}

	var keys []string

	for k := range keyMap {
		keys = append(keys, k)
	}

	return keys
}
