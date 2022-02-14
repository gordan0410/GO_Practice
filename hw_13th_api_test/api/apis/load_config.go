package apis

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	API    struct {
		Port  string `json:"port"`
		Key   string `json:"key"`
		Value string `json:"value"`
		Path1 struct {
			ID string `json:"id"`
		} `json:"path1"`
		Path2 struct {
			ID string `json:"id"`
		} `json:"path2"`
	} `json:"api"`
	V1 struct {
		Port  string `json:"port"`
		Key   string `json:"key"`
		Value string `json:"value"`
		Path1 struct {
			ID string `json:"id"`
		} `json:"path1"`
		Path2 struct {
			ID string `json:"id"`
		} `json:"path2"`
	} `json:"v1"`
	V2 struct {
		Port  string `json:"port"`
		Key   string `json:"key"`
		Value string `json:"value"`
		Path1 struct {
			ID string `json:"id"`
		} `json:"path1"`
		Path2 struct {
			ID string `json:"id"`
		} `json:"path2"`
	} `json:"v2"`
}

var conf *Config

func Load_config(addr string) {
	file, err := os.Open(addr)
	if err != nil {
		log.Println(err)
	}
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		log.Println(err)
	}
}
