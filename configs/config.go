package configs

import (
	"github.com/spf13/viper"
)

var ViperUnmarshal = viper.Unmarshal

type conf struct {
	WebServerPort      string `mapstructure:"WEB_SERVER_PORT"`
	CepAPIUrl          string `mapstructure:"CEP_API_URL"`
	WeatherAPIUrl      string `mapstructure:"WEATHER_API_URL"`
	WeatherAPIKey      string `mapstructure:"WEATHER_API_KEY"`
	WeatherAPILanguage string `mapstructure:"WEATHER_LANGUAGE"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = ViperUnmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	if cfg.WebServerPort == "" || cfg.CepAPIUrl == "" || cfg.WeatherAPIUrl == "" || cfg.WeatherAPIKey == "" {
		panic("missing required configuration")
	}

	return cfg, err
}
