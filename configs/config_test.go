package configs

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigPanicOnUnmarshalError(t *testing.T) {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBufferString(`
app_name: MeuApp
port: invalid_number
`))
	if err != nil {
		t.Fatalf("Error when simulate configs: %v", err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Code should have returned a panic")
		}
	}()

	if _, err := LoadConfig(""); err == nil {
		t.Errorf("Expected error in LoadConfig, but no error was returned")
	}

}

func TestLoadConfigReadInConfigFails(t *testing.T) {
	invalidPath := "./invalid"

	assert.Panics(t, func() {
		_, _ = LoadConfig(invalidPath)
	}, "LoadConfig should panic when ReadInConfig fails")
}

func TestLoadConfig(t *testing.T) {
	envContent := `
WEB_SERVER_PORT=8080
CEP_API_URL=http://example.com/cep
WEATHER_API_URL=http://example.com/weather
WEATHER_API_KEY=testkey
WEATHER_LANGUAGE=en
`
	envFilePath := ".env"
	err := os.WriteFile(envFilePath, []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(envFilePath)

	cfg, err := LoadConfig(".")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "8080", cfg.WebServerPort)
	assert.Equal(t, "http://example.com/cep", cfg.CepAPIUrl)
	assert.Equal(t, "http://example.com/weather", cfg.WeatherAPIUrl)
	assert.Equal(t, "testkey", cfg.WeatherAPIKey)
	assert.Equal(t, "en", cfg.WeatherAPILanguage)
}

func TestLoadConfigUnmarshalFails(t *testing.T) {
	envContent := `
WEB_SERVER_PORT=8080
CEP_API_URL=1234
WEATHER_API_URL=
WEATHER_API_KEY=
WEATHER_LANGUAGE=
`
	envFilePath := ".env"
	err := os.WriteFile(envFilePath, []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(envFilePath)

	assert.Panics(t, func() {
		_, _ = LoadConfig(".")
	}, "LoadConfig should panic when Unmarshal fails")
}
