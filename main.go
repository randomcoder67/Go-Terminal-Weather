package main

import (
	"fmt"
	//"strconv"
	//"strings"
	"encoding/json"
)

func main() {
	//fmt.Println("In main.go")
	bbcJSON := GetJSON()
	fmt.Println("Testing")
	for _, day := range bbcJSON {
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
	metofficeJSON, dayNames := ParseHTML()
	_ = metofficeJSON
	thingA, _ := json.MarshalIndent(bbcJSON, "", "  ")
	//fmt.Println(string(thingA))
	_ = thingA
	_ = dayNames
}
