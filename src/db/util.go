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

func createPrerequisites() error {

	transaction, err := connection.Begin()
	if err != nil {
		return err
	}

	statement1, _ := transaction.Prepare("CREATE TABLE IF NOT EXISTS path_mappings (path_mapping_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY , path VARCHAR(255) NOT NULL, `table` VARCHAR(255) NOT NULL) AUTO_INCREMENT=20000;")
	statement2, _ := transaction.Prepare("CREATE TABLE IF NOT EXISTS key_mappings (key_mapping_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY , `key` VARCHAR(255) NOT NULL, `column` VARCHAR(255) NOT NULL) AUTO_INCREMENT=30000;")
	statement3, _ := transaction.Prepare("CREATE TABLE IF NOT EXISTS behaviors (behavior_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY , path_mapping_id INT NOT NULL, key_mapping_id INT NOT NULL) AUTO_INCREMENT=10000;")

	_, err = statement1.Exec()
	if err != nil {

		err := transaction.Rollback()
		if err != nil {
			return err
		}

	}

	_, err = statement2.Exec()
	if err != nil {

		err := transaction.Rollback()
		if err != nil {
			return err
		}

	}

	_, err = statement3.Exec()
	if err != nil {

		err := transaction.Rollback()
		if err != nil {
			return err
		}

	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	return err
}

func FactoryReset() (err error) {

	transaction, err := connection.Begin()

	if err != nil {
		return err
	}

	statement1, err := transaction.Prepare("DROP TABLE `behaviors`;")
	statement2, err := transaction.Prepare("DROP TABLE `path_mappings`;")
	statement3, err := transaction.Prepare("DROP TABLE `key_mappings`;")

	_, err = statement1.Exec()
	if err != nil {

		err := transaction.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	_, err = statement2.Exec()
	if err != nil {

		err := transaction.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	_, err = statement3.Exec()
	if err != nil {

		err := transaction.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	err = createPrerequisites()
	if err != nil {
		return err
	}

	return nil

}
