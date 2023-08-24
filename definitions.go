package main

import (
	"time"
)

const FILE_PATH string = "/Programs/myRepos/goWeather/notPartOfFinal/"
// /home/user string, initailised in main.go
var homeDir string

// Weather symbols

// Credit: I used this website to convert the JSON into a struct. https://mholt.github.io/json-to-go/ 
// This struct is for the json that is contained in the BBC html page. 
type BBCWeatherFormat struct {
	Data struct {
		Forecasts []struct {
			Detailed struct {
				IssueDate   time.Time `json:"issueDate"`
				LastUpdated time.Time `json:"lastUpdated"`
				Reports	 []struct {
					EnhancedWeatherDescription		string `json:"enhancedWeatherDescription"`
					ExtendedWeatherType			   int	`json:"extendedWeatherType"`
					FeelsLikeTemperatureC			 int	`json:"feelsLikeTemperatureC"`
					FeelsLikeTemperatureF			 int	`json:"feelsLikeTemperatureF"`
					GustSpeedKph					  int	`json:"gustSpeedKph"`
					GustSpeedMph					  int	`json:"gustSpeedMph"`
					Humidity						  int	`json:"humidity"`
					LocalDate						 string `json:"localDate"`
					PrecipitationProbabilityInPercent int	`json:"precipitationProbabilityInPercent"`
					PrecipitationProbabilityText	  string `json:"precipitationProbabilityText"`
					Pressure						  int	`json:"pressure"`
					TemperatureC					  int	`json:"temperatureC"`
					TemperatureF					  int	`json:"temperatureF"`
					Timeslot						  string `json:"timeslot"`
					TimeslotLength					int	`json:"timeslotLength"`
					Visibility						string `json:"visibility"`
					WeatherType					   int	`json:"weatherType"`
					WeatherTypeText				   string `json:"weatherTypeText"`
					WindDescription				   string `json:"windDescription"`
					WindDirection					 string `json:"windDirection"`
					WindDirectionAbbreviation		 string `json:"windDirectionAbbreviation"`
					WindDirectionFull				 string `json:"windDirectionFull"`
					WindSpeedKph					  int	`json:"windSpeedKph"`
					WindSpeedMph					  int	`json:"windSpeedMph"`
				} `json:"reports"`
			} `json:"detailed"`
			Summary struct {
				IssueDate   time.Time `json:"issueDate"`
				LastUpdated time.Time `json:"lastUpdated"`
				Report	  struct {
					EnhancedWeatherDescription		string `json:"enhancedWeatherDescription"`
					GustSpeedKph					  int	`json:"gustSpeedKph"`
					GustSpeedMph					  int	`json:"gustSpeedMph"`
					LocalDate						 string `json:"localDate"`
					LowermaxTemperatureC			  any	`json:"lowermaxTemperatureC"`
					LowermaxTemperatureF			  any	`json:"lowermaxTemperatureF"`
					LowerminTemperatureC			  any	`json:"lowerminTemperatureC"`
					LowerminTemperatureF			  any	`json:"lowerminTemperatureF"`
					MaxTempC						  int	`json:"maxTempC"`
					MaxTempF						  int	`json:"maxTempF"`
					MinTempC						  int	`json:"minTempC"`
					MinTempF						  int	`json:"minTempF"`
					MostLikelyHighTemperatureC		int	`json:"mostLikelyHighTemperatureC"`
					MostLikelyHighTemperatureF		int	`json:"mostLikelyHighTemperatureF"`
					MostLikelyLowTemperatureC		 int	`json:"mostLikelyLowTemperatureC"`
					MostLikelyLowTemperatureF		 int	`json:"mostLikelyLowTemperatureF"`
					PollenIndex					   int	`json:"pollenIndex"`
					PollenIndexBand				   string `json:"pollenIndexBand"`
					PollenIndexIconText			   string `json:"pollenIndexIconText"`
					PollenIndexText				   string `json:"pollenIndexText"`
					PollutionIndex					int	`json:"pollutionIndex"`
					PollutionIndexBand				string `json:"pollutionIndexBand"`
					PollutionIndexIconText			string `json:"pollutionIndexIconText"`
					PollutionIndexText				string `json:"pollutionIndexText"`
					PrecipitationProbabilityInPercent int	`json:"precipitationProbabilityInPercent"`
					PrecipitationProbabilityText	  string `json:"precipitationProbabilityText"`
					Sunrise						   string `json:"sunrise"`
					Sunset							string `json:"sunset"`
					UppermaxTemperatureC			  any	`json:"uppermaxTemperatureC"`
					UppermaxTemperatureF			  any	`json:"uppermaxTemperatureF"`
					UpperminTemperatureC			  any	`json:"upperminTemperatureC"`
					UpperminTemperatureF			  any	`json:"upperminTemperatureF"`
					UvIndex						   int	`json:"uvIndex"`
					UvIndexBand					   string `json:"uvIndexBand"`
					UvIndexIconText				   string `json:"uvIndexIconText"`
					UvIndexText					   string `json:"uvIndexText"`
					WeatherType					   int	`json:"weatherType"`
					WeatherTypeText				   string `json:"weatherTypeText"`
					WindDescription				   string `json:"windDescription"`
					WindDirection					 string `json:"windDirection"`
					WindDirectionAbbreviation		 string `json:"windDirectionAbbreviation"`
					WindDirectionFull				 string `json:"windDirectionFull"`
					WindSpeedKph					  int	`json:"windSpeedKph"`
					WindSpeedMph					  int	`json:"windSpeedMph"`
				} `json:"report"`
			} `json:"summary"`
		} `json:"forecasts"`
		IsNight	 bool	  `json:"isNight"`
		IssueDate   time.Time `json:"issueDate"`
		LastUpdated time.Time `json:"lastUpdated"`
		Location	struct {
			Container string  `json:"container"`
			ID		string  `json:"id"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Name	  string  `json:"name"`
		} `json:"location"`
		Message string `json:"message"`
		Night   bool   `json:"night"`
	} `json:"data"`
	Environment  string `json:"environment"`
	FeatureFlags struct {
		UseAlgorithmicText bool `json:"useAlgorithmicText"`
	} `json:"featureFlags"`
	LocatorKey string `json:"locatorKey"`
	Options	struct {
		Day		string `json:"day"`
		Locale	 string `json:"locale"`
		LocationID string `json:"location_id"`
	} `json:"options"`
	UasKey		string `json:"uasKey"`
	WeatherAPIURI string `json:"weatherApiUri"`
}

// This struct is the target format for the JSON (removed a lot of unnecessary parts from the original BBC Weather JSON. Both the BBC Weather and MetOffice data will be converted to this struct format 
type weatherFormat struct {
	EnhancedWeatherDescription		string `json:"enhancedWeatherDescription"`
	ExtendedWeatherType			   int	`json:"extendedWeatherType"`
	Date							  string `json:"date"`
	Time							  string `json:"time"`
	FeelsLikeTemperatureC			 int	`json:"feelsLikeTemperatureC"`
	GustSpeedMph					  int	`json:"gustSpeedMph"`
	Humidity						  int	`json:"humidity"`
	PrecipitationProbabilityInPercent int	`json:"precipitationProbabilityInPercent"`
	TemperatureC					  int	`json:"temperatureC"`
	TimeslotLength					int	`json:"timeslotLength"`
	Visibility						string `json:"visibility"`
	WeatherType					   int	`json:"weatherType"`
	WeatherTypeText				   string `json:"weatherTypeText"`
	WindDirectionAbbreviation		 string `json:"windDirectionAbbreviation"`
	WindSpeedMph					  int	`json:"windSpeedMph"`
}
