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
)

func DoThing() {
	homeDir, _ := os.UserHomeDir()
	dat, err := os.ReadFile(homeDir + FILE_PATH + "output.html")
	if err != nil {
		panic(err)
	}
	
	doc, err := html.Parse(strings.NewReader(string(dat)))
	if err != nil {
		panic(err)
	}
	var daysParsed int = 0 // There are 8 tables but only 7 days, so need to return before borking
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			fmt.Println("NEW TABLE")
			daysParsed++ // Return if done all 8 days
			if daysParsed > 7 {
				return
			}
			// Get node for header of table
			tHead := n.FirstChild.NextSibling
			_ = tHead.FirstChild
			// Parse header row and get items out
			var curRowNum int = 0 // To keep track of which row as each require different parsing
			
			// Get body and first row
			tBody := tHead.NextSibling.NextSibling
			curRow := tBody.FirstChild
			// Go through the siblings of the first row, i.e. iterate through rows
			for curRow.NextSibling != nil {
				if curRow.Data == "tr" { // If the sibling is another row, process it
					curRowNum++
					curCell := curRow.FirstChild // Get the first cell of the row
					for curCell.NextSibling != nil { // Iterate through cells 
						if curCell.Data == "td" { // If sibling is a cell
							switch { // Parse according to row number
							case curRowNum == 1: // Weather Symbol
								fmt.Println("Currently in: 1 (Weather Symbol)")
								fmt.Println(curCell.FirstChild.NextSibling.Attr[3].Val)
							case curRowNum == 2: // Chance of Precipitation
								fmt.Println("Currently in: 2 (Chance of rain)")
								fmt.Printf("%s\n", strings.Trim(curCell.FirstChild.Data, "\n"))
							case curRowNum == 3: // Actual Temperature
								fmt.Println("Currently in: 3 (Actual Temperature)")
								fmt.Printf("%s\n", curCell.FirstChild.NextSibling.Attr[0].Val)
							case curRowNum == 4 : // Feels Like Temperature
								fmt.Println("Currently in: 4 (Feels Like Temperature)")
								fmt.Printf("%s\n", curCell.Attr[2].Val)
							case curRowNum == 5: // Wind direction and speed
								fmt.Println("Currently in: 5 (Wind Direction and Speed)")
								fmt.Printf("%s\n", curCell.FirstChild.NextSibling.FirstChild.NextSibling.Attr[2].Val)
								fmt.Printf("%s\n", strings.Trim(curCell.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.Data, "\n"))
							case curRowNum == 6: // Wind Gust
								fmt.Println("Currently in: 6 (Wind Gust)")
								fmt.Printf("%s\n", curCell.FirstChild.NextSibling.FirstChild.Data)
							case curRowNum == 7: // Visibility
								fmt.Println("Currently in: 7 (Visibility)")
								fmt.Printf("%s\n", curCell.FirstChild.NextSibling.FirstChild.Data)
							case curRowNum == 8: // Humidity
								fmt.Println("Currently in: 8 (Humidity)")
								fmt.Printf("%s\n", strings.Trim(curCell.FirstChild.Data, "\n"))
							case curRowNum == 9: // UV
								fmt.Println("Currently in: 9 (UV)")
								fmt.Printf("%+v\n", curCell.FirstChild.NextSibling.Attr[4].Val)
							}
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
}
