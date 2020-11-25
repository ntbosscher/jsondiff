package jsondiff

import (
	"encoding/json"
)

// Diff compares oldValue with newValue and returns a json tree of
// the changed values.
func Diff(oldValue interface{}, newValue interface{}) (json.RawMessage, error) {
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
	calculateDiff(diff, oldMap, newMap)

	return json.Marshal(diff)
}

// calculateDiff calculates the difference between the old and new maps
// and fills diffResult with the result
func calculateDiff(
	diffResult map[string]interface{},
	oldMap map[string]interface{},
	newMap map[string]interface{},
) {

	// iterate over keys
	for k, oldProp := range oldMap {
		newProp := newMap[k]

		// check if the values are maps themselves
		mpOld, oldIsMap := oldProp.(map[string]interface{})
		mpNew, newIsMap := newProp.(map[string]interface{})

		// one is a map, the other is not, must be a change
		if oldIsMap != newIsMap {
			diffResult[k] = newProp
			continue
		}

		// both are maps, check subkeys for changes
		if oldIsMap && newIsMap {
			subResult := map[string]interface{}{}
			calculateDiff(subResult, mpOld, mpNew)

			// has subkey differences, add to diff
			if len(subResult) > 0 {
				diffResult[k] = subResult
			}

			continue
		}

		// regular value comparison (non-map values)
		if newProp != oldProp {
			diffResult[k] = newProp
		}
	}
}
