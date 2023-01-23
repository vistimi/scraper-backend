package util

import "encoding/json"

func ConvertMap[K comparable, VF, VT any](from map[K]VF) (map[K]VT, error) {
	to := make(map[K]VT, len(from))
	for k,v :=  range(from){
		var t VT
		temp, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(temp, &t)
		to[k] = t
	}
	return to, nil
}
