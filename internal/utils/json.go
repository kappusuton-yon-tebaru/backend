package utils

import "encoding/json"

func MapToStruct(mp map[string]any, s interface{}) error {
	bs, err := json.Marshal(mp)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, s)
	if err != nil {
		return err
	}

	return nil
}
