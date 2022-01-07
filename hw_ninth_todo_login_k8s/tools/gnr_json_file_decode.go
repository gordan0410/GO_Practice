package tools

import (
	"encoding/json"
	"os"
)

func Json_file_convert(addr string, str interface{}) (interface{}, error) {
	file, err := os.Open(addr)
	if err != nil {
		return nil, err
	}

	f_data := json.NewDecoder(file)
	err = f_data.Decode(str)
	if err != nil {
		return nil, err
	}
	return str, nil
}
