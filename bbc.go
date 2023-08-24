package main

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	"os"
	//"sort"
	//"net/http"
	"golang.org/x/net/html"
	//"encoding/csv"
	"strings"
	"encoding/json"
	"strconv"
)



func tidyJSON(inputJSON BBCWeatherFormat) map[int]map[int]weatherFormat {
	var finalJSON map[int]map[int]weatherFormat = make(map[int]map[int]weatherFormat)
	//finalJSON["230824"] = make(map[string]weatherFormat)
	//finalJSON["230824"]["1400"] = testStruct{"thingA", 1276}
	
	for _, day := range inputJSON.Data.Forecasts {
		for _, timeSlot := range day.Detailed.Reports {
			curDate, _ := strconv.Atoi(strings.ReplaceAll(timeSlot.LocalDate, "-", "")[2:])
			curTime, _ := strconv.Atoi(timeSlot.Timeslot[0:2] + timeSlot.Timeslot[3:5])
			if finalJSON[curDate] == nil {
				finalJSON[curDate] = make(map[int]weatherFormat)
			}
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
				Visibility: timeSlot.Visibility,
				WeatherType: timeSlot.WeatherType,
				WeatherTypeText: timeSlot.WeatherTypeText,
				WindDirectionAbbreviation: timeSlot.WindDirectionAbbreviation,
				WindSpeedMph: timeSlot.WindSpeedMph,
			}
			finalJSON[curDate][curTime] = structToAdd
		}
	}
	return finalJSON
}	

func GetJSON() map[int]map[int]weatherFormat {
	//fmt.Println("In bbc.go")
	var finalJSONString string = ""
	homeDir, _ := os.UserHomeDir()
	dat, err := os.ReadFile(homeDir + FILE_PATH + "bbc.html")
	if err != nil {
		panic(err)
	}
	
	doc, err := html.Parse(strings.NewReader(string(dat)))
	if err != nil {
		panic(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
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
	
	// Format into JSON
	var formattedJSON BBCWeatherFormat
	if err := json.Unmarshal([]byte(finalJSONString), &formattedJSON); err != nil {
		panic(err)
	}
	return tidyJSON(formattedJSON)
}
