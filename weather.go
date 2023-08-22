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

const FILE_PATH string = "/Programs/smallPrograms/metOffice/go/"
var homeDir string

func main() {
	homeDir, _ := os.UserHomeDir()
	/* 
	fmt.Println("Weather program")
	
	// Constuct filename and open file
	homeDir, _ = os.UserHomeDir()
	var fileName string = homeDir + FILE_PATH + "locations.csv"
	dat, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	if dat == nil {
		fmt.Println("Error, no locations")
	}
	
	// Read data into an array and print it, then get user input
	r := csv.NewReader(strings.NewReader(string(dat)))
	r.Comma = '|'
	records, _ := r.ReadAll()
	for i:=0; i<len(records); i++ {
		fmt.Printf("%d: %s\n", i+1, records[i][0])
	}
	fmt.Printf("Choose Location: ")
	var chosenLocation int
	_, err = fmt.Scanf("%d\n", &chosenLocation)
	
	// Get location code from user input, then create link and send get request
	var locationCode string = records[chosenLocation-1][1]
	var metOfficeLink string = "https://www.metoffice.gov.uk/weather/forecast/" + locationCode
	response, err := http.Get(metOfficeLink)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	
	// Read content of response into string
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var fullBody string = string(content)
	*/
	
	
	// Parse HTML 
	
	dat, err := os.ReadFile(homeDir + FILE_PATH + "output.html")
	if err != nil {
		panic(err)
	}
	/*
	
	z := html.NewTokenizer(strings.NewReader(string(dat)))
	
	depth := 0
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			fmt.Println(z.Err())
			return
		case html.TextToken:
			if depth > 0 {
				// emitBytes should copy the []byte it receives,
				// if it doesn't process it immediately.
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			if len(tn) == 1 && tn[0] == 'a' {
				if tt == html.StartTagToken {
					_, thing, _ := z.TagAttr()
					fmt.Println(string(thing))
					depth++
				} else {
					depth--
				}
			}
		}
	}
	*/
	
	/*
	var doThing bool = false
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.TextToken && doThing {
			//fmt.Println(tt)
			//t := z.Token()
			//fmt.Println(t.Data)
			doThing = false
		}
		if tt == html.StartTagToken || tt == html.EndTagToken {
			token := z.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						fmt.Println(attr.Val)
					}
				}
			}
		}
		//fmt.Println(tt)
	}
	*/
	
	doc, err := html.Parse(strings.NewReader(string(dat)))
	if err != nil {
		panic(err)
	}
	//var test int = 0
	//inTable, inBody, inRow, inCell := false, false, false, false
	var inTable bool = false
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			fmt.Println(n.Data)
			if !inTable {
				inTable = true
			} else if inTable {
				os.Exit(0)
			}
		} else if inTable {
			//fmt.Println(n.Data)
		}
		fmt.Println(n.Data)
		//fmt.Println(n.Data)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}
