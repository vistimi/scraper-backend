package utils

import ("encoding/json")

// To transform the structured extracted data from html into json for printing.
// log.Println(toJson(data))
func ToJson(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}