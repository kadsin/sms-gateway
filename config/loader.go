package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	env "github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

func init() {
	_, thisFile, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(thisFile)

	if testing.Testing() {
		Load(fmt.Sprintf("%v/../.env.testing", basepath))
	} else {
		Load(fmt.Sprintf("%v/../.env", basepath))
	}
}

func Load(envFilePath string) {
	if err := godotenv.Load(envFilePath); err != nil {
		log.Printf("Error on loading .env file from %v: %+v\n", envFilePath, err)
	}

	if _, err := env.UnmarshalFromEnviron(&Env); err != nil {
		log.Fatalf("Error on unmarshaling .env file: %+v\n", err)
	}
}

var Env struct {
	App struct {
		Name        string `env:"APP_NAME"`
		Environment string `env:"APP_ENV"`
		Debug       bool   `env:"APP_DEBUG"`
	}

	DB struct {
		Connection string `env:"DB_CONNECTION"`
		Host       string `env:"DB_HOST"`
		Port       string `env:"DB_PORT"`
		Username   string `env:"DB_USERNAME"`
		Password   string `env:"DB_PASSWORD"`
		Name       string `env:"DB_NAME"`
	}

	Analytics struct {
		Connection string `env:"ANALYTICS_DB_CONNECTION"`
		Host       string `env:"ANALYTICS_DB_HOST"`
		Cluster    string `env:"ANALYTICS_DB_CLUSTER,default=default"`
		Port       string `env:"ANALYTICS_DB_PORT"`
		Username   string `env:"ANALYTICS_DB_USERNAME"`
		Password   string `env:"ANALYTICS_DB_PASSWORD"`
		Name       string `env:"ANALYTICS_DB_NAME"`
	}

	Doc struct {
		Auth struct {
			Username string `env:"DOC_AUTH_USERNAME"`
			Password string `env:"DOC_AUTH_PASSWORD"`
		}
	}

	Queue struct {
		Driver string `env:"QUEUE_DRIVER"`
		Host   string `env:"QUEUE_HOST"`
		Port   string `env:"QUEUE_PORT"`
		Names  struct {
			Regular string `env:"QUEUE_NAME_REGULAR"`
			Express string `env:"QUEUE_NAME_EXPRESS"`
		}
	}
}
