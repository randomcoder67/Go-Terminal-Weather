package main

import (
	"golang.org/x/net/html"
	"strings"
	"encoding/json"
	"fmt"
	//"strconv"
)

// Tidy up the JSON, reformatting it and removing unneeded information
func tidyJSON(inputJSON BBCWeatherFormat) map[int]map[string]weatherFormat {
	// Create the map
	var finalJSON map[int]map[string]weatherFormat = make(map[int]map[string]weatherFormat)
	//finalJSON["230824"] = make(map[string]weatherFormat)
	//finalJSON["230824"]["1400"] = testStruct{"thingA", 1276}
	//thing, _ := json.MarshalIndent(inputJSON, "", "  ")
	//fmt.Println(string(thing))
	// Iterate through all the days
	for i, day := range inputJSON.Data.Forecasts {
		// Iterate through all the time slots
		for _, timeSlot := range day.Detailed.Reports {
			// Get the current date in format yymmdd as an int
			//curDate, _ := strconv.Atoi(strings.ReplaceAll(timeSlot.LocalDate, "-", "")[2:])
			// Get the current time in format HHMM as an int
			curTime := timeSlot.Timeslot[0:2] + timeSlot.Timeslot[3:5]
			// Use current date to initialise parts of map
			if finalJSON[i] == nil {
				finalJSON[i] = make(map[string]weatherFormat)
			}
			// Create the struct to add and populate it
			var structToAdd weatherFormat = weatherFormat{
				EnhancedWeatherDescription: timeSlot.EnhancedWeatherDescription,
				ExtendedWeatherType: timeSlot.ExtendedWeatherType,
				Date: timeSlot.LocalDate,
				Time: timeSlot.Timeslot,
				FeelsLikeTemperatureC: timeSlot.FeelsLikeTemperatureC,
				GustSpeedMph: timeSlot.GustSpeedMph,
				Humidity: timeSlot.Humidity,
				PrecipitationProbabilityInPercent: timeSlot.PrecipitationProbabilityInPercent,
				TemperatureC: timeSlot.TemperatureC,
				TimeslotLength: timeSlot.TimeslotLength,
				Visibility: bbcWeatherVisibility[timeSlot.Visibility],
				WeatherType: timeSlot.WeatherType,
				WeatherTypeText: timeSlot.WeatherTypeText,
				WindDirectionAbbreviation: timeSlot.WindDirectionAbbreviation,
				WindSpeedMph: timeSlot.WindSpeedMph,
			}
			// Add the struct at the correct time slot
			finalJSON[i][curTime] = structToAdd
		}
	}
	// return the tidied JSON
	return finalJSON
}	

// Get the JSON out of BBC Weather HTML
func GetBBCWeatherFormatted(bbcWeatherHTML []byte) map[int]map[string]weatherFormat {
	fmt.Println("In bbcWeather.go")
	// Initialise the string to fill with JSON
	var finalJSONString string = ""
	doc, err := html.Parse(strings.NewReader(string(bbcWeatherHTML)))
	if err != nil {
		panic(err)
	}
	// Function to parse HTML
	var f func(*html.Node)
	f = func(n *html.Node) {
		// Get the JSON
		if n.Data == "script" {
			if len(n.Attr) > 1 {
				if n.Attr[1].Key == "data-state-id" {
					//fmt.Printf("%s\n", n.FirstChild.Data)
					finalJSONString = n.FirstChild.Data
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	
	// Format string into JSON
	var formattedJSON BBCWeatherFormat
	if err := json.Unmarshal([]byte(finalJSONString), &formattedJSON); err != nil {
		panic(err)
	}
	// Tidy the JSON up and return it to main
	return tidyJSON(formattedJSON)
}
