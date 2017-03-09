package config

import (
	"bufio"
	"os"
	"reflect"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Config struct {
	DemoAddr    string `env:"DEMO_ADDR" required:"true"`
	DemoAppAddr string `env:"DEMO_APP_ADDR" required:"true"`
}

var config *Config

func GetConfig() *Config {
	return config
}

func InitConfig(envFile string) *Config {
	configFileParams := LoadConfigFile(envFile)

	config = new(Config)
	if err := LoadConfig(config, configFileParams); err != nil {
		log.Error("LoadConfig got error: ", err)
		os.Exit(1)
	}

	return config
}

func LoadConfigFile(envfile string) map[string]string {
	log.Debug("envfile: ", envfile)
	configMap := make(map[string]string)
	fp, err := os.Open(envfile)
	if err != nil {
		log.Errorf("Failed to open config file %s: %s", envfile, err.Error())
		return nil
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		trimed := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(trimed, "#") && len(trimed) > 0 {
			parts := strings.SplitN(trimed, "=", 2)
			if len(parts) != 2 {
				log.Warnf("Invalid config file %s, %v", envfile, parts)
				continue
			}

			key, val := strings.Trim(parts[0], " "), strings.Trim(parts[1], ` "'`)
			configMap[strings.ToUpper(key)] = val
		}
	}

	return configMap
}

func exitMissingEnv(env string) {
	log.Errorf("program exit missing config for env %s", env)
	os.Exit(1)
}

func exitCheckEnv(env string, err error) {
	log.Errorf("Check env %s got error: %s", env, err.Error())
}

func LoadConfig(configEntry interface{}, configFileParams map[string]string) error {
	val := reflect.ValueOf(configEntry).Elem()

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		required := typeField.Tag.Get("required")
		envKey := typeField.Tag.Get("env")

		env := os.Getenv(envKey)
		if env == "" {
			env = configFileParams[envKey]
		}

		if env == "" {
			env = typeField.Tag.Get("default")
		}

		if env == "" && required == "true" {
			exitMissingEnv(envKey)
		}

		var configEntryValue interface{}
		var err error
		valueFiled := val.Field(i).Interface()
		value := val.Field(i)
		switch valueFiled.(type) {
		case int:
			configEntryValue, err = strconv.Atoi(env)
		case int64:
			configEntryValue, err = strconv.ParseInt(env, 10, 64)
		case int16:
			configEntryValue, err = strconv.ParseInt(env, 10, 16)
			_, ok := configEntryValue.(int64)
			if !ok {
				exitCheckEnv(typeField.Name, err)
			}
			configEntryValue = int16(configEntryValue.(int64))
		case uint16:
			configEntryValue, err = strconv.ParseUint(env, 10, 16)

			_, ok := configEntryValue.(uint64)
			if !ok {
				exitCheckEnv(typeField.Name, err)
			}
			configEntryValue = uint16(configEntryValue.(uint64))
		case uint64:
			configEntryValue, err = strconv.ParseUint(env, 10, 64)
		case bool:
			configEntryValue, err = strconv.ParseBool(env)
		case []string:
			configEntryValue = strings.SplitN(env, ",", -1)
		default:
			configEntryValue = env
		}

		if err != nil {
			exitCheckEnv(typeField.Name, err)
		}
		value.Set(reflect.ValueOf(configEntryValue))
	}

	return nil
}
