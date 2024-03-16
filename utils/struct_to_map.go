package utils

import "encoding/json"

// Converts a struct to a generic Map type
func StructToMap(obj interface{}) (returnMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &returnMap) // Convert to a map
	return
}
