package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	//fmt.Println("In main.go")
	testThing := GetJSON()
	for _, day := range testThing {
		for _, timeSlot := range day {
			fmt.Printf("%+v\n", timeSlot)
			dateIndex, _ := strconv.Atoi(strings.ReplaceAll(timeSlot.Date, "-", "")[2:])
			timeSlotIndex, _ := strconv.Atoi(timeSlot.Time[0:2] + timeSlot.Time[3:5])
			fmt.Printf("%+v\n", testThing[dateIndex][timeSlotIndex])
		}
	}
}
