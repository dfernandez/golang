package mynegroni

import (
	"code.google.com/p/gcfg"
	"os"
)

type Config struct {
	OAuth map[string]*struct {
		ClientId string
		SecretId string
	}
	Database map[string]*struct {
		Connector string
		Dns       string
	}
	Domain map[string]*struct {
		Url string
	}
	Stathat struct {
		Account string
	}
}

func LoadConfig() *Config {

	var config Config

	path := os.Getenv("GOPATH") + "/cfg/app.gcfg"
	err := gcfg.ReadFileInto(&config, path)

	if err != nil {
		panic(err)
	}

	return &config
}
