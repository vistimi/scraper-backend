package utils

import ("encoding/json")

// ToJSON transform the structured extracted data from html into json for printing.
// log.Println(ToJSON(data))
func ToJSON(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}