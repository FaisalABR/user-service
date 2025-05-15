package util

import (
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func BindFromJSON(dest any, filename, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
		return err
	}

	return nil
}

func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)
	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
		return err
	}

	for k, v := range env {
		var (
			valOf = reflect.ValueOf(v)
			val   string
		)

		switch valOf.Kind() {
		case reflect.String:
			val = valOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Uint()))
		case reflect.Float32:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Float64:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		}

		err = os.Setenv(k, val)
		if err != nil {
			logrus.Errorf("failed to set env: %v", err)
			return err
		}
	}

	return nil
}

func BindFromConsul(dest any, endpoint, path string) error {
	// initialize viper
	v := viper.New()

	// set config type and add remote provider
	v.SetConfigType("json")

	err := v.AddRemoteProvider("consul", endpoint, path)
	if err != nil {
		logrus.Errorf("failed to add remote: %v", err)
		return err
	}

	// read remote config
	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("failed to read remote: %v", err)
		return err
	}

	// unmarshal the dest
	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to read unmarhsal: %v", err)
		return err
	}

	// set env from consul kv
	err = SetEnvFromConsulKV(v)
	if err != nil {
		logrus.Errorf("failed to read remote: %v", err)
		return err
	}

	return nil
}
