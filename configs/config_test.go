package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
