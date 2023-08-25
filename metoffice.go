package main

import (
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"os"
	//"sort"
	//"net/http"
	"golang.org/x/net/html"
	//"encoding/csv"
	"strings"
	"strconv"
)

// Function to parse the MetOffice HTML into a map of maps of weatherFormat structs, and an array of the avalible dates
func ParseHTML() (map[int]map[string]weatherFormat, [7]string) {
	// Array for all the day names
	var dayNames = [7]string{}
	var dateNums = [7]int{}
	// final JSON to return
	var finalJSON map[int]map[string]weatherFormat = make(map[int]map[string]weatherFormat)
	// Open file and parse HTML
	homeDir, _ := os.UserHomeDir()
	dat, err := os.ReadFile(homeDir + FILE_PATH + "output.html")
	if err != nil {
		panic(err)
	}
	doc, err := html.Parse(strings.NewReader(string(dat)))
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
					finalJSON[dateNum] = make(map[string]weatherFormat)
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
					finalJSON[dateNums[daysParsed-1]][timeSlot] = weatherA  
				}
				curHeaderCell = curHeaderCell.NextSibling
			}
			
			// Parse the main table body
			var curRowNum int = 0 // To keep track of which row as each require different parsing
			var currentDateKey int = dateNums[daysParsed-1]
			fmt.Println(currentDateKey)
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
								fmt.Println("Currently in: 1 (Weather Symbol)")
								fmt.Println(timeSlots[atWhichTime])
								fmt.Println(curCell.FirstChild.NextSibling.Attr[3].Val)
							case curRowNum == 2: // Chance of Precipitation
								fmt.Println("Currently in: 2 (Chance of rain)")
								fmt.Println(timeSlots[atWhichTime])
								fmt.Printf("%s\n", strings.Trim(curCell.FirstChild.Data, "\n"))
							case curRowNum == 3: // Actual Temperature
								fmt.Println("Currently in: 3 (Actual Temperature)")
								fmt.Println(timeSlots[atWhichTime])
								fmt.Printf("%s\n", curCell.FirstChild.NextSibling.Attr[0].Val)
							case curRowNum == 4 : // Feels Like Temperature
								fmt.Println("Currently in: 4 (Feels Like Temperature)")
								fmt.Println(timeSlots[atWhichTime])
								if curCell.Attr[2].Key == "data-value" {
									fmt.Println(curCell.Attr[2].Val)
								} else {
									fmt.Println(curCell.Attr[1].Val)
								}
							case curRowNum == 5: // Wind direction and speed
								fmt.Println("Currently in: 5 (Wind Direction and Speed)")
								fmt.Println(timeSlots[atWhichTime])
								fmt.Printf("%s\n", curCell.FirstChild.NextSibling.FirstChild.NextSibling.Attr[2].Val)
								fmt.Printf("%s\n", strings.Trim(curCell.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.Data, "\n"))
							case curRowNum == 6: // Wind Gust
								fmt.Println("Currently in: 6 (Wind Gust)")
								fmt.Println(timeSlots[atWhichTime])
								fmt.Printf("%s\n", curCell.FirstChild.NextSibling.FirstChild.Data)
							case curRowNum == 7: // Visibility
								fmt.Println("Currently in: 7 (Visibility)")
								fmt.Println(timeSlots[atWhichTime])
								fmt.Printf("%s\n", curCell.FirstChild.NextSibling.FirstChild.Data)
							case curRowNum == 8: // Humidity
								fmt.Println("Currently in: 8 (Humidity)")
								fmt.Println(timeSlots[atWhichTime])
								fmt.Printf("%s\n", strings.Trim(curCell.FirstChild.Data, "\n"))
							case curRowNum == 9: // UV
								fmt.Println("Currently in: 9 (UV)")
								fmt.Println(timeSlots[atWhichTime])
								fmt.Printf("%+v\n", curCell.FirstChild.NextSibling.Attr[4].Val)
							}
							atWhichTime++
						}
						curCell = curCell.NextSibling
					} 
				}
				curRow = curRow.NextSibling
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
