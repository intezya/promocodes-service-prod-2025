package pkg

import (
	"reflect"
)

func RecursiveRemoveNulls(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if isNil(value) {
				delete(v, key)
				continue
			}

			switch val := value.(type) {
			case map[string]interface{}:
				RecursiveRemoveNulls(val)
				if len(val) == 0 {
					delete(v, key)
				}
			case []map[string]interface{}:
				for i := range val {
					RecursiveRemoveNulls(val[i])
				}
				if len(val) == 0 {
					delete(v, key)
				}
			}
		}

	case []map[string]interface{}:
		for i := range v {
			RecursiveRemoveNulls(v[i])
		}
	}
}

func isNil(value interface{}) bool {
	if value == nil {
		return true
	}

	rv := reflect.ValueOf(value)
	return rv.Kind() == reflect.Ptr && rv.IsNil()
}
