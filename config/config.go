package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

//Database config
type DBConfig struct {
	Database string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
	Schema   string `yaml:"schema"`
}

// Config struct of application config
type Config struct {
	Database DBConfig  `yaml:"database"`

	Envs []string `yaml:",flow"`
}

// loadEnvs load system environments with keys listed in ...e
// then write values to map and return map[key]value
func loadEnvs(e ...string) map[string]string {
	envs := make(map[string]string)

	for _, k := range e {
		envs[k] = os.Getenv(k)
	}

	return envs
}

// New Create new instance of app config with given path to (../config.yml)
func New(p ...string) Config {
	path := "./dev.yml"
	if len(p) > 0 {
		path = p[0]
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	config := Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	envs := loadEnvs(config.Envs...)

	v, ok := envs["db_pass"]
	if ok && v != "" {
		config.Database.Password = v
	}
	v, ok = envs["db_user"]
	if ok && v != "" {
		config.Database.User = v
	}
	v, ok = envs["db_address"]
	if ok && v != "" {
		config.Database.Address = v
	}

	return config
}