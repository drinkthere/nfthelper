package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	// 日志配置
	LogLevel string
	LogPath  string

	// 电报是否开启debug模式
	IsTgDebug bool
}

func LoadConfig(filename string) *Config {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := new(Config)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	return data
}
