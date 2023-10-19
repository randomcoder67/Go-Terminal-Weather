package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	"sort"
	"strconv"
	//"time"
	//"encoding/json"
)

var curDay int = 0
var curPage int = 0
var showLegend bool = false
var currentForecast int = 0
var styles = map[string]tcell.Style{}
var termHeight int
var termWidth int

// Function to allow adding strings to screen instead of just individual runes (chars)
func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

// Function to draw the current day and page
func drawData(s tcell.Screen, weatherData map[int]map[string]weatherFormat, dayNames [7]string) {
	s.Clear()
	// Get currently selected provider
	var currentProvider string = "MetOffice"
	if currentForecast == 1 { currentProvider = "BBC Weather" }
	
	// Draw the day name
	drawText(s, 0, 0, 50, 10, styles["pink"], dayNames[curDay] + " (" + currentProvider + "):")
	// Draw the row labels
	for i, label := range rowLabels {
		drawText(s, 1, 2+2*i, 50, 20, styles["yellow"], label + ":")
	}
	
	// Make slice of all the timeslots (used as maps are not ordered)
	timeSlots := []string{}
	for timeSlot, _ := range weatherData[curDay] {
		timeSlots = append(timeSlots, timeSlot)
	}
	sort.SliceStable(timeSlots, func(i, j int) bool{ // Sort timeslots 
		return timeSlots[i] < timeSlots[j]
	})
	
	
	var curTimeSlot int = 0
	
	render := func(splitDay bool) {
		// Handle days with more than 9 timeslots, can't fit on one screen
		timeSlotsFinal := timeSlots // By default, the time slots array is fine
		if splitDay { // Unless told to split
			if len(timeSlots) <= 9 { // If the length is <= 9, (i.e. at 2100), do nothing and reset curPage
				curPage = 0
			} else { // Otherwise 
				if len(timeSlots) <= 18 && curPage == 2 { // If length is <= 18, don't allow curPage > 1
					curPage--
				} else { // Change timeSlotsFinal to the current page
					var startIndex int = curPage*9
					var endIndex int = startIndex + 9
					if endIndex > len(timeSlots) { endIndex = len(timeSlots) }
					timeSlotsFinal = timeSlots[startIndex:endIndex]
				}
			}
		}
		
		for _, timeSlot := range timeSlotsFinal {
			curWeatherStruct := weatherData[curDay][timeSlot]
			
			// Time
			drawText(s, 30+curTimeSlot*8, 2, 50+curTimeSlot*8, 20, styles["pink"], timeSlot[0:2] + ":" + timeSlot[2:4])
			
			// Weather Symbol
			drawText(s, 30+curTimeSlot*8, 4, 50+curTimeSlot*8, 20, styles[weatherColours[curWeatherStruct.WeatherType]], weatherSymbols[curWeatherStruct.WeatherType])
			
			// Chance of Precipitation
			var precipInt int = curWeatherStruct.PrecipitationProbabilityInPercent
			var precipString string = strconv.Itoa(precipInt) + "%"
			var precipColour tcell.Style
			if precipInt < 15 {
				precipColour = styles["green"]
			} else if precipInt < 30 {
				precipColour = styles["blue"]
			} else if precipInt < 65 {
				precipColour = styles["yellow"]
			} else {
				precipColour = styles["red"]
			}
			if precipInt == 5 {
				precipString = "<" + precipString
			}
			drawText(s, 30+curTimeSlot*8, 6, 50+curTimeSlot*8, 20, precipColour, precipString)
			
			// Temperature
			var tempInt int = curWeatherStruct.TemperatureC
			var tempColour tcell.Style
			if tempInt < 4 {
				tempColour = styles["blue"]
			} else if tempInt < 14 {
				tempColour = styles["magenta"]
			} else if tempInt < 18 {
				tempColour = styles["green"]
			} else if tempInt < 23 {
				tempColour = styles["yellow"]
			} else {
				tempColour = styles["red"]
			}
			drawText(s, 30+curTimeSlot*8, 8, 50+curTimeSlot*8, 20, tempColour, strconv.Itoa(tempInt) + "°")
			
			// Feels like temperature
			var feelsTempInt int = curWeatherStruct.FeelsLikeTemperatureC
			var feelsTempColour tcell.Style
			if feelsTempInt < 4 {
				feelsTempColour = styles["blue"]
			} else if feelsTempInt < 14 {
				feelsTempColour = styles["magenta"]
			} else if feelsTempInt < 18 {
				feelsTempColour = styles["green"]
			} else if feelsTempInt < 23 {
				feelsTempColour = styles["yellow"]
			} else {
				feelsTempColour = styles["red"]
			}
			drawText(s, 30+curTimeSlot*8, 10, 50+curTimeSlot*8, 20, feelsTempColour, strconv.Itoa(feelsTempInt) + "°")
			
			// Wind direction and speed
			drawText(s, 30+curTimeSlot*8, 12, 50+curTimeSlot*8, 20, styles["white"], windDirection[curWeatherStruct.WindDirectionAbbreviation])
			
			var windInt int = curWeatherStruct.WindSpeedMph
			var windColour tcell.Style
			if windInt < 4 {
				windColour = styles["white"]
			} else if windInt < 10 {
				windColour = styles["green"]
			} else if windInt < 18 {
				windColour = styles["magenta"]
			} else if windInt < 25 {
				windColour = styles["yellow"]
			} else {
				windColour = styles["red"]
			}
			drawText(s, 30+curTimeSlot*8, 13, 50+curTimeSlot*8, 20, windColour, strconv.Itoa(windInt) + " mph")
			
			// Wind gust
			var gustInt int = curWeatherStruct.GustSpeedMph
			var gustColour tcell.Style
			if gustInt < 4 {
				gustColour = styles["white"]
			} else if gustInt < 12 {
				gustColour = styles["green"]
			} else if gustInt < 24 {
				gustColour = styles["magenta"]
			} else if gustInt < 30 {
				gustColour = styles["yellow"]
			} else {
				gustColour = styles["red"]
			}
			drawText(s, 30+curTimeSlot*8, 14, 50+curTimeSlot*8, 20, gustColour, strconv.Itoa(gustInt) + " mph")
			
			// Visibility
			var visString string = curWeatherStruct.Visibility
			var visColour tcell.Style
			if visString == "UN" {
				visColour = styles["white"]
			} else if visString == "VP" {
				visColour = styles["red"]
			} else if visString == "P" {
				visColour = styles["red"]
			} else if visString == "M" {
				visColour = styles["yellow"]
			} else if visString == "G" {
				visColour = styles["yellow"]
			} else if visString == "VG" {
				visColour = styles["green"]
			} else if visString == "E" {
				visColour = styles["green"]
			}
			drawText(s, 30+curTimeSlot*8, 16, 50+curTimeSlot*8, 20, visColour, visString)
			
			// Humidity
			drawText(s, 30+curTimeSlot*8, 18, 50+curTimeSlot*8, 20, styles["blue"], strconv.Itoa(curWeatherStruct.Humidity) + "%")
			
			//drawText(s, 
			curTimeSlot++
			
			//drawText(s, termWidth/2 - 2, termHeight-3, termWidth/2 + 6, termHeight-1, styles["white"], strconv.Itoa(curPage) + " " + strconv.Itoa(len(timeSlots)))
			
			// Render down arrows if necessary
			if curPage*9+9 < len(timeSlots) && splitDay {
				drawText(s, termWidth/2 - 2, termHeight-1, termWidth/2 + 6, termHeight-1, styles["white"], DOWN_ARROWS)
			}
			// Render up arrows if necessary
			if curPage*9+9 > 9 && splitDay {
				drawText(s, termWidth/2 - 2, 0, termWidth/2 + 6, 0, styles["white"], UP_ARROWS)
			}
		}
	}
	
	// Only today and tomorrow need multiple pages
	if curDay < 2 || currentForecast == 1 {
		render(true)
	} else {
		render(false)
	}
}
	

