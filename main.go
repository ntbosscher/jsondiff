package jsondiff

import (
	"encoding/json"
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

func DiffFormat(oldValue interface{}, newValue interface{}, formatter Formatter) (json.RawMessage, error) {
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
	calculateDiff(diff, oldMap, newMap, formatter)

	return json.Marshal(diff)
}

// calculateDiff calculates the difference between the old and new maps
// and fills diffResult with the result
func calculateDiff(
	diffResult map[string]interface{},
	oldMap map[string]interface{},
	newMap map[string]interface{},
	formatter Formatter,
) {

	// iterate over keys
	for _, k := range allKeys(oldMap, newMap) {

		newProp := newMap[k]
		oldProp := oldMap[k]

		// check if the values are maps themselves
		mpOld, oldIsMap := oldProp.(map[string]interface{})
		mpNew, newIsMap := newProp.(map[string]interface{})

		// one is a map, the other is not, must be a change
		if oldIsMap != newIsMap {
			diffResult[k] = formatter(oldProp, newProp)
			continue
		}

		// both are maps, check subkeys for changes
		if oldIsMap && newIsMap {
			subResult := map[string]interface{}{}
			calculateDiff(subResult, mpOld, mpNew, formatter)

			// has subkey differences, add to diff
			if len(subResult) > 0 {
				diffResult[k] = subResult
			}

			continue
		}

		// use deepEquals to determine equality b/c we don't dive into array-diffing
		// we just show the entire array as changed
		if !deepEquals(oldProp, newProp) {
			diffResult[k] = formatter(oldProp, newProp)
		}
	}
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