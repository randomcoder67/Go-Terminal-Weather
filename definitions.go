package main

import (
	"time"
)

const FILE_PATH string = "/Programs/myRepos/goWeather/notPartOfFinal/"
// /home/user string, initailised in main.go
var homeDir string

// Row Labels
var rowLabels = [10]string{"Time", "Weather Symbol", "Change of Precipitation", "Temperature (°C)", "Feels like temperature (°C)", "Wind direction and speed", "Wind gust", "Visibility", "Humidity", "UV"}

// Weather symbols
var weatherSymbols = map[int]string{
	0: "",
	1: "",
	2: "  ",
	3: "",
	4: " ",
	5: "",
	6: "FOG",
	7: "",
	8: " ",
	9: "",
	10: "",
	11: "DRIZ",
	12: "",
	13: " ",
	14: " ",
	15: "",
	16: " ",
	17: " ",
	18: "",
	19: "HAIL ",
	20: "HAIL ",
	21: "HAIL",
	22: " ",
	23: " ",
	24: "",
	25: " ",
	26: " ",
	27: "",
	28: " ",
	29: " ",
	30: "",
	31: "",
	32: " ",
}

// Met Office weather name to BBC weather numbers (Not sure if MetOffice has Sandstorm, Hazy or Tropical storm but added them just incase)
var nameToWeatherNum = map[string]int{
	"Clear night": 0,
	"Sunny day": 1,
	"Partly cloudy (night)": 2,
	"Sunny intervals": 3,
	"Sandstorm": 4,
	"Mist": 5,
	"Fog": 6,
	"Cloudy": 7,
	"Overcast": 8,
	"Light shower (night)": 9,
	"Light shower (day)": 10,
	"Drizzle": 11,
	"Light rain": 12,
	"Heavy shower (night)": 13,
	"Heavy shower (day)": 14,
	"Heavy rain": 15,
	"Sleet shower (night)": 16,
	"Sleet shower (day)": 17,
	"Sleet": 18,
	"Hail shower (night)": 19,
	"Hail shower (day)": 20,
	"Hail": 21,
	"Light snow shower (night)": 22,
	"Light snow shower (day)": 23,
	"Light snow": 24,
	"Heavy snow shower (night)": 25,
	"Heavy snow shower (day)": 26,
	"Heavy snow": 27,
	"Thunder shower (night)": 28,
	"Thunder shower (day)": 29,
	"Thunder": 30,
	"Hazy": 31,
	"Tropical storm": 32,
}

// Wind direction abbreviation to symbols
var windDirection = map[string]string{
	"S": "",
	"W": "",
	"N": "",
	"E": "",
	"SW": " ",
	"NW": " ",
	"SE": " ",
	"NE": " ",
	"SSW": "",
	"WNW": "",
	"WSW": "",
	"NNW": "",
	"SSE": "",
	"ESE": "",
	"ENE": "",
	"NNE": "",
}

// Colours for the weather symbols, tried to make the colours so you could get an idea of the weather by just glancing at the symbol. Yellow = good and sunny, white = okay and cloudy/night, blue = wet, magenta = cold and snowy and red = danger (thunder or fog)
var weatherColours = map[int]string{
	0: "colorWhite",
	1: "colorYellow",
	2: "colorWhite",
	3: "colorYellow",
	4: "colorRed",
	5: "colorBlue",
	6: "colorRed",
	7: "colorWhite",
	8: "colorWhite",
	9: "colorBlue",
	10: "colorBlue",
	11: "colorBlue",
	12: "colorBlue",
	13: "colorBlue",
	14: "colorBlue",
	15: "colorBlue",
	16: "colorMagenta",
	17: "colorMagenta",
	18: "colorMagenta",
	19: "colorBlue",
	20: "colorBlue",
	21: "colorBlue",
	22: "colorMagenta",
	23: "colorMagenta",
	24: "colorMagenta",
	25: "colorMagenta",
	26: "colorMagenta",
	27: "colorMagenta",
	28: "colorRed",
	29: "colorRed",
	30: "colorRed",
	31: "colorRed",
	32: "colorRed",
}

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
