package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Env represents the application environment.
type Env struct {
	logger         string
	mySQLHost      string
	mySQLName      string
	mySQLPassword  string
	mySQLPort      int
	mySQLUser      string
	port           int
	redisHost      string
	redisPort      int
	redisKeyPrefix string
	redisPassword  string
}

const (
	envKeyLogger         = "LOGGER"
	envKeyMySQLHost      = "MYSQL_HOST"
	envKeyMySQLName      = "MYSQL_NAME"
	envKeyMySQLPassword  = "MYSQL_PASSWORD"
	envKeyMySQLPort      = "MYSQL_PORT"
	envKeyMySQLUser      = "MYSQL_USER"
	envKeyPort           = "PORT"
	envKeyRedisHost      = "REDIS_HOST"
	envKeyRedisPort      = "REDIS_PORT"
	envKeyRedisPassword  = "REDIS_PASSWORD"
	envKeyRedisKeyPrefix = "REDIS_KEY_PREFIX"
)

const (
	envFileName = ".env"
	envFilePath = "."
	envFileType = "env"
)

const (
	defaultPort   = 3000
	defaultLogger = "PRODUCTION"
)

// ReadEnv returns the application environment configuration.
func ReadEnv() (*Env, error) { //nolint: gocyclo
	env := Env{}

	v := viper.New()
	v.AutomaticEnv()
	v.AddConfigPath(envFilePath)
	v.SetConfigType(envFileType)
	v.SetConfigName(envFileName)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if port := v.GetInt(envKeyPort); port != 0 {
		env.port = port
	} else {
		env.port = defaultPort
	}

	if logger := v.GetString(envKeyLogger); logger != "" {
		env.logger = logger
	} else {
		env.logger = defaultLogger
	}

	if mySQLHost := v.GetString(envKeyMySQLHost); mySQLHost != "" {
		env.mySQLHost = mySQLHost
	} else {
		return nil, missingEnvError(envKeyMySQLHost)
	}

	if mySQLPort := v.GetInt(envKeyMySQLPort); mySQLPort != 0 {
		env.mySQLPort = mySQLPort
	} else {
		return nil, missingEnvError(envKeyMySQLPort)
	}

	if mySQLName := v.GetString(envKeyMySQLName); mySQLName != "" {
		env.mySQLName = mySQLName
	} else {
		return nil, missingEnvError(envKeyMySQLName)
	}

	if mySQLUser := v.GetString(envKeyMySQLUser); mySQLUser != "" {
		env.mySQLUser = mySQLUser
	}

	if mySQLPassword := v.GetString(envKeyMySQLPassword); mySQLPassword != "" {
		env.mySQLPassword = mySQLPassword
	}

	if redisHost := v.GetString(envKeyRedisHost); redisHost != "" {
		env.redisHost = redisHost
	} else {
		return nil, missingEnvError(envKeyRedisHost)
	}

	if redisPort := v.GetInt(envKeyRedisPort); redisPort != 0 {
		env.redisPort = redisPort
	} else {
		return nil, missingEnvError(envKeyRedisPort)
	}

	if redisPassword := v.GetString(envKeyRedisPassword); redisPassword != "" {
		env.redisPassword = redisPassword
	}

	if redisKeyPrefix := v.GetString(envKeyRedisKeyPrefix); redisKeyPrefix != "" {
		env.redisKeyPrefix = redisKeyPrefix
	}

	return &env, nil
}

// GetMySQLURL returns the MySQL URL.
func (e *Env) GetMySQLURL() string {
	if e.mySQLUser == "" || e.mySQLPassword == "" {
		return fmt.Sprintf("tcp(%s:%d)/%s?parseTime=true",
			e.mySQLHost,
			e.mySQLPort,
			e.mySQLName)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		e.mySQLUser,
		e.mySQLPassword,
		e.mySQLHost,
		e.mySQLPort,
		e.mySQLName)
}

// GetLogger returns the logger type.
func (e *Env) GetLogger() string {
	return e.logger
}

// GetRedisHost returns the Redis host.
func (e *Env) GetRedisHost() string {
	return e.redisHost
}

// GetRedisPort returns the Redis port.
func (e *Env) GetRedisPort() int {
	return e.redisPort
}

// GetRedisPassword returs the Redis password.
func (e *Env) GetRedisPassword() string {
	return e.redisPassword
}

// GetPort returns the port.
func (e *Env) GetPort() int {
	return e.port
}

// GetRedisKeyPrefix returns the redis cache key prefix.
func (e *Env) GetRedisKeyPrefix() string {
	return e.redisKeyPrefix
}

func missingEnvError(key string) error {
	return fmt.Errorf("%w: %s", ErrorMissing, key)
}
