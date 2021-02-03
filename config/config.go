package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

// Authorization config
type Authorization struct {
	// Jwt salt is randomly string which will be additional added to jwt token
	JwtSalt string `yaml:"jwt_salt"`
	// When jwt is expired
	JwtExpires time.Duration `yaml:"jwt_expires"`
	// Who created jwt token
	JwtIss string `yaml:"jwt_iss"`
	// Who can use the token
	//
	// by default everyone can use the token
	JwtAud string `yaml:"jwt_aud"`
}

//Database config
type DBConfig struct {
	// Database name
	Database string `yaml:"name"`
	// Postgres user
	//
	// for example: postgres
	User string `yaml:"user"`
	// Postgres password
	//
	// for example: (nothing :) )
	Password string `yaml:"password"`
	// Postgres machine hostname
	//
	// for example: localhost, mysite.com
	Address string `yaml:"address"`
	// Port on which postgres
	//
	// for example: 5432
	Port string `yaml:"port"`
	// The path to the file from which the database schema will be generated
	Schema string `yaml:"schema"`
	// Instance of redis config
	Redis RedisConfig `yaml:"redis"`
}

// Redis cache config
type RedisConfig struct {
	// Redis machine hostname
	//
	// for example: localhost, mysite.com
	Address  string `yaml:"address"`
	// Port on which redis
	//
	// for example 6379
	Port     string `yaml:"port"`
	// Password to redis
	//
	// by default: (empty)
	Password string `yaml:"password"`
	// Database name
	//
	// by default: 0
	Database int    `yaml:"database"`
}

// Config struct of application config
type Config struct {
	// Database configs
	Database DBConfig      `yaml:"database"`
	// Auth config
	Auth     Authorization `yaml:"auth"`

	// Sys envs
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

// New Create new instance of apprepo config with given path to (../config.yml)
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
	if err = yaml.Unmarshal(data, &config); err != nil {
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