func DoDisplay(metofficeJSON map[int]map[string]weatherFormat, bbcWeatherJSON map[int]map[string]weatherFormat, dayNames [7]string) {

	//thing, _ := json.MarshalIndent(bbcWeatherJSON, "", "  ")
	//fmt.Println(string(thing))

	fmt.Println("In display.go")
	
	var weatherData = [2]map[int]map[string]weatherFormat{metofficeJSON, bbcWeatherJSON}
	
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = s.Init()
	if err != nil {
		panic(err)
	}
	
	// Create colours
	yellowColour := tcell.NewHexColor(0xffff80)
	redColour := tcell.NewHexColor(0xe64747)
	blueColour := tcell.NewHexColor(0x38ffea)
	magentaColour := tcell.NewHexColor(0x9a7cff)
	greenColour := tcell.NewHexColor(0x52ff6d)
	pinkColour := tcell.NewHexColor(0xff76c1)
	
	// Use the colours to make styles and add to "styles" map
	styles["white"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	styles["yellow"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(yellowColour)
	styles["red"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(redColour)
	styles["blue"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(blueColour)
	styles["magenta"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(magentaColour)
	styles["green"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(greenColour)
	styles["pink"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(pinkColour)
	
	s.SetStyle(styles["white"])
	s.Clear()
	
	drawData(s, weatherData[currentForecast], dayNames)
	termWidth, termHeight = s.Size()
	
	//thing := time.Now()
	quit := func() {
		s.Fini()
		//fmt.Println("After Display Setup:", thing)
		os.Exit(0)
	}
	for {
		s.Show()
		
		ev := s.PollEvent()
		
		switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyLeft { // If left pressed, curDay-- unless it's 0
					if curDay != 0 { curDay-- }
					drawData(s, weatherData[currentForecast], dayNames)
				} else if ev.Key() == tcell.KeyRight { // If right pressed, curDay++ unless it's 2
					if curDay != 6 { curDay++ }
					drawData(s, weatherData[currentForecast], dayNames)
				} else if ev.Key() == tcell.KeyUp { // If up pressed, go up unless at top
					if curPage != 0 { curPage-- }
					drawData(s, weatherData[currentForecast], dayNames)
				} else if ev.Key() == tcell.KeyDown  { // If down pressed, go down unless at bottom
					if curPage != 2 { curPage++ }
					drawData(s, weatherData[currentForecast], dayNames)
				// Esc, Ctrl-C and q all quit
				} else if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					quit()
				} else if ev.Rune() == ' ' { // Space switches between MetOffice and BBC Weather data
					if currentForecast == 0 {
						currentForecast = 1
					} else {
						currentForecast = 0
					}
					drawData(s, weatherData[currentForecast], dayNames)
				}
			}
		}
}
