package domain

import "errors"

type WeatherResponse struct {
	Location LocationData   `json:"location"`
	Current  CurrentWeather `json:"current"`
}

type LocationData struct {
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Timezone  string  `json:"tz_id"`
}

type WeatherCondition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
}

type CurrentWeather struct {
	TempK       float64
	TempC       float64          `json:"temp_c"`
	TempF       float64          `json:"temp_f"`
	Humidity    int              `json:"humidity"`
	WindKph     float64          `json:"wind_kph"`
	Condition   WeatherCondition `json:"condition"`
	LastUpdated string           `json:"last_updated"`
}

func (w *WeatherResponse) PopulateFromMap(data map[string]interface{}) error {
	location, ok := data["location"].(map[string]interface{})
	if !ok {
		return errors.New("invalid location data")
	}
	w.Location.Name = location["name"].(string)
	w.Location.Region = location["region"].(string)
	w.Location.Country = location["country"].(string)

	current, ok := data["current"].(map[string]interface{})
	if !ok {
		return errors.New("invalid current weather data")
	}
	w.Current.TempC = current["temp_c"].(float64)
	w.Current.TempF = current["temp_f"].(float64)
	w.Current.TempK = w.Current.TempC + 273.15

	condition, ok := current["condition"].(map[string]interface{})
	if ok {
		w.Current.Condition.Text = condition["text"].(string)
		w.Current.Condition.Icon = condition["icon"].(string)
	}
	return nil
}
