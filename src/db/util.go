package db

import "fmt"

func CompareBehaviors(behavior1 Behavior, behavior2 Behavior) bool {

	if !ComparePathMappings(behavior1.PathMapping, behavior2.PathMapping) {
		return false
	}

	if !(len(behavior1.KeyMappings) == len(behavior2.KeyMappings)) {
		return false
	}

	for i := range behavior1.KeyMappings {

		if !CompareKeyMappings(behavior1.KeyMappings[i], behavior2.KeyMappings[i]) {
			return false
		}

	}

	return true
}

func CompareKeyMappings(keyMapping1 KeyMapping, keyMapping2 KeyMapping) bool {

	if keyMapping1.Key == keyMapping2.Key && keyMapping1.Column == keyMapping2.Column {
		return true
	} else {
		return false
	}

}

func ComparePathMappings(pathMapping1 PathMapping, pathMapping2 PathMapping) bool {

	if pathMapping1.Path == pathMapping2.Path && pathMapping1.Table == pathMapping2.Table {
		return true
	} else {
		return false
	}

}

func parseInterfacesToMapSlice(unparsedData []map[string]interface{}) (parsedData []map[string]string) {

	for index := range unparsedData {

		var parsedDatum = map[string]string{}

		for k, v := range unparsedData[index] {

			parsedDatum[k] = fmt.Sprintf("%s", v)
		}

		parsedData = append(parsedData, parsedDatum)

	}
	return

}
