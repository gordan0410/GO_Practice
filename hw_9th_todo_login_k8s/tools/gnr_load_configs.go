package tools

import (
	"encoding/json"
	"os"
)

type Config_data struct {
	Jwt struct {
		JwtKey    string `json:"Jwt_key"`
		JwtMaxage int    `json:"Jwt_maxage"`
	} `json:"jwt"`
	Mysql struct {
		Username     string `json:"Username"`
		Password     string `json:"Password"`
		Addr         string `json:"Addr"`
		Database     string `json:"Database"`
		MaxLifetime  int    `json:"Max_lifetime"`
		MaxOpenconns int    `json:"Max_openconns"`
		MaxIdleconns int    `json:"Max_idleconns"`
	} `json:"mysql"`
	Redis struct {
		Size     int    `json:"Size"`
		Network  string `json:"Network"`
		Address  string `json:"Address"`
		Password string `json:"Password"`
	} `json:"redis"`
	Session struct {
		SessionName   string `json:"Session_name"`
		SessionPrefix string `json:"Session_prefix"`
		SessionKey    string `json:"Session_key"`
		SessionMaxage int    `json:"Session_maxage"`
	} `json:"session"`
}

func Load_configs(addr string) (*Config_data, error) {
	var configs Config_data
	file, err := os.Open(addr)
	if err != nil {
		return nil, err
	}

	f_data := json.NewDecoder(file)
	err = f_data.Decode(&configs)
	if err != nil {
		return nil, err
	}
	return &configs, nil
}
