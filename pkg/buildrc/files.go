package buildrc

import "encoding/json"

func (me *BuildrcJSON) Files() (map[string]string, error) {
	ok, err := json.Marshal(me)
	if err != nil {
		return nil, err
	}

	var res map[string]string

	err = json.Unmarshal(ok, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
