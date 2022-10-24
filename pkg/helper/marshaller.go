package helper

import "encoding/json"

func JsonToJson(data interface{}, js interface{}) error {
	body, err := json.Marshal(js)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	return err
}
