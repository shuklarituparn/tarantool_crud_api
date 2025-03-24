package utils

import "fmt"

func ConvertToJSONSafe(data interface{}) interface{} {
	switch val := data.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			strKey := fmt.Sprintf("%v", k)
			result[strKey] = ConvertToJSONSafe(v)
		}
		return result
	case []interface{}:
		for i, item := range val {
			val[i] = ConvertToJSONSafe(item)
		}
		return val
	default:
		return data
	}
}
