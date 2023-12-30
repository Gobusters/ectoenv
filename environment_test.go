package ectoenv

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Example struct for testing
type Config struct {
	StringValue      string   `env:"STRING_VALUE" env-default:"default"`
	IntValue         int      `env:"INT_VALUE" env-default:"42"`
	BoolValue        bool     `env:"BOOL_VALUE" env-default:"true"`
	IntSliceValue    []int    `env:"INT_SLICE_VALUE" env-default:"1,2,3,4,5"`
	StringSliceValue []string `env:"STRING_SLICE_VALUE" env-default:"a,b,c,d,e"`
	BoolSliceValue   []bool   `env:"BOOL_SLICE_VALUE" env-default:"true,false,true,false,true"`
	SubStruct        struct {
		SubStringValue string `env:"SUB_STRING_VALUE" env-default:"subdefault"`
	}
}

func TestBindEnv_ValidStructPointer(t *testing.T) {
	os.Setenv("STRING_VALUE", "test")
	os.Setenv("INT_VALUE", "42")
	os.Setenv("BOOL_VALUE", "true")
	os.Setenv("INT_SLICE_VALUE", "1,2,3,4,5")
	os.Setenv("STRING_SLICE_VALUE", "a,b,c,d,e")
	os.Setenv("BOOL_SLICE_VALUE", "true,false,true,false,true")
	os.Setenv("SUB_STRING_VALUE", "subtest")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)

	assert.Nil(t, err)
	assert.Equal(t, "test", config.StringValue)
	assert.Equal(t, 42, config.IntValue)
	assert.Equal(t, true, config.BoolValue)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, config.IntSliceValue)
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, config.StringSliceValue)
	assert.Equal(t, []bool{true, false, true, false, true}, config.BoolSliceValue)
	assert.Equal(t, "subtest", config.SubStruct.SubStringValue)
}

func TestBindEnv_NilPointer(t *testing.T) {
	err := BindEnv(nil)
	assert.NotNil(t, err)
}

func TestBindEnv_NonPointerInput(t *testing.T) {
	var config Config
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_NonStructPointer(t *testing.T) {
	var intValue int
	err := BindEnv(&intValue)
	assert.NotNil(t, err)
}

func TestBindEnv_NonSettableField(t *testing.T) {
	var config struct {
		NonSettable string `env:"NON_SETTABLE"`
	}
	err := BindEnv(&config)
	assert.Nil(t, err)
}

func TestBindEnv_InvalidIntValue(t *testing.T) {
	os.Setenv("INT_VALUE", "notanint")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_InvalidBoolValue(t *testing.T) {
	os.Setenv("BOOL_VALUE", "notabool")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_InvalidIntSliceValue(t *testing.T) {
	os.Setenv("INT_SLICE_VALUE", "notanint")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_InvalidBoolSliceValue(t *testing.T) {
	os.Setenv("BOOL_SLICE_VALUE", "notabool")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_InvalidSliceValueElement(t *testing.T) {
	os.Setenv("INT_SLICE_VALUE", "1,2,3,4,notanint")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_InvalidSliceValueElementBool(t *testing.T) {
	os.Setenv("BOOL_SLICE_VALUE", "true,false,true,false,notabool")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_InvalidSliceValueElementBoolLength(t *testing.T) {
	os.Setenv("BOOL_SLICE_VALUE", "true,false,true,false,")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_InvalidSliceValueElementInt(t *testing.T) {
	os.Setenv("INT_SLICE_VALUE", "1,2,3,4,notanint")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_InvalidSliceValueElementIntLength(t *testing.T) {
	os.Setenv("INT_SLICE_VALUE", "1,2,3,4,")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_EmptySliceValue(t *testing.T) {
	os.Setenv("INT_SLICE_VALUE", "")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.Nil(t, err)
}

func TestBindEnv_EmptySliceValueElement(t *testing.T) {
	os.Setenv("STRING_SLICE_VALUE", "a,b,c,d,")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.Nil(t, err)
	assert.Equal(t, []string{"a", "b", "c", "d", ""}, config.StringSliceValue)
}

func TestBindEnv_EmptySliceValueElementBool(t *testing.T) {
	os.Setenv("BOOL_SLICE_VALUE", "true,false,true,false,")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestBindEnv_EmptySliceValueElementInt(t *testing.T) {
	os.Setenv("INT_SLICE_VALUE", "1,2,3,4,")
	defer os.Clearenv()

	config := &Config{}
	err := BindEnv(config)
	assert.NotNil(t, err)
}

func TestEnvironmentDefaultTag(t *testing.T) {
	config := &Config{}
	err := BindEnv(config)

	assert.Nil(t, err)
	assert.Equal(t, "default", config.StringValue)
	assert.Equal(t, 42, config.IntValue)
	assert.Equal(t, true, config.BoolValue)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, config.IntSliceValue)
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, config.StringSliceValue)
	assert.Equal(t, []bool{true, false, true, false, true}, config.BoolSliceValue)
	assert.Equal(t, "subdefault", config.SubStruct.SubStringValue)
}

// test refresh env
func TestRefreshEnv(t *testing.T) {
	// set env
	os.Setenv("STRING_VALUE", "test")
	defer os.Clearenv()

	config := &Config{}
	interval := 1
	AUTO_REFRESH_INTERVAL = interval
	err := BindEnvWithAutoRefresh(config)
	assert.Nil(t, err)

	assert.Equal(t, "test", config.StringValue)

	// change env
	os.Setenv("STRING_VALUE", "test2")

	// wait 1 second for refresh
	<-time.After(time.Millisecond * 1100)

	assert.Equal(t, "test2", config.StringValue)
}
