package main

import (
	"fmt"
	"encoding/csv"
	"os"
	"strings"
	"bufio"
	"strconv"
	"time"
)


func main() {
	fmt.Println("In main.go")
	// Get home dir
	homeDir, _ = os.UserHomeDir()
	// Read locations file 
	dat, err := os.ReadFile(homeDir + LOCATIONS_FILE_LOC)
	if err != nil {
		panic(err)
	}
	// Parse as csv
	r := csv.NewReader(strings.NewReader(string(dat)))
	r.Comma = '|'
	locations, _ := r.ReadAll()
	
	var metofficeCode string
	var bbcWeatherCode string
	// Display options to user if more than one, otherwise use only entry
	if len(locations) > 1 {
		// Print info to user
		fmt.Println("Welcome to Go Weather, select your desired location:")
		for i, location := range locations {
			fmt.Printf("%d. %s\n", i+1, location[0])
		}
		// Get input
		in := bufio.NewReader(os.Stdin)
		givenIndex, err := in.ReadString('\n')
		// Convert to integer
		givenIndexInt, err := strconv.Atoi(strings.ReplaceAll(givenIndex, "\n", ""))
		// Check input was valid
		if err != nil {
			fmt.Println("Error, invalid input")
			os.Exit(1)
		}
		// And in range
		if givenIndexInt > len(locations) || givenIndexInt == 0 {
			fmt.Println("Error, input out of range")
			os.Exit(1)
		}
		// Get codes
		metofficeCode = locations[givenIndexInt-1][1]
		bbcWeatherCode = locations[givenIndexInt-1][2]
	} else {
		metofficeCode = locations[0][1]
		bbcWeatherCode = locations[0][2]
	}
	fmt.Println(metofficeCode, bbcWeatherCode)
	fmt.Println("Before Download:", time.Now())
	metofficeHTML, bbcWeatherHTML := RetrieveHTML(metofficeCode, bbcWeatherCode)
	fmt.Println("After Download:", time.Now())
	
	metofficeJSON, dayNames := GetMetOfficeFormatted(metofficeHTML)
	fmt.Println("After MetOffice:", time.Now())
	bbcWeatherJSON := GetBBCWeatherFormatted(bbcWeatherHTML)
	fmt.Println("After BBC Weather:", time.Now())
	/*
	for _, day := range bbcWeatherJSON {
		for _, timeSlot := range day {
			_ = timeSlot
			//fmt.Println(timeSlot.WeatherType)
			//fmt.Println(timeSlot.ExtendedWeatherType)
			//fmt.Println("TIME SLOT COMPLETE")
			//fmt.Printf("%+v\n", timeSlot)
			//dateIndex, _ := strconv.Atoi(strings.ReplaceAll(timeSlot.Date, "-", "")[2:])
			//timeSlotIndex, _ := strconv.Atoi(timeSlot.Time[0:2] + timeSlot.Time[3:5])
			//fmt.Printf("%+v\n", testThing[dateIndex][timeSlotIndex])
		}
	}
	*/
	//_ = metofficeJSON
	//thingA, _ := json.MarshalIndent(metofficeJSON, "", "  ")
	//fmt.Println(string(thingA))
	//_ = thingA
	//_ = dayNames
	//fmt.Printf("%+v\n", metofficeJSON[230803]["0100"].Date)
	//fmt.Println(bbcWeatherJSON)
	//os.Exit(0)
	DoDisplay(metofficeJSON, bbcWeatherJSON, dayNames)
}
