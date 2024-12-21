package utils

import (
	"fmt"
	"strconv"
)

func ConvertToStringSlice(interfaceSlice []any) ([]string, error) {
	fmt.Println("converting to string slice", interfaceSlice)
	var stringSlice []string
	for i, v := range interfaceSlice {
		switch strVal := v.(type) {
		case string:
			stringSlice = append(stringSlice, strVal)
		case int:
			stringSlice = append(stringSlice, strconv.Itoa(strVal))
		case int64:
			stringSlice = append(stringSlice, strconv.FormatInt(strVal, 10))
		case float64:
			stringSlice = append(stringSlice, strconv.FormatFloat(strVal, 'G', 10, 64))
		case bool:
			stringSlice = append(stringSlice, strconv.FormatBool(strVal))
		case []any:
			subSlice, err := ConvertToStringSlice(strVal)
			if err != nil {
				return nil, err
			}
			stringSlice = append(stringSlice, subSlice...)
		default:
			return nil, fmt.Errorf("element at index %d is of unsupported type %T", i, v)
		}
	}
	return stringSlice, nil
}
