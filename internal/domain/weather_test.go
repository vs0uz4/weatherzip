package domain

import (
	"errors"
	"testing"
)

func TestWeatherResponsePopulateFromMap(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		expectErr error
		output    WeatherResponse
	}{
		{
			name: "Valid Data",
			input: map[string]interface{}{
				"location": map[string]interface{}{"name": "Cidade C", "region": "Região R", "country": "País P"},
				"current": map[string]interface{}{
					"temp_c": 25.0, "temp_f": 77.0, "condition": map[string]interface{}{"text": "Sunny", "icon": "icon_url"},
				},
			},
			output: WeatherResponse{
				Location: LocationData{Name: "Cidade C", Region: "Região R", Country: "País P"},
				Current: CurrentWeather{
					TempC: 25.0, TempF: 77.0, TempK: 298.15,
					Condition: WeatherCondition{Text: "Sunny", Icon: "icon_url"},
				},
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result WeatherResponse
			err := result.PopulateFromMap(tt.input)

			if !errors.Is(err, tt.expectErr) {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}

			if result != tt.output {
				t.Errorf("Expected output %+v, got %+v", tt.output, result)
			}
		})
	}
}
