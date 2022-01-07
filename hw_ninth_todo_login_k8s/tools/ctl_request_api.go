package tools

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Respond_msg struct {
	Data   map[string]interface{} `json:"data"`
	Status string                 `json:"status"`
	Msg    string                 `json:"msg"`
}

func Request_api(url string, method string, payload io.ReadCloser, data map[string]string) (*Respond_msg, error) {
	// set new request
	client := &http.Client{}
	var req *http.Request
	var err error
	if payload != nil && data != nil {
		// payload combine data
		// payload read to byte
		payload_map := make(map[string]string)
		payload_byte, err := io.ReadAll(payload)
		if err != nil {
			return nil, err
		}

		if len(payload_byte) != 0 {
			// payload byte to json map
			err = json.Unmarshal(payload_byte, &payload_map)
			if err != nil {
				return nil, err
			}

			// add data to json map
			for k, v := range payload_map {
				data[k] = v
			}
		}

		// json map to byte
		json_data, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		// set request
		req, err = http.NewRequest(method, url, bytes.NewBuffer(json_data))
		if err != nil {
			return nil, err
		}
	} else if payload != nil {
		// set request
		req, err = http.NewRequest(method, url, payload)
		if err != nil {
			return nil, err
		}
	} else if data != nil {
		// set data to byte
		json_data, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		// set request
		req, err = http.NewRequest(method, url, bytes.NewBuffer(json_data))
		if err != nil {
			return nil, err
		}

	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}

	}
	if err != nil {
		return nil, err
	}

	// send request and receive respond
	res, err := client.Do(req)
	if err != nil {
		return &Respond_msg{}, err
	}
	defer res.Body.Close()

	// respond decode
	var r_msg *Respond_msg
	err = json.NewDecoder(res.Body).Decode(&r_msg)
	if err != nil {
		return &Respond_msg{}, err
	}
	return r_msg, nil
}
