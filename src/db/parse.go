package db

import (
	"fmt"
)

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
