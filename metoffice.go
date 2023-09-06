package main

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
	"strconv"
	"math"
)

// Function to parse the MetOffice HTML into a map of maps of weatherFormat structs, and an array of the avalible dates
func GetMetOfficeFormatted(metofficeHTML []byte) (map[int]map[string]weatherFormat, [7]string) {
	fmt.Println("In metoffice.go")
	// Array for all the day names
	var dayNames = [7]string{}
	var dateNums = [7]int{}
	// final JSON to return
	var finalJSON map[int]map[string]weatherFormat = make(map[int]map[string]weatherFormat)
	doc, err := html.Parse(strings.NewReader(string(metofficeHTML)))
	if err != nil {
		panic(err)
	}
	var curDay int = 0 // Parse the 7 day headers to get the string (e.g. Today (24 August 2023))
	var daysParsed int = 0 // There are 8 tables but only 7 days, so need to return before borking
	// Iterate through the HTML and parse it as needed
	var f func(*html.Node)
	f = func(n *html.Node) {
		// Gets the avalible dates
		if n.Type == html.ElementNode && n.Data == "div" {
			if len(n.Attr) == 3 { // Check for the string "daynHeader" where n is 0-6
				if n.Attr[1].Val == "day" + strconv.Itoa(curDay) + "Header" {
					// Get the date in format yymmdd as an integer to use as key for map
					dateNum, _ := strconv.Atoi(strings.ReplaceAll(n.Attr[0].Val, "-", "")[2:])
					// Get day name (e.g. "Today", "Thursday")
					var dayName string = n.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.FirstChild.Data
					// Get full day name (e.g. "24 August 2023")
					var fullDayName string = n.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.FirstChild.NextSibling.FirstChild.Data
					// Add correctly formatted day names to the dayNames array (in format "Today (24 August 2023)")
					dayNames[curDay] = dayName + fullDayName
					// Initialse the day map using dateNum as the key
					finalJSON[curDay] = make(map[string]weatherFormat)
					dateNums[curDay] = dateNum
					// Advance to next day
					curDay++
				}
			}
		}
		// Parse the tables
		if n.Type == html.ElementNode && n.Data == "table" {
			var timeSlots = []string{}
			//fmt.Println("NEW TABLE")
			daysParsed++ // Return if done all 8 days
			if daysParsed > 7 {
				return
			}
			
			// Get node for header of table
			tHead := n.FirstChild.NextSibling
			//fmt.Printf("%+v\n", tHead.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.Attr[1].Val)
			// Iterate through header cells and get the time
			curHeaderCell := tHead.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling
			for curHeaderCell.NextSibling != nil {
				if curHeaderCell.Data == "th" {
					// Gets the time in format HHMM
					timeSlot := strings.ReplaceAll(curHeaderCell.Attr[1].Val, ":", "")
					timeSlots = append(timeSlots, timeSlot)
					var weatherA weatherFormat
					//fmt.Println(curDay)
					finalJSON[daysParsed-1][timeSlot] = weatherA  
				}
				curHeaderCell = curHeaderCell.NextSibling
			}
			//fmt.Println(timeSlots)
			
			
			
			// Parse the main table body
			var curRowNum int = 0 // To keep track of which row as each require different parsing
			var currentDateKey int = dateNums[daysParsed-1]
			//fmt.Println(currentDateKey)
			
			// Create array to hold struct and add empty structs equal to the number of time slots in the current day
			var tempStructHolder = []weatherFormat{}
			for i, _ := range timeSlots {
				var tempStruct weatherFormat
				tempStructHolder = append(tempStructHolder, tempStruct)
				tempStructHolder[i].Date = strconv.Itoa(currentDateKey)
			}
			
			// Get body and first row
			tBody := tHead.NextSibling.NextSibling
			curRow := tBody.FirstChild
			// Go through the siblings of the first row, i.e. iterate through rows
			for curRow.NextSibling != nil {
				if curRow.Data == "tr" { // If the sibling is another row, process it
					curRowNum++
					curCell := curRow.FirstChild // Get the first cell of the row
					var atWhichTime int = 0
					for curCell.NextSibling != nil { // Iterate through cells
						if curCell.Data == "td" { // If sibling is a cell
							switch { // Parse according to row number
							case curRowNum == 1: // Weather Symbol
								tempStructHolder[atWhichTime].Time = timeSlots[atWhichTime]
								tempStructHolder[atWhichTime].WeatherType = nameToWeatherNum[curCell.FirstChild.NextSibling.Attr[3].Val]
							case curRowNum == 2: // Chance of Precipitation
								//fmt.Println("Currently in: 2 (Chance of rain)")
								//fmt.Println(timeSlots[atWhichTime])
								value, _ := strconv.Atoi(strings.Trim(
									strings.Trim(strings.Trim(curCell.FirstChild.Data, "\n"), "%"), "<"))
								tempStructHolder[atWhichTime].PrecipitationProbabilityInPercent = value
							case curRowNum == 3: // Actual Temperature
								//fmt.Println("Currently in: 3 (Actual Temperature)")
								//fmt.Println(timeSlots[atWhichTime])
								//fmt.Printf("%s\n", curCell.FirstChild.NextSibling.Attr[0].Val)
								value, _ := strconv.ParseFloat(curCell.FirstChild.NextSibling.Attr[0].Val, 64)
								valueInt := int(math.Round(value))
								tempStructHolder[atWhichTime].TemperatureC = valueInt
							case curRowNum == 4 : // Feels Like Temperature
								//fmt.Println("Currently in: 4 (Feels Like Temperature)")
								//fmt.Println(timeSlots[atWhichTime])
								var valueString string
								if curCell.Attr[2].Key == "data-value" {
									valueString = curCell.Attr[2].Val
								} else {
									valueString = curCell.Attr[1].Val
								}
								value, _ := strconv.ParseFloat(valueString, 64)
								valueInt := int(math.Round(value))
								tempStructHolder[atWhichTime].FeelsLikeTemperatureC = valueInt
							case curRowNum == 5: // Wind direction and speed
								//fmt.Println("Currently in: 5 (Wind Direction and Speed)")
								//fmt.Println(timeSlots[atWhichTime])
								tempStructHolder[atWhichTime].WindDirectionAbbreviation = curCell.FirstChild.NextSibling.FirstChild.NextSibling.Attr[2].Val
								value, _ := strconv.Atoi(strings.Trim(curCell.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.Data, "\n"))
								tempStructHolder[atWhichTime].WindSpeedMph = value
							case curRowNum == 6: // Wind Gust
								//fmt.Println("Currently in: 6 (Wind Gust)")
								//fmt.Println(timeSlots[atWhichTime])
								value, _ := strconv.Atoi(curCell.FirstChild.NextSibling.FirstChild.Data)
								tempStructHolder[atWhichTime].GustSpeedMph = value
							case curRowNum == 7: // Visibility
								//fmt.Println("Currently in: 7 (Visibility)")
								//fmt.Println(timeSlots[atWhichTime])
								tempStructHolder[atWhichTime].Visibility = curCell.FirstChild.NextSibling.FirstChild.Data
							case curRowNum == 8: // Humidity
								//fmt.Println("Currently in: 8 (Humidity)")
								//fmt.Println(timeSlots[atWhichTime])
								value, _ := strconv.Atoi(strings.Trim(strings.Trim(curCell.FirstChild.Data, "\n"), "%"))
								tempStructHolder[atWhichTime].Humidity = value
							case curRowNum == 9: // UV
								//fmt.Println("Currently in: 9 (UV)")
								//fmt.Println(timeSlots[atWhichTime])
								//fmt.Printf("%+v\n", curCell.FirstChild.NextSibling.Attr[4].Val)
							}
							// Set time slot length correctly
							if curRowNum == 1 {
								if curDay == 1 || curDay == 2 { // Today and Tomorrow = per hour
									tempStructHolder[atWhichTime].TimeslotLength = 1
								} else if curDay > 3 { // 3 days away+ = per 3 hour
									tempStructHolder[atWhichTime].TimeslotLength = 3
								} else if curDay == 3 { // 2 days away = 3 hour except first (0000)
									if timeSlots[atWhichTime] == "0000" {
										tempStructHolder[atWhichTime].TimeslotLength = 1
									} else {
										tempStructHolder[atWhichTime].TimeslotLength = 1
									}
								}
							}
							atWhichTime++
						}
						curCell = curCell.NextSibling
					} 
				}
				curRow = curRow.NextSibling
			}
			//fmt.Printf("%+v\n", tempStructHolder[0])
			for i, timeSlot := range timeSlots {
				//fmt.Println(i, timeSlot)
				finalJSON[daysParsed-1][timeSlot] = tempStructHolder[i]
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	// Return the weather data and array of day names
	return finalJSON, dayNames
}
