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

// BindEnv sets the values of the provided struct based on the values of the environment variables
// defined in the struct's tags. The struct must be a non-nil pointer to a struct.
// v: a non-nil pointer to a struct
// returns: an error if the provided value is not a non-nil pointer to a struct or if the value of an environment variable
func BindEnv(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("provided value must be a non-nil pointer to a struct")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("provided value must be a pointer to a struct")
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if !field.CanSet() {
			continue
		}

		// if the field is a struct, recursively set the environment variables
		if field.Kind() == reflect.Struct {
			if err := BindEnv(field.Addr().Interface()); err != nil {
				return fmt.Errorf("unable to set value for field %s: %w", field.Type().Name(), err)
			}
			continue
		}

		// get the env tag for the field
		envTag := rt.Field(i).Tag.Get(ENV_TAG)
		if envTag == "" {
			continue
		}

		// get the value of the environment variable
		envValue := os.Getenv(envTag)
		if envValue == "" {
			// if the environment variable is not set, check if there is a default value
			defaultTag := rt.Field(i).Tag.Get(ENV_DEFAULT_TAG)
			if defaultTag != "" {
				envValue = defaultTag
			} else {
				continue
			}
		}

		switch field.Kind() {
		case reflect.String:
			// set the value of the field
			field.SetString(envValue)
		case reflect.Int:
			// convert the string value to an int
			val, err := strconv.Atoi(envValue)
			if err != nil {
				return fmt.Errorf("unable to set value for field %s. failed to parse %s as int: %w", field.Type().Name(), envValue, err)
			}
			field.SetInt(int64(val))
		case reflect.Bool:
			// convert the string value to a bool
			val, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("unable to set value for field %s. failed to parse %s as bool: %w", field.Type().Name(), envValue, err)
			}
			field.SetBool(val)
		case reflect.Float64:
			// convert the string value to a float64
			val, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return fmt.Errorf("unable to set value for field %s. failed to parse %s as float64: %w", field.Type().Name(), envValue, err)
			}
			field.SetFloat(val)
		case reflect.Struct:
			// recursively set the environment variables
			if err := BindEnv(field.Addr().Interface()); err != nil {
				return fmt.Errorf("unable to set value for field %s: %w", field.Type().Name(), err)
			}
		case reflect.Slice:
			// split the string value by comma and set the value of the field
			switch field.Type().Elem().Kind() {
			case reflect.String:
				split := strings.Split(envValue, ",")
				field.Set(reflect.ValueOf(split))
			case reflect.Bool:
				split := strings.Split(envValue, ",")
				boolSlice := make([]bool, 0, len(split))
				for _, str := range split {
					val, err := strconv.ParseBool(str)
					if err != nil {
						return fmt.Errorf("unable to set value for field %s. failed to parse %s as bool: %w", field.Type().Name(), envValue, err)
					}
					boolSlice = append(boolSlice, val)
				}
				field.Set(reflect.ValueOf(boolSlice))
			case reflect.Float64:
				split := strings.Split(envValue, ",")
				floatSlice := make([]float64, 0, len(split))
				for _, str := range split {
					val, err := strconv.ParseFloat(str, 64)
					if err != nil {
						return fmt.Errorf("unable to set value for field %s. failed to parse %s as float64: %w", field.Type().Name(), envValue, err)
					}
					floatSlice = append(floatSlice, val)
				}
				field.Set(reflect.ValueOf(floatSlice))
			case reflect.Int:
				split := strings.Split(envValue, ",")
				intSlice := make([]int, 0, len(split))
				for _, str := range split {
					val, err := strconv.Atoi(str)
					if err != nil {
						return fmt.Errorf("unable to set value for field %s. failed to parse %s as int: %w", field.Type().Name(), envValue, err)
					}
					intSlice = append(intSlice, val)
				}
				field.Set(reflect.ValueOf(intSlice))
			}
		}
	}

	return nil
}

// BindEnvWithAutoRefresh sets the values of the provided struct based on the values of the environment variables
// defined in the struct's tags. The struct must be a non-nil pointer to a struct.
// v: a non-nil pointer to a struct
// interval: the interval in seconds to refresh the environment variables
// returns: an error if the provided value is not a non-nil pointer to a struct or if the value of an environment variable
func BindEnvWithAutoRefresh(v interface{}, interval *int) error {
	if err := BindEnv(v); err != nil {
		return err
	}

	if interval == nil {
		defaultInterval := 60
		interval = &defaultInterval
	}

	refresh(*interval, v)

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
