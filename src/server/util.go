package server

import "grest/db"

func correctData(behavior db.Behavior, data map[string]string) (correctedData map[string]string, unknownKeys []string) {

	correctedData = map[string]string{}

	for k, v := range data {

		found := false

		for i := range behavior.KeyMappings {

			bKey := behavior.KeyMappings[i].Key
			bColumn := behavior.KeyMappings[i].Column

			if k == bColumn {

				correctedData[bKey] = v
				found = true
				break

			}

		}

		if !found {

			unknownKeys = append(unknownKeys, k)

		}

	}

	return

}
