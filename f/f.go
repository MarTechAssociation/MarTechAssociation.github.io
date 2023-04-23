// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package f

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func StructToMap(s interface{}) map[string]interface{} {
	var mapped map[string]interface{}
	inRecord, _ := json.Marshal(s)
	json.Unmarshal(inRecord, &mapped)

	return mapped
}

func InterfaceToString(x interface{}) string {
	if x == nil {
		return ""
	}

	// When convert from JSON number to number in Golang
	// If the number is big, the json.Marshal will use float64 as type of conversion
	// even it just a plain integer and has nothing relate to the decimal number
	// The block of code below will deal with this situation when the number is float32 or float64
	// and we want to transform it to string

	// This is how to convert float32 to string by not using scietific notation
	// The default %v will convert float32 using scientific notation, that we don't want
	v32, ok := x.(float32)
	if ok {
		return strconv.FormatFloat(float64(v32), 'f', -1, 32)
	}
	// This is how to convert float64 to string by not using scietific notation
	// The default %v will convert float64 using scientific notation, that we don't want
	v64, ok := x.(float64)
	if ok {
		return strconv.FormatFloat(v64, 'f', -1, 64)
	}

	// Try pass to string
	str, ok := x.(string)
	if ok {
		return str
	}

	// Try pass to []byte
	b, ok := x.([]byte)
	if ok {
		return string(b)
	}

	// Default case use Sprintf
	return fmt.Sprintf("%v", x)
}
