package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port   string `json:"port"`
	Name   string `json:"name"`
	Ollama Ollama `json:"ollama"`
	Milvus struct {
		Host       string `json:"host"`
		Port       string `json:"port"`
		User       string `json:"user"`
		Pass       string `json:"pass"`
		Dbname     string `json:"dbname"`
		SearchPath string `json:"search_path"`
		IdleConns  int    `json:"idle_conns"`
		OpenConns  int    `json:"open_conns"`
		DriverName string `json:"driver_name"`
	} `json:"milvus"`
}

type Ollama struct {
	Url       string            `json:"url"`
	Endpoints map[string]string `json:"endpoints"`
	Timeout   time.Duration     `json:"timeout"`
}

func LoadConfig(path string) (*Config, error) {
	config := new(Config)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can not open the config file: %s(%w)", path, err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println("error closing body: " + err.Error())
		}
	}(file)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
