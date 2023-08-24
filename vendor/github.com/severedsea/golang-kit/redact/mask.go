package redact

import (
	"encoding/json"
	"strings"
)

const (
	// secretMask is the masked text for secrets
	secretMask = "XXXXXXXX"

	// maskChar is the masked character used to mask PII (Personally identifiable information) values
	maskChar = "X"
	// piiClearCount is the number of characters from the start of the string that are not masked (clear)
	piiClearCount = 4
)

// MaskMap attempts to mask a map[string]interface{} input
func MaskMap(input map[string]interface{}, redactionKeys Keys) map[string]interface{} {
	result := MaskObject(input, redactionKeys)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return input
	}

	return resultMap
}

// MaskObject attempts to mask an unknown input
func MaskObject(input any, redactionKeys Keys) any {
	// interface is of type map[string]interface{}, iterate over items and mask if value is string
	// or slice of strings, recurse otherwise
	m, ok := input.(map[string]interface{})
	if ok {
		for key, value := range m {
			switch v := value.(type) {
			case []interface{}:
				for i, it := range v {
					if s, ok := it.(string); ok && !isJSONObject(s) { // primitive string
						v[i] = getRedactedValue(key, s, redactionKeys)
					} else { // JSON object or string
						v[i] = MaskObject(it, redactionKeys)
					}
				}

			case string:
				if !isJSONObject(v) { // primitive string
					(m)[key] = getRedactedValue(key, v, redactionKeys)
				} else { // JSON string
					(m)[key] = MaskObject(value, redactionKeys)
				}

			default:
				(m)[key] = MaskObject(value, redactionKeys)
			}
		}

		return m
	}

	// interface is not of type map[string]interface{}, attempt to convert it to
	// map[string]interface{}
	switch value := input.(type) {
	case []interface{}:
		for i, it := range value {
			value[i] = MaskObject(it, redactionKeys)
		}

		return value

	// attempts to unmarshal string into JSON before recursing
	case string:
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(value), &m); err == nil {
			return MaskObject(m, redactionKeys)
		}

		var s []interface{}
		if err := json.Unmarshal([]byte(value), &s); err == nil {
			return MaskObject(s, redactionKeys)
		}

		return value

	// attempts to unmarshal bytes into JSON before recursing
	case []byte:
		var m map[string]interface{}
		if err := json.Unmarshal(value, &m); err == nil {
			return MaskObject(m, redactionKeys)
		}

		var s []interface{}
		if err := json.Unmarshal(value, &s); err == nil {
			return MaskObject(s, redactionKeys)
		}

		return value

	// attempts to convert struct or primitive to map[string]interface{} before recursing - if value
	// is a primitive (e.g. int, bool), the unmarshalling will fail and the value will be returned,
	// ending the recursion
	default:
		bytes, err := json.Marshal(&value)
		if err != nil {
			return value
		}

		var result map[string]interface{}
		if err := json.Unmarshal(bytes, &result); err != nil {
			return value
		}

		return MaskObject(result, redactionKeys)
	}
}

// getRedactedValue determines the value to be set after checking whether the key provided is in the Keys list
func getRedactedValue(key, value string, redactKeys Keys) string {
	// Redact secrets completely
	if isSensitiveData(key, redactKeys.Secrets) {
		return secretMask
	}

	// Redact value by masking the middle part with maskChar
	if isSensitiveData(key, redactKeys.UIN) {
		return MaskUIN(value)
	}

	// Redact value by masking the middle part with maskChar
	if isSensitiveData(key, redactKeys.Masked) {
		return Mask(value, maskChar, piiClearCount)
	}

	return value
}

// isSensitiveData validates whether a key is in the sensitive data keys provided
func isSensitiveData(key string, sensitiveKeys []string) bool {
	for _, k := range sensitiveKeys {
		if strings.Contains(sanitizeKey(key), sanitizeKey(k)) {
			return true
		}
	}

	return false
}

// Mask masks the middle portion with the character/s provided as c in the argument
// clearCount represents the number of characters that are not masked
//
// Sample:
// | s            | c | clearCount | result       |
// |:------------:|:-:|:----------:|:------------:|
// | P            | x | -1         | x            |
// | PE           | x | -1         | xx           |
// | PLEASEMASKME | x | -1         | xxxxxxxxxxxx |
// | P            | x | 0          | x            |
// | PE           | x | 0          | xx           |
// | PLEASEMASKME | x | 0          | xxxxxxxxxxxx |
// | P            | x | 1          | x            |
// | PE           | x | 1          | Px           |
// | PLEASEMASKME | x | 1          | Pxxxxxxxxxxx |
// | P            | x | 2          | x            |
// | PE           | x | 2          | Px           |
// | PLE          | x | 2          | Pxx          |
// | PLEASEMASKME | x | 2          | PxxxxxxxxxxE |
// | P            | x | 3          | x            |
// | PE           | x | 3          | Px           |
// | PLE          | x | 3          | Pxx          |
// | PLSE         | x | 3          | PxxE         |
// | PLASE        | x | 3          | PxxxE        |
// | PLEASEMASKME | x | 3          | PLxxxxxxxxxE |
// | P            | x | 4          | x            |
// | PE           | x | 4          | Px           |
// | PLE          | x | 4          | Pxx          |
// | PLSE         | x | 4          | PxxE         |
// | PLASE        | x | 4          | PxxxE        |
// | PLEASME      | x | 4          | PLxxxxE      |
// | PLEASEME     | x | 4          | PLxxxxME     |
// | PLEASEMASKME | x | 4          | PLxxxxxxxxME |
// | P            | x | 5          | x            |
// | PE           | x | 5          | Px           |
// | PLE          | x | 5          | Pxx          |
// | PLSE         | x | 5          | PxxE         |
// | PLASE        | x | 5          | PxxxE        |
// | PLEASE       | x | 5          | PLxxxE       |
// | PLEASEME     | x | 5          | PLxxxxME     |
// | PLEASEMASKME | x | 5          | PLExxxxxxxME |
func Mask(input, maskChar string, clearCount int) string {
	strLen := len(input)
	if clearCount <= 0 {
		return strings.Repeat(maskChar, strLen)
	}

	// Handle short string length relative to clear count
	if strLen < clearCount*2 {
		clearCount = strLen / 2
	}

	// Right clear string
	rightClearCount := clearCount / 2
	leftClearCount := clearCount - rightClearCount
	leftClear := input[0:leftClearCount]

	// Masked string
	maskedCount := strLen - clearCount
	masked := strings.Repeat(maskChar, maskedCount)

	// Right clear string
	var rightClear string
	if rightClearCount > 0 {
		rightClear = input[leftClearCount+maskedCount:]
	}

	return leftClear + masked + rightClear
}

// MaskUIN masks the front part of the UIN
//
// example:
//
// S1234567D = SXXXX567D
func MaskUIN(input string) string {
	if len(input) < 9 {
		return input
	}

	res := string(input[0])
	res += "XXXX"
	res += input[len(input)-4:]

	return res
}

func isJSONObject(str string) bool {
	var js map[string]interface{}

	return json.Unmarshal([]byte(str), &js) == nil
}
