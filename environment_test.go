package ectoenv

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestBindEnv(t *testing.T) {
	type Config struct {
		StringValue  string    `env:"TEST_STRING"`
		IntValue     int       `env:"TEST_INT"`
		BoolValue    bool      `env:"TEST_BOOL"`
		FloatValue   float64   `env:"TEST_FLOAT"`
		StringSlice  []string  `env:"TEST_STRING_SLICE"`
		IntSlice     []int     `env:"TEST_INT_SLICE"`
		BoolSlice    []bool    `env:"TEST_BOOL_SLICE"`
		FloatSlice   []float64 `env:"TEST_FLOAT_SLICE"`
		DefaultValue string    `env:"TEST_DEFAULT" env-default:"default"`
	}

	tests := []struct {
		name     string
		envVars  map[string]string
		expected Config
	}{
		{
			name: "All values set",
			envVars: map[string]string{
				"TEST_STRING":       "test",
				"TEST_INT":          "42",
				"TEST_BOOL":         "true",
				"TEST_FLOAT":        "3.14",
				"TEST_STRING_SLICE": "a,b,c",
				"TEST_INT_SLICE":    "1,2,3",
				"TEST_BOOL_SLICE":   "true,false,true",
				"TEST_FLOAT_SLICE":  "1.1,2.2,3.3",
			},
			expected: Config{
				StringValue:  "test",
				IntValue:     42,
				BoolValue:    true,
				FloatValue:   3.14,
				StringSlice:  []string{"a", "b", "c"},
				IntSlice:     []int{1, 2, 3},
				BoolSlice:    []bool{true, false, true},
				FloatSlice:   []float64{1.1, 2.2, 3.3},
				DefaultValue: "default",
			},
		},
		{
			name:    "Default value",
			envVars: map[string]string{},
			expected: Config{
				DefaultValue: "default",
			},
		},
		{
			name: "Partial values set",
			envVars: map[string]string{
				"TEST_STRING": "partial",
				"TEST_BOOL":   "false",
			},
			expected: Config{
				StringValue:  "partial",
				BoolValue:    false,
				DefaultValue: "default",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// Clean up environment variables
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			var config Config
			err := BindEnv(&config)
			if err != nil {
				t.Fatalf("BindEnv() error = %v", err)
			}

			if !reflect.DeepEqual(config, tt.expected) {
				t.Errorf("BindEnv() got = %v, want %v", config, tt.expected)
			}
		})
	}
}

func TestBindEnvWithInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: "provided value must be a non-nil pointer to a struct",
		},
		{
			name:     "Non-pointer input",
			input:    struct{}{},
			expected: "provided value must be a non-nil pointer to a struct",
		},
		{
			name:     "Pointer to non-struct",
			input:    new(int),
			expected: "provided value must be a pointer to a struct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BindEnv(tt.input)
			if err == nil || err.Error() != tt.expected {
				t.Errorf("BindEnv() error = %v, want %v", err, tt.expected)
			}
		})
	}
}

func TestBindEnvWithInvalidValues(t *testing.T) {
	type InvalidConfig struct {
		IntValue   int     `env:"INVALID_INT"`
		BoolValue  bool    `env:"INVALID_BOOL"`
		FloatValue float64 `env:"INVALID_FLOAT"`
	}

	tests := []struct {
		name    string
		envVars map[string]string
		field   string
	}{
		{
			name:    "Invalid int",
			envVars: map[string]string{"INVALID_INT": "not_an_int"},
			field:   "IntValue",
		},
		{
			name:    "Invalid bool",
			envVars: map[string]string{"INVALID_BOOL": "not_a_bool"},
			field:   "BoolValue",
		},
		{
			name:    "Invalid float",
			envVars: map[string]string{"INVALID_FLOAT": "not_a_float"},
			field:   "FloatValue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// Clean up environment variables
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			var config InvalidConfig
			err := BindEnv(&config)
			if err == nil {
				t.Fatalf("BindEnv() expected error, got nil")
			}
			if !reflect.DeepEqual(config, InvalidConfig{}) {
				t.Errorf("BindEnv() unexpectedly modified config: %v", config)
			}
		})
	}
}

func TestBindEnvWithAutoRefresh(t *testing.T) {
	type Config struct {
		Value string `env:"TEST_AUTO_REFRESH"`
	}

	os.Setenv("TEST_AUTO_REFRESH", "initial")
	defer os.Unsetenv("TEST_AUTO_REFRESH")

	var config Config
	err := BindEnvWithAutoRefresh(&config)
	if err != nil {
		t.Fatalf("BindEnvWithAutoRefresh() error = %v", err)
	}

	if config.Value != "initial" {
		t.Errorf("Initial value not set correctly, got %v, want %v", config.Value, "initial")
	}

	// Change the environment variable
	os.Setenv("TEST_AUTO_REFRESH", "updated")

	// Wait for the refresh to occur
	time.Sleep(time.Duration(AUTO_REFRESH_INTERVAL+1) * time.Second)

	if config.Value != "updated" {
		t.Errorf("Value not updated after refresh, got %v, want %v", config.Value, "updated")
	}
}
