package ectoenv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ENV_TAG = "env"

var ENV_DEFAULT_TAG = "env-default"

// AUTO_REFRESH_INTERVAL is the interval in seconds to refresh the environment variables
var AUTO_REFRESH_INTERVAL = 60

// BindEnv sets the values of the provided struct based on the values of the environment variables
// defined in the struct's tags. The struct must be a non-nil pointer to a struct.
// v: a non-nil pointer to a struct
// returns: an error if the provided value is not a non-nil pointer to a struct or if the value of an environment variable
func BindEnv(v interface{}) error {
	rv, err := validateInput(v)
	if err != nil {
		return err
	}

	return setFieldValues(rv)
}

func validateInput(v interface{}) (reflect.Value, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return reflect.Value{}, errors.New("provided value must be a non-nil pointer to a struct")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("provided value must be a pointer to a struct")
	}

	return rv, nil
}

func setFieldValues(rv reflect.Value) error {
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if !field.CanSet() {
			continue
		}

		if field.Kind() == reflect.Struct {
			if err := BindEnv(field.Addr().Interface()); err != nil {
				return fmt.Errorf("unable to set value for field %s: %w", field.Type().Name(), err)
			}
			continue
		}

		envTag := rt.Field(i).Tag.Get(ENV_TAG)
		if envTag == "" {
			continue
		}

		envValue := getEnvValue(rt.Field(i), envTag)
		if envValue == "" {
			continue
		}

		if err := setFieldValue(field, envValue); err != nil {
			return err
		}
	}

	return nil
}

func getEnvValue(field reflect.StructField, envTag string) string {
	envValue := os.Getenv(envTag)
	if envValue == "" {
		defaultTag := field.Tag.Get(ENV_DEFAULT_TAG)
		if defaultTag != "" {
			envValue = defaultTag
		}
	}
	return envValue
}

func setFieldValue(field reflect.Value, envValue string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(envValue)
	case reflect.Int:
		return setIntField(field, envValue)
	case reflect.Bool:
		return setBoolField(field, envValue)
	case reflect.Float64:
		return setFloat64Field(field, envValue)
	case reflect.Slice:
		return setSliceField(field, envValue)
	}
	return nil
}

func setIntField(field reflect.Value, envValue string) error {
	val, err := strconv.Atoi(envValue)
	if err != nil {
		return fmt.Errorf("unable to set value for field %s. failed to parse %s as int: %w", field.Type().Name(), envValue, err)
	}
	field.SetInt(int64(val))
	return nil
}

func setBoolField(field reflect.Value, envValue string) error {
	val, err := strconv.ParseBool(envValue)
	if err != nil {
		return fmt.Errorf("unable to set value for field %s. failed to parse %s as bool: %w", field.Type().Name(), envValue, err)
	}
	field.SetBool(val)
	return nil
}

func setFloat64Field(field reflect.Value, envValue string) error {
	val, err := strconv.ParseFloat(envValue, 64)
	if err != nil {
		return fmt.Errorf("unable to set value for field %s. failed to parse %s as float64: %w", field.Type().Name(), envValue, err)
	}
	field.SetFloat(val)
	return nil
}

func setSliceField(field reflect.Value, envValue string) error {
	split := strings.Split(envValue, ",")
	switch field.Type().Elem().Kind() {
	case reflect.String:
		field.Set(reflect.ValueOf(split))
	case reflect.Bool:
		return setBoolSlice(field, split)
	case reflect.Float64:
		return setFloat64Slice(field, split)
	case reflect.Int:
		return setIntSlice(field, split)
	}
	return nil
}

func setBoolSlice(field reflect.Value, split []string) error {
	boolSlice := make([]bool, 0, len(split))
	for _, str := range split {
		val, err := strconv.ParseBool(str)
		if err != nil {
			return fmt.Errorf("unable to set value for field %s. failed to parse %s as bool: %w", field.Type().Name(), str, err)
		}
		boolSlice = append(boolSlice, val)
	}
	field.Set(reflect.ValueOf(boolSlice))
	return nil
}

func setFloat64Slice(field reflect.Value, split []string) error {
	floatSlice := make([]float64, 0, len(split))
	for _, str := range split {
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("unable to set value for field %s. failed to parse %s as float64: %w", field.Type().Name(), str, err)
		}
		floatSlice = append(floatSlice, val)
	}
	field.Set(reflect.ValueOf(floatSlice))
	return nil
}

func setIntSlice(field reflect.Value, split []string) error {
	intSlice := make([]int, 0, len(split))
	for _, str := range split {
		val, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("unable to set value for field %s. failed to parse %s as int: %w", field.Type().Name(), str, err)
		}
		intSlice = append(intSlice, val)
	}
	field.Set(reflect.ValueOf(intSlice))
	return nil
}

// BindEnvWithAutoRefresh sets the values of the provided struct based on the values of the environment variables
// defined in the struct's tags. The struct must be a non-nil pointer to a struct. This function also refreshes the
// environment variables on a interval set with `AUTO_REFRESH_INTERVAL`.
// v: a non-nil pointer to a struct
// returns: an error if the provided value is not a non-nil pointer to a struct or if the value of an environment variable
func BindEnvWithAutoRefresh(v interface{}) error {
	if err := BindEnv(v); err != nil {
		return err
	}

	refresh(AUTO_REFRESH_INTERVAL, v)

	return nil
}

// refresh refreshes the environment variables
func refresh(interval int, v interface{}) {
	go func() {
		for {
			// sleep for the interval
			<-time.After(time.Duration(interval) * time.Second)
			if err := BindEnv(v); err != nil {
				fmt.Printf("failed to refresh environment variables: %s", err)
			}
		}
	}()
}
